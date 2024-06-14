package main

import (
	"context"
	"geogracom-test/config"
	"geogracom-test/internal/api"
	"geogracom-test/internal/route"
	"geogracom-test/pkg/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// ---------------- may fail ----------------
	var (
		cfg = config.GetConfig()
		db  = repository.New(ctx, cfg.PostgresConnection)
	)

	// --------------- can't fail ---------------
	var (
		routeApi = route.NewMethod(route.NewRepo(db))
		srv      = api.New()
	)

	routeApi.MapHandlers(srv.App.Group("api/route"))

	srv.Start(ctx, cfg.ServerAddress)
}
