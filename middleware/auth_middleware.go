package middleware

import (
	"net/http"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/TesyarRAz/go-open-music/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	Db *gorm.DB
}

func (a *AuthMiddleware) Auth(c *gin.Context) {
	token, err := service.ValidateAuthorization(c.GetHeader("Authorization"))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "unauthorized",
			"error":   err.Error(),
		})

		return
	}

	if token == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "unauthorized",
			"error":   err.Error(),
		})
	}

	userId, ok := token.Get("userId")

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "unauthorized",
			"error":   err.Error(),
		})

		return
	}

	var user model.User

	if err := a.Db.First(&user, "id = ?", userId).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.Set("user", &user)
}
