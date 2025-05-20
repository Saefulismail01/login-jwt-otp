package delivery

import (
	"context"
	"database/sql"
	"fmt"
	"login-jwt-otp/config"
	"login-jwt-otp/delivery/controller"
	"login-jwt-otp/middleware"
	"login-jwt-otp/repository"
	"login-jwt-otp/usecase"
	"login-jwt-otp/utils/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Server struct {
	userRepo       repository.UserRepoInterface
	userUsecase    usecase.UserUsecase
	authUsecase    *usecase.AuthUsecase
	jwtService     service.JwtService
	authMiddleware *middleware.AuthMiddleware
	engine         *gin.Engine
	host           string
	db             *sql.DB
	server         *http.Server
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api")

	// Auth routes
	authGroup := rg.Group("/auth")
	controller.NewAuthController(authGroup, s.jwtService, *s.authUsecase, s.userUsecase).Route()
}

func (s *Server) Run() {

	s.initRoute()

	s.server = &http.Server{
		Addr:    s.host,
		Handler: s.engine,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server running on %s\n", s.host)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start server: %v", err))
		}
	}()

	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	if err := s.db.Close(); err != nil {
		fmt.Printf("Error closing database: %v\n", err)
	}

	fmt.Println("Server gracefully stopped 󱠡")

}

func NewServer() *Server {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	db, cfg, err := config.ConnectDB()
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return nil
	}

	// Initialize database repositories
	userRepo := repository.NewUserRepo(db)

	// Initialize service dependencies
	jwtSecret := os.Getenv("JWT_SECRET")
	appName := os.Getenv("APP_NAME")
	jwtExpiryStr := os.Getenv("JWT_EXPIRY")
	jwtExpiry, err := time.ParseDuration(jwtExpiryStr)
	if err != nil {
		jwtExpiry = 24 * time.Hour
		fmt.Printf("Warning: Could not parse JWT_EXPIRY value '%s', using default of 24h: %v\n", jwtExpiryStr, err)
	}

	jwtService := service.NewJwtService(jwtSecret, appName, jwtExpiry)

	// Initialize middleware components
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Initialize usecase layer components
	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)

	// Initialize Gin engine
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		userRepo:       userRepo,
		userUsecase:    userUsecase,
		authUsecase:    authUsecase,
		jwtService:     jwtService,
		authMiddleware: authMiddleware,
		engine:         engine,
		host:           host,
		db:             db,
	}
}
