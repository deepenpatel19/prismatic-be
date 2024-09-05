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
	r.GET("/", Hello)
	r.POST("/user", api.CreateUser) // Open endpoint

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		logger.Logger.Error("MAIN :: No route found", zap.Any("claims", claims))
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/token", authMiddleware.LoginHandler)
	r.GET("/logout", authMiddleware.LogoutHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	// Auth Group
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())

	// User APIs
	auth.GET("/user/:userId/me", api.Me)
	auth.PUT("/user/:userId", api.UpdateUser)
	auth.DELETE("/user/:userId", api.DeleteUser)

	// add/remove connection APIs
	auth.POST("/user/:userId/addConnection", api.AddConnection)
	auth.POST("/user/:userId/removeConnection", api.RemoveConnection)

	// Post APIs
	auth.GET("/user/:userId/posts", api.FetchPosts)
	auth.POST("/user/:userId/post", api.CreatePost)
	auth.PUT("/user/:userId/post/:postId", api.UpdatePost)
	auth.DELETE("/user/:userId/post/:postId", api.DeletePost)

	// Post comment APIs
	auth.GET("/user/:userId/post/:postId/comments", api.FetchPostComments)
	auth.POST("/user/:userId/post/:postId/comment", api.CreatePostComment)
	auth.PUT("/user/:userId/post/:postId/comment/:postCommentId", api.UpdatePostComment)
	auth.DELETE("/user/:userId/post/:postId/comment/:postCommentId", api.DeletePostComment)

	// Starting server
	if err := r.Run(":8000"); err != nil {
		logger.Logger.Fatal("Failed to start the server:", zap.Error(err))
		return
	}
}
