package embed

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/ollama/ollama/api"
	log "github.com/sirupsen/logrus"
)

const (
	defaultModel   = "bge-m3:567m"
	defaultBaseURL = "http://localhost:11434"
)

var client *api.Client

func init() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	client, err = api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
}

func Embed(words []string) ([][]float32, error) {
	ctx := context.Background()
	req := &api.EmbedRequest{
		Model: defaultModel,
		Input: words,
	}

	resp, err := client.Embed(ctx, req)
	return resp.Embeddings, err
}
