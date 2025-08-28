package embed

import (
	"context"

	"github.com/ollama/ollama/api"
	log "github.com/sirupsen/logrus"
)

const (
	defaultModel = "bge-m3:567m"
)

var client *api.Client

func init() {
	var err error

	client, err = api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
}

func Embed(words ...string) ([][]float32, error) {
	ctx := context.Background()
	req := &api.EmbedRequest{
		Model: defaultModel,
		Input: words,
	}

	resp, err := client.Embed(ctx, req)
	return resp.Embeddings, err
}
