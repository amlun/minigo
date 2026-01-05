package http

import (
	"minigo/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"

	appsvc "minigo/internal/application/service"
	configx "minigo/internal/infrastructure/config"
	infrarepo "minigo/internal/infrastructure/repository"
	"minigo/internal/infrastructure/tx"
	"minigo/internal/interfaces/middleware"
)

// BuildRouter builds the gin engine with routes and middleware.
func BuildRouter(db *bun.DB) *gin.Engine {
	// 设置Gin模式
	if configx.IsDevEnv() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()

	// CORS middleware first to handle preflight requests
	engine.Use(middleware.CORSMiddleware())
	// global error handler
	engine.Use(middleware.ErrorHandlerMiddleware())
	engine.Use(gin.Recovery())
	// logging and performance middlewares per design doc
	engine.Use(middleware.RequestLoggerMiddleware())

	// repositories
	userRepo := infrarepo.NewBunUserRepository(db)

	// transaction manager
	txManager := tx.NewManager(db)

	// services
	authSvc := appsvc.NewAuthService(userRepo)
	userSvc := appsvc.NewUserService(userRepo, txManager)

	// infrastructure services
	//ossService := oss.NewOSSService()

	// handlers
	authHandler := handlers.NewAuthHandler(authSvc, userSvc)
	// health
	engine.GET("/api/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// admin routes (no shop context)
	apiGroup := engine.Group("/api")
	{
		apiGroup.POST("/auth/login", authHandler.Login)
	}

	return engine
}
