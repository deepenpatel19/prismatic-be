package main

import (
	"net/http"
	"os"
	"time"

	// JWT
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/deepenpatel19/prismatic-be/middleware"

	// Gin
	"github.com/gin-gonic/gin"

	// Zap logger
	"go.uber.org/zap"

	// Internal packages
	"github.com/deepenpatel19/prismatic-be/api"
	"github.com/deepenpatel19/prismatic-be/core"
	"github.com/deepenpatel19/prismatic-be/logger"

	// "github.com/prismatic-be/middleware"
	"github.com/deepenpatel19/prismatic-be/models"
)

func Hello(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ok"})
}

func main() {
	core.ReadEnvFile()        // Configure ENV File
	logger.LoggerInit()       // Configure Logger
	models.RunMigrations()    // Run migrations to sync db schema related changes
	models.CreateConnection() // Create DB connection pool

	authMiddleware, err := middleware.GetAuthMiddleware()
	if err != nil {
		logger.Logger.Error("MAIN :: Error while configuring auth middleware", zap.Error(err))
		os.Exit(1)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Timeout(60*time.Second, middleware.NewServiceUnavailable()))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.POST("/user", api.CreateUser) // Open endpoint

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logger.Logger.Error("MAIN :: No route found", zap.Any("claims", claims))
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/token", authMiddleware.LoginHandler)
	r.GET("/logout", authMiddleware.LogoutHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	r.GET("/users", api.AllUser)

	// Auth Group
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())

	// User APIs
	auth.GET("/user/:userId/me", api.Me)
	auth.PUT("/user/:userId", api.UpdateUser)
	auth.DELETE("/user/:userId", api.DeleteUser)

	// Starting server
	if err := r.Run(":8000"); err != nil {
		logger.Logger.Fatal("Failed to start the server:", zap.Error(err))
		return
	}
}
