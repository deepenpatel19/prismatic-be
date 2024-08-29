package middleware

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/deepenpatel19/prismatic-be/core"
	"github.com/deepenpatel19/prismatic-be/logger"

	"github.com/deepenpatel19/prismatic-be/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	identityKey = "id"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*models.UserSchema); ok {
			return jwt.MapClaims{
				identityKey: v.Email,
			}
		}
		return jwt.MapClaims{}
	}
}

func identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &models.UserSchema{
			Email: claims[identityKey].(string),
		}
	}
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals login
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userEmail := loginVals.Username
		password := loginVals.Password

		// userDataFromDb := models.FetchUserForAuth(userEmail)
		userDataFromDb := models.FetchUserForAuthV1(userEmail)
		if userDataFromDb.Id == 0 {
			return "", errors.New("no account found")
		}

		err := bcrypt.CompareHashAndPassword([]byte(userDataFromDb.Password), []byte(password))
		if err != nil {
			logger.Logger.Error("AUTH :: Error while comparing hash of password", zap.Error(err))
		} else {
			logger.Logger.Info("AUTH :: user data ", zap.Any("data", userDataFromDb))
			return &userDataFromDb, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*models.UserSchema); ok {
			// userDataFromDb := models.FetchUserForAuth(v.Email)
			userDataFromDb := models.FetchUserForAuthV1(v.Email)
			logger.Logger.Info("AUTH :: authorizator ", zap.Any("data", userDataFromDb))
			return userDataFromDb.Id != 0
		}
		return false
	}
}

func unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       core.Config.AuthRealm,
		Key:         []byte(core.Config.AuthSecretKey),
		Timeout:     12 * time.Hour,
		MaxRefresh:  12 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: payloadFunc(),

		IdentityHandler: identityHandler(),
		Authenticator:   authenticator(),
		Authorizator:    authorizator(),
		Unauthorized:    unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	return authMiddleware, err
}
