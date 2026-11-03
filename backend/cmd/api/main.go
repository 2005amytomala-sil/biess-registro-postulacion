package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/config"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/handler"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/middleware"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/jwt"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/mocks"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/postgres"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/platform/redis"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/repository"
	"github.com/2005amytomala-sil/biess-registro-postulacion/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using defaults")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set Gin mode
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize PostgreSQL
	pg, err := postgres.New(context.Background(), &cfg.Postgres)
	if err != nil {
		log.Fatalf("Error connecting to postgres: %v", err)
	}
	defer pg.Close()

	// Initialize Redis
	redisClient := redis.New(&cfg.Redis)
	defer redisClient.Close()

	// Initialize dependencies
	postulanteRepo := repository.NewPostulanteRepository(pg)
	registroCivil := mocks.NewRegistroCivilMock()
	intentosLimiter := redis.NewIntentosLimiter(redisClient.Client())
	tokenStorage := redis.NewTokenStorage(redisClient.Client())
	tokenManager := jwt.NewTokenManager(cfg.JWT.SecretKey, cfg.JWT.ExpirationHrs)
	authService := service.NewAuthService(postulanteRepo, registroCivil, intentosLimiter, tokenManager, tokenStorage)
	authHandler := handler.NewAuthHandler(authService)

	// Catalog dependencies
	catalogoRepo := repository.NewCatalogoRepository(pg.Pool())
	catalogoCache := redis.NewCatalogoCache(redisClient)
	catalogoService := service.NewCatalogoService(catalogoRepo, catalogoCache, log.Default())
	catalogoHandler := handler.NewCatalogoHandler(catalogoService)

	h := handler.NewHandler(pg, redisClient, authHandler, catalogoHandler)

	// Setup router
	r := gin.New()

	// Add custom middlewares
	r.Use(middleware.Recover())
	r.Use(middleware.Logging())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORS(cfg.CORS.AllowedOrigins))

	h.RegisterRoutes(r)

	// Server configuration
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
