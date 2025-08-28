package main

import (
	"context"
	"rag/milvus"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	ctx := context.Background()
	m := milvus.NewMigration(ctx)
	m.Migrate()
}
