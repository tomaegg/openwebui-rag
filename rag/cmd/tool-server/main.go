package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	v1 "rag/generated/ragtools/v1"
	"rag/generated/ragtools/v1/handler"
	"rag/schemas"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.Info("server is to be started")
	gracefully()
}

const defaultAddr = ":8080"

func gracefully() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // 或 ["*"] 放行所有
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	h := handler.NewToolServer()
	v1.RegisterHandlersWithBaseURL(e, h, "")
	e.GET("/openapi.json", func(ctx echo.Context) error {
		j := json.RawMessage(schemas.ToolsOpenAPI)
		ctx.JSON(http.StatusOK, j)
		return nil
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		log.Info("server is running")
		if err := e.Start(defaultAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.Release(ctx); err != nil {
		log.WithError(err).Error("failed to release server resource")
	}
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	log.Println("server exited")
}
