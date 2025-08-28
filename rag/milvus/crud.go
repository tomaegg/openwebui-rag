package milvus

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type RagCli struct {
	cli *milvusclient.Client
}

type VecIdx = []float32

type InsertReq struct {
	Vector  VecIdx
	Content []byte
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

func (c *RagCli) Insert(ctx context.Context, req ...InsertReq) (int64, error) {
	n := len(req)
	vecs := make([]VecIdx, 0, n)
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

	ret, err := c.cli.Insert(ctx, opts)

	return ret.InsertCount, err
}

func (c *RagCli) Query(ctx context.Context) {
}
