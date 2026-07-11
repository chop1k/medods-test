package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/chop1k/medods-test/internal/app/config"
	"github.com/chop1k/medods-test/internal/database"
	"github.com/chop1k/medods-test/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	backend *http.Server
}

func NewServer(cfg *config.ServerConfig, db *sql.DB) *HttpServer {
	router := gin.Default()

	templateStorage := database.NewTemplateStorage(db)

	v1 := router.Group("/v1")

	RegisterTemplateRoutes(v1, handler.NewTemplateHandler(templateStorage))
	RegisterTaskRoutes(v1, handler.NewTaskHandler())
	RegisterTagRoutes(v1, handler.NewTagHandler())
	RegisterSchedulingRoutes(v1, handler.NewSchedulingHandler(templateStorage))

	return &HttpServer{
		backend: &http.Server{
			Addr:         cfg.Addr(),
			Handler:      router,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *HttpServer) Listen() error {
	if err := s.backend.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.backend.WriteTimeout)
	defer cancel()

	return s.backend.Shutdown(shutdownCtx)
}
