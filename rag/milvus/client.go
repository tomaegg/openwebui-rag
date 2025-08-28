package milvus

import (
	"context"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	log "github.com/sirupsen/logrus"
)

func NewDefaultClient(ctx context.Context) (*milvusclient.Client) {
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: defaultMilvusAddr,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("created milvus client: address: %s", defaultMilvusAddr)
	return client
}
