package api

import (
	"github.com/deepenpatel19/prismatic-be/core"
	"github.com/deepenpatel19/prismatic-be/logger"
	"github.com/deepenpatel19/prismatic-be/models"
	"github.com/deepenpatel19/prismatic-be/schemas"
	"github.com/deepenpatel19/prismatic-be/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var userData models.UserCreateSchema
	if err := c.Bind(&userData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with user create schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(userData.Password), core.Config.PasswordHashCost)
	if err != nil {
		logger.Logger.Error("API :: Error while generating password hash", zap.String("requestId", uuidString), zap.Error(err))
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}
	userData.Password = string(passwordHashBytes)

	id, err := userData.Insert(uuidString)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": id,
	})
}

func Me(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	data := models.FetchUserForMeV1(uuidString, uri.UserId)
	c.JSON(200, gin.H{
		"message": data,
	})
}

func UpdateUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	var userData models.UserCreateSchema
	if err := c.Bind(&userData); err != nil {
		logger.Logger.Error("API :: Error while binding request data with user update schema.",
			zap.String("requestId", uuidString),
			zap.Error(err),
		)
		c.JSON(400, gin.H{
			"message": "something went wrong - please check request body",
		})
		return
	}

	id, err := userData.Update(uuidString, uri.UserId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": id,
	})
}

func DeleteUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	var uri schemas.URI
	if err := c.ShouldBindUri(&uri); err != nil {
		logger.Logger.Error("API :: Error while uri binding", zap.Error(err), zap.String("requestId", uuidString))
		c.JSON(400, gin.H{"message": err})
		return
	}

	status, err := models.DeleteUserFromDB(uuidString, uri.UserId)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": status,
	})
}

func AllUser(c *gin.Context) {
	uuidString := utils.GetUUID()
	c.Header("X-REQUEST-ID", uuidString)

	data, count, err := models.FetchAllUsers(uuidString)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": data,
		"count":   count,
	})
}
