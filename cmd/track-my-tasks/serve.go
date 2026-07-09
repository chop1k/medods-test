package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/chop1k/medods-test/internal/config"
	"github.com/chop1k/medods-test/internal/database"
	"github.com/chop1k/medods-test/internal/handlers"
	"github.com/chop1k/medods-test/internal/routes"
)

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	serverCfg := config.RegisterServerFlags(fs)
	dbCfg := config.RegisterDatabaseFlags(fs)

	if err := fs.Parse(args); err != nil {
		fatal("serve: failed to parse flags: %v", err)
	}

	db, err := database.Connect(*dbCfg)

	if err != nil {
		fatal("serve: failed to connect to db: %v", err)
	}

	startServer(*serverCfg, db)
}

// startServer builds the router and blocks until the server exits or a
// shutdown signal is received.
func startServer(serverCfg config.ServerConfig, db *sql.DB) {
	router := newRouter(db)

	srv := &http.Server{
		Addr:         serverCfg.Addr(),
		Handler:      router,
		ReadTimeout:  serverCfg.ReadTimeout,
		WriteTimeout: serverCfg.WriteTimeout,
		IdleTimeout:  serverCfg.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("http server listening on %s", serverCfg.Addr())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fatal("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down http server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), serverCfg.WriteTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fatal("http server shutdown error: %v", err)
	}
}

// newRouter wires up the application's Gin router. Kept separate from
// startServer so it can also be reused directly by tests.
func newRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	templateStorage := database.NewTemplateStorage(db)

	v1 := router.Group("/v1")
	routes.RegisterTemplateRoutes(v1, handlers.NewTemplateHandler(templateStorage))
	routes.RegisterTaskRoutes(v1, handlers.NewTaskHandler())
	routes.RegisterTagRoutes(v1, handlers.NewTagHandler())
	routes.RegisterSchedulingRoutes(v1, handlers.NewSchedulingHandler(templateStorage))

	return router
}
