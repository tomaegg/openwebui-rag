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
	defaultMilvusAddr = "localhost:19530"
	defaultDB         = "openwebui_rag"
	defaultCollection = "rag_passage"
	defaultDim        = 1024
	defaultChunk      = 2048
)

const (
	vectorField  = "passage_vector"
	idField      = "passage_id"
	contentField = "content"
)

type Migration struct {
	cli *milvusclient.Client
}

func NewMigration(ctx context.Context) *Migration {
	return &Migration{cli: NewDefaultClient(ctx)}
}

func (m *Migration) Close() error {
	ctx := context.TODO()
	return m.cli.Close(ctx)
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

func (m *Migration) CreateCollection(ctx context.Context, name string, recreate bool) {
	var err error
	schema := entity.NewSchema()
	schema = schema.WithDynamicFieldEnabled(true)

	// NOTE: fields idx: sha256 64 bytes
	schema = schema.WithField(entity.NewField().WithName(idField).WithIsAutoID(false).WithDataType(entity.FieldTypeVarChar).WithMaxLength(64).
		WithIsPrimaryKey(true))

	schema = schema.WithField(entity.NewField().WithName(vectorField).WithDataType(entity.FieldTypeFloatVector).WithDim(defaultDim))

	schema = schema.WithField(entity.NewField().WithName(contentField).WithDataType(entity.FieldTypeVarChar).WithMaxLength(defaultChunk))

	indexOptions := []milvusclient.CreateIndexOption{
		milvusclient.NewCreateIndexOption(name,
			vectorField,
			index.NewAutoIndex(entity.COSINE),
		),
	}

	if recreate {
		m.DropCollection(ctx, name)
		log.WithField("name", name).Warn("old collection dropped")
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

func (m *Migration) Migrate(recreate bool) {
	ctx := context.Background()
	m.CreateDB(ctx, defaultDB)
	m.CreateCollection(ctx, defaultCollection, recreate)
}
