package main

import (
	"context"
	"flag"
	"rag/milvus"

	log "github.com/sirupsen/logrus"
)

// 定义命令行 flag
var recreate = flag.Bool("recreate", false, "drop and recreate collections before migration")

func main() {
	log.SetLevel(log.InfoLevel)

	flag.Parse()

	ctx := context.Background()
	m := milvus.NewMigration(ctx)
	m.Migrate(*recreate)
}
