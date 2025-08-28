package milvus

import (
	"context"
	"os"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	log "github.com/sirupsen/logrus"
)

func NewDefaultClient(ctx context.Context) (*milvusclient.Client) {
	addr := os.Getenv("MILVUS_ENDPOINT")
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: addr,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("created milvus client: address: %s", addr)
	return client
}
