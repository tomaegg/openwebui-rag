package handler

import (
	"context"
	"errors"
	"net/http"
	v1 "rag/generated/ragtools/v1"
	"rag/milvus"
	"rag/utils/embed"

	"github.com/labstack/echo/v4"
	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type ToolServer struct {
	cli *milvus.RagCli
}

var _ v1.ServerInterface = &ToolServer{}

func NewToolServer() *ToolServer {
	return &ToolServer{
		cli: milvus.NewRagCli(context.TODO()),
	}
}

func (s *ToolServer) Release(ctx context.Context) error {
	return s.cli.Close(ctx)
}

func (s *ToolServer) SearchByTopic(ctx echo.Context) error {
	var req v1.SearchRequest
	err := ctx.Bind(&req)
	if err != nil {
		return err
	}
	if len(req.Topic) == 0 {
		return errors.New("cannot search empty string")
	}

	vec, err := embed.Embed(req.Topic)
	if err != nil {
		log.WithError(err).Error("error encode topic")
		return err
	}

	qvec := lo.Map(vec, func(item []float32, _ int) entity.Vector {
		return entity.FloatVector(item)
	})

	resp, err := s.cli.Search(context.TODO(), 10, qvec...)
	if err != nil {
		log.WithError(err).Error("error query topic")
		return err
	}

	searchResp := v1.SearchResponse{
		Count: len(resp),
		Results: lo.Map(resp, func(item milvus.QueryResp, _ int) v1.SearchResultItem {
			return v1.SearchResultItem{
				Score: item.Score,
				Text:  item.Content,
			}
		}),
	}

	log.WithField("topic", req.Topic).Info("search topic done")

	ctx.JSON(http.StatusOK, searchResp)

	return nil
}
