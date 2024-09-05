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

func CreatePostComment(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var postCommentData models.PostComment
	if err := c.Bind(&postCommentData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with post comment schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}
	postCommentData.UserId = uri.UserId

	id, err := postCommentData.Insert(uuidString)
	if err != nil {
		logger.Logger.Error("API :: Error while inserting post comment", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": id,
	})

}

func UpdatePostComment(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var postCommentData models.PostComment
	if err := c.Bind(&postCommentData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with post comment schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}
	postCommentData.UserId = uri.UserId

	id, err := postCommentData.Update(uuidString, uri.PostCommentId)
	if err != nil {
		logger.Logger.Error("API :: Error while updating post comment", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": id,
	})

}

func DeletePostComment(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	status, err := models.DeletePostComment(uuidString, uri.PostCommentId, uri.UserId)
	if err != nil {
		logger.Logger.Error("API :: Error while delete post comment", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": status,
	})

}

func FetchPostComments(c *gin.Context) {
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

	var limit int
	var offset int

	limit, _ = strconv.Atoi(limitQueryStr)
	offset, _ = strconv.Atoi(offsetQueryStr)

	if limit > 50 {
		limit = 50
	} else if limit == 0 {
		limit = 10
	}

	postCommentData, count, err := models.FetchPostComments(uuidString, limit, offset)
	if err != nil {
		logger.Logger.Error("API :: Error while fetch post comments", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": postCommentData,
		"count":   count,
	})

}
