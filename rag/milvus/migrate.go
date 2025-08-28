package milvus

import (
	"context"
	"slices"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/index"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	log "github.com/sirupsen/logrus"
)

const (
	milvusAddr        = "localhost:19530"
	defaultDB         = "openwebui_rag"
	defaultCollection = "rag_passage"
	defaultDim        = 1024
	defaultChunk      = 1024
)

type Migration struct {
	cli *milvusclient.Client
}

func NewMigration(ctx context.Context) *Migration {
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: milvusAddr,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("created milvus client: addr: %s", milvusAddr)
	return &Migration{cli: client}
}

func (m *Migration) ListDB(ctx context.Context) []string {
	// List all existing databases
	databases, err := m.cli.ListDatabase(ctx, milvusclient.NewListDatabaseOption())
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("existed databases: %+v", databases)
	return databases
}

func (m *Migration) CreateDB(ctx context.Context, dbName string) {
	dbs := m.ListDB(ctx)
	if slices.Contains(dbs, dbName) {
		log.Infof("%s already exists, skipped", dbName)
		return
	}
	// create db
	err := m.cli.CreateDatabase(ctx, milvusclient.NewCreateDatabaseOption(dbName))
	if err != nil {
		log.Fatal(err)
	}

	db, err := m.cli.DescribeDatabase(ctx, milvusclient.NewDescribeDatabaseOption(dbName))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%+v", *db)
}

func (m *Migration) DropCollection(ctx context.Context, name string) {
	// NOTE: 暂时丢弃并且重新创建
	err := m.cli.DropCollection(ctx, milvusclient.NewDropCollectionOption(name))
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Migration) CreateCollection(ctx context.Context, name string) {
	var err error
	const (
		vectorField  = "passage_vector"
		idField      = "passage_id"
		contentField = "content"
	)
	schema := entity.NewSchema()
	schema = schema.WithDynamicFieldEnabled(true)

	// fields
	schema = schema.WithField(entity.NewField().WithName(idField).WithIsAutoID(false).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true))
	schema = schema.WithField(entity.NewField().WithName(vectorField).WithDataType(entity.FieldTypeFloatVector).WithDim(defaultDim))
	schema = schema.WithField(entity.NewField().WithName(contentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(defaultChunk))

	indexOptions := []milvusclient.CreateIndexOption{
		milvusclient.NewCreateIndexOption(name,
			vectorField,
			index.NewAutoIndex(entity.COSINE),
		),
	}

	// collection
	err = m.cli.CreateCollection(
		ctx,
		milvusclient.NewCreateCollectionOption(name, schema).WithIndexOptions(indexOptions...))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("collection created: %v", name)
}

func (m *Migration) Migrate() {
	ctx := context.Background()
	defer m.cli.Close(ctx)

	m.CreateDB(ctx, defaultDB)
	m.CreateCollection(ctx, defaultCollection)
}
