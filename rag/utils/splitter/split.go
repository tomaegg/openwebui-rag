package splitter

import (
	"context"
	"os"
	"rag/generated/split"
)

var cli *split.APIClient

func init() {
	cfg := split.NewConfiguration()
	cfg.Servers[0].URL = os.Getenv("SPLIT_HOST")
	cli = split.NewAPIClient(cfg)
}

func SplitContent(content string) ([]string, error) {
	api := cli.DefaultAPI
	ctx := context.TODO()
	req := split.NewSplitRequest([]string{content})
	splitResp, _, err := api.SplitTextSplitPost(ctx).SplitRequest(*req).Execute()
	if err != nil {
		return nil, err
	}
	return splitResp.Chunks, nil
}
