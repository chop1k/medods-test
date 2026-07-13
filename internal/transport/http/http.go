package http

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/chop1k/medods-test/internal/app/config"
	"github.com/chop1k/medods-test/internal/repository"
	"github.com/chop1k/medods-test/internal/transport/http/handler"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	backend *http.Server
}

func NewRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	templateStorage := repository.NewTemplateStorage(db)
	taskStorage := repository.NewTaskStorage(db)
	tagStorage := repository.NewTagStorage(db)

	v1 := router.Group("/v1")

	RegisterTemplateRoutes(v1, handler.NewTemplateHandler(templateStorage))
	RegisterTaskRoutes(v1, handler.NewTaskHandler(taskStorage))
	RegisterTagRoutes(v1, handler.NewTagHandler(tagStorage))
	RegisterSchedulingRoutes(v1, handler.NewSchedulingHandler(templateStorage, taskStorage))

	return router
}

func NewServer(cfg *config.ServerConfig, db *sql.DB) *HttpServer {
	return &HttpServer{
		backend: &http.Server{
			Addr:         cfg.Addr(),
			Handler:      NewRouter(db),
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
