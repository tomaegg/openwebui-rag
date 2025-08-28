package milvus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	log "github.com/sirupsen/logrus"
)

type RagCli struct {
	cli *milvusclient.Client
}

type InsertReq struct {
	Vector  []float32
	Content []byte
}

type QueryResp struct {
	Scores  float32
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

func (c *RagCli) Upsert(ctx context.Context, req ...InsertReq) (int64, error) {
	n := len(req)
	if n < 1 {
		return 0, nil
	}
	vecs := make([][]float32, 0, n)
	contents := make([]string, 0, n)
	idx := make([]string, 0, n)

	for i := range n {
		vecs = append(vecs, req[i].Vector)
		contents = append(contents, string(req[i].Content))
		idx = append(idx, Sha256(req[i].Content))
	}

	opts := milvusclient.NewColumnBasedInsertOption(defaultCollection)
	opts = opts.WithFloatVectorColumn(vectorField, defaultDim, vecs)
	opts = opts.WithVarcharColumn(contentField, contents)
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

	ret, err := c.cli.Search(ctx, opts)
	if err != nil {
		return nil, err
	}

	var qresp []QueryResp

	for _, resp := range ret {
		col := resp.GetColumn(contentField)

		for j := range col.Len() {
			content, err := col.GetAsString(j)
			if err != nil {
				log.WithError(err).Error("error fetch search result")
				continue
			}
			qresp = append(qresp, QueryResp{
				Scores:  resp.Scores[j],
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
