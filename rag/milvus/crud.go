package milvus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	log "github.com/sirupsen/logrus"
)

type RagCli struct {
	cli *milvusclient.Client
}

type UpsertReq struct {
	Vector  [][]float32
	Content []string
}

type QueryResp struct {
	Score   float32
	Content string
}

func NewRagCli(ctx context.Context) *RagCli {
	return &RagCli{cli: NewDefaultClient(ctx)}
}

func Sha256(content []byte) string {
	// 计算 SHA256
	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])
	return hashStr
}

func (c *RagCli) Close(ctx context.Context) error {
	return c.cli.Close(ctx)
}

func (c *RagCli) Upsert(ctx context.Context, req UpsertReq) (int64, error) {
	if len(req.Content) != len(req.Vector) {
		return -1, fmt.Errorf("vecs: %d != content: %d", len(req.Content), len(req.Vector))
	}
	n := len(req.Content)
	if n < 1 {
		return 0, nil
	}

	idx := make([]string, 0, n)
	for i := range n {
		b := []byte(req.Content[i])
		idx = append(idx, Sha256(b))
	}

	opts := milvusclient.NewColumnBasedInsertOption(defaultCollection)
	opts = opts.WithFloatVectorColumn(vectorField, defaultDim, req.Vector)
	opts = opts.WithVarcharColumn(contentField, req.Content)
	opts = opts.WithVarcharColumn(idField, idx)

	ret, err := c.cli.Upsert(ctx, opts)

	return ret.UpsertCount, err
}

// https://milvus.io/api-reference/go/v2.6.x/Vector/Search.md
func (c *RagCli) Search(ctx context.Context, limit int, qvecs ...entity.Vector) ([]QueryResp, error) {
	n := len(qvecs)
	if n < 1 {
		return nil, nil
	}

	opts := milvusclient.NewSearchOption(
		defaultCollection,
		limit,
		qvecs,
	)
	opts = opts.WithConsistencyLevel(entity.ClStrong)
	opts = opts.WithANNSField(vectorField)
	opts = opts.WithOutputFields(contentField)

	ret, err := c.cli.Search(ctx, opts)
	if err != nil {
		return nil, err
	}

	var qresp []QueryResp

	for _, resp := range ret {
		log.WithField("resp_len", resp.Len()).Info("data fetched")
		log.Infof("%+v", resp.Fields)

		col := resp.GetColumn(contentField)

		for j := range col.Len() {
			content, err := col.GetAsString(j)
			if err != nil {
				log.WithError(err).Error("error fetch search result")
				continue
			}
			qresp = append(qresp, QueryResp{
				Score:   resp.Scores[j],
				Content: content,
			})
		}

	}

	log.WithFields(log.Fields{
		"search_vec": len(qvecs),
		"result_set": len(ret),
	}).Info("search done")

	return qresp, nil
}
