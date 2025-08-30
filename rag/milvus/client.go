package milvus

import (
	"context"
	"os"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
	log "github.com/sirupsen/logrus"
)

func NewDefaultClient(ctx context.Context) *milvusclient.Client {
	addr := os.Getenv("MILVUS_ENDPOINT")
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: addr,
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := client.UseDatabase(ctx, milvusclient.NewUseDatabaseOption(defaultDB)); err != nil {
		log.Fatal(err)
	}
	log.WithField("db", defaultDB).WithField("address", addr).Info("created milvus client")
	return client
}
