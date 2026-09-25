package router

import (
	"context"
	"log/slog"
	"time"

	"github.com/artvault/artvault/internal/config"
	"github.com/artvault/artvault/internal/handler"
	"github.com/artvault/artvault/internal/middleware"
	"github.com/artvault/artvault/internal/repository"
	"github.com/artvault/artvault/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handlers 汇聚所有处理器。
type Handlers struct {
	Auth        *handler.AuthHandler
	Artwork     *handler.ArtworkHandler
	Exhibition  *handler.ExhibitionHandler
	Artist      *handler.ArtistHandler
	Interaction *handler.InteractionHandler
	Review      *handler.ReviewHandler
	Audit       *handler.AuditHandler
}

// NewRouter 装配路由、CORS、限流、健康检查与鉴权。
func NewRouter(cfg *config.Config, logger *slog.Logger, db *mongo.Database, hs *Handlers, auditSvc *service.AuditLogService, artworkRepo *repository.ArtworkRepository, exhibitionRepo *repository.ExhibitionRepository) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RateLimiter(cfg.RateLimitPerMin))
	r.Use(middleware.ErrorHandler(logger))

	healthz := func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) }
	readyz := func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.Client().Ping(ctx, nil); err != nil {
			c.JSON(503, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	}
	r.GET("/healthz", healthz)
	r.GET("/readyz", readyz)
	r.GET("/api/healthz", healthz)
	r.GET("/api/readyz", readyz)

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		authGroup.Use(middleware.RateLimiterWindow(cfg.LoginRateLimit, cfg.LoginRateWindow))
		{
			authGroup.POST("/login", hs.Auth.Login)
			authGroup.POST("/register", hs.Auth.Register)
		}

		public := v1.Group("")
		public.Use(middleware.ContentReview(artworkRepo, exhibitionRepo))
		{
			public.GET("/artworks", hs.Artwork.List)
			public.GET("/artworks/:id", hs.Artwork.Get)
			public.GET("/exhibitions", hs.Exhibition.List)
			public.GET("/exhibitions/:id", hs.Exhibition.Get)
			public.GET("/exhibitions/:id/readiness", hs.Exhibition.Readiness)
			public.GET("/artists", hs.Artist.List)
			public.GET("/artists/:id", hs.Artist.Get)
			public.GET("/interactions", hs.Interaction.List)
		}

		secured := v1.Group("")
		secured.Use(middleware.Auth(cfg))
		secured.Use(middleware.AuditLogger(auditSvc, logger))
		{
			secured.POST("/interactions", hs.Interaction.Create)
			secured.DELETE("/interactions", hs.Interaction.Cancel)

			artistOnly := secured.Group("")
			artistOnly.Use(middleware.RequireRole("Admin", "Artist"))
			{
				artistOnly.POST("/artworks", hs.Artwork.Create)
				artistOnly.PATCH("/artworks/:id/status", hs.Artwork.ChangeStatus)
				artistOnly.PATCH("/artists/me", hs.Artist.Upsert)
			}

			curatorOnly := secured.Group("")
			curatorOnly.Use(middleware.RequireRole("Admin", "Curator"))
			{
				curatorOnly.POST("/exhibitions", hs.Exhibition.Create)
				curatorOnly.PATCH("/exhibitions/:id/status", hs.Exhibition.ChangeStatus)
				curatorOnly.POST("/exhibitions/:id/artworks/:artworkId", hs.Exhibition.AddArtwork)
			}

			adminOnly := secured.Group("")
			adminOnly.Use(middleware.RequireRole("Admin", "Curator"))
			{
				adminOnly.PATCH("/artworks/:id/review", hs.Artwork.Review)
				adminOnly.POST("/reviews/:targetId", hs.Review.Create)
				adminOnly.GET("/reviews", hs.Review.List)
				adminOnly.GET("/audit-logs", hs.Audit.List)
			}
		}
	}
	return r
}
