package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/artvault/artvault/database/seeds"
	"github.com/artvault/artvault/internal/config"
	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/database"
	"github.com/artvault/artvault/internal/handler"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/router"
	"github.com/artvault/artvault/internal/service"
	"github.com/artvault/artvault/internal/util"
)

func main() {
	cfg := config.Load()
	log := util.NewLogger(cfg.AppEnv)
	if err := cfg.Validate(); err != nil {
		log.Error("invalid config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoConnectTimeout)
	defer cancel()
	db, client, err := database.Connect(ctx, cfg.MongoURI, cfg.MongoDBName, cfg.MongoConnectTimeout)
	if err != nil {
		log.Error("connect mongodb failed", "error", err)
		os.Exit(1)
	}
	log.Info(constants.LogMongoConnected, "db", cfg.MongoDBName)
	seeds.Seed(ctx, db, log)

	userRepo := repository.NewUserRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	artworkRepo := repository.NewArtworkRepository(db)
	exhibitionRepo := repository.NewExhibitionRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	authSvc := service.NewAuthService(userRepo, log, cfg)
	artistSvc := service.NewArtistService(artistRepo, log)
	artworkSvc := service.NewArtworkService(artworkRepo, log)
	exhibitionSvc := service.NewExhibitionService(exhibitionRepo, log)
	interactionSvc := service.NewInteractionService(interactionRepo, artworkRepo, log)
	reviewSvc := service.NewReviewService(reviewRepo, log)
	auditSvc := service.NewAuditLogService(auditRepo, log)

	hs := &router.Handlers{
		Auth:        handler.NewAuthHandler(authSvc),
		Artwork:     handler.NewArtworkHandler(artworkSvc),
		Exhibition:  handler.NewExhibitionHandler(exhibitionSvc),
		Artist:      handler.NewArtistHandler(artistSvc),
		Interaction: handler.NewInteractionHandler(interactionSvc),
		Review:      handler.NewReviewHandler(reviewSvc),
		Audit:       handler.NewAuditHandler(auditSvc),
	}

	ginRouter := router.NewRouter(cfg, log, db, hs, auditSvc, artworkRepo, exhibitionRepo)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      ginRouter,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info(constants.LogServerStart, "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownWait)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server forced to shutdown", "error", err)
	}
	_ = client.Disconnect(context.Background())
	log.Info("server exited")
}
