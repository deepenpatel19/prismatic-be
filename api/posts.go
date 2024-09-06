package api

import (
	"strconv"

	"github.com/deepenpatel19/prismatic-be/logger"
	"github.com/deepenpatel19/prismatic-be/models"
	"github.com/deepenpatel19/prismatic-be/schemas"
	"github.com/deepenpatel19/prismatic-be/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePost(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var postData models.Post
	if err := c.Bind(&postData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with post schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}
	postData.UserId = uri.UserId

	id, err := postData.Insert(uuidString)
	if err != nil {
		logger.Logger.Error("API :: Error while inserting post", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": id,
	})
}

func UpdatePost(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var postData models.Post
	if err := c.Bind(&postData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with post schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}
	postData.UserId = uri.UserId

	id, err := postData.Update(uuidString, uri.PostId)
	if err != nil {
		logger.Logger.Error("API :: Error while updating post", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": id,
	})

}

func DeletePost(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	status, err := models.DeletePost(uuidString, uri.PostId, uri.UserId)
	if err != nil {
		logger.Logger.Error("API :: Error while delete post", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": status,
	})
}

func FetchPosts(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	limitQueryStr := c.DefaultQuery("limit", "10")
	offsetQueryStr := c.DefaultQuery("offset", "0")
	queryVersion := c.DefaultQuery("query_version", "v0")

	var limit int
	var offset int
	var postData []models.PostResponse
	var count int64
	var err error

	limit, _ = strconv.Atoi(limitQueryStr)
	offset, _ = strconv.Atoi(offsetQueryStr)

	if limit > 50 {
		limit = 50
	} else if limit == 0 {
		limit = 10
	}

	if queryVersion == "v1" {
		postData, count, err = models.FetchPostsV1(uuidString, limit, offset)
	} else {
		postData, count, err = models.FetchPosts(uuidString, limit, offset)
	}

	if err != nil {
		logger.Logger.Error("API :: Error while fetch posts", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": postData,
		"count":   count,
	})

}
