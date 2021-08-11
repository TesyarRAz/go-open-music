package user

import (
	"net/http"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/TesyarRAz/go-open-music/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	Db *gorm.DB
}

func (u *UserController) Store(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	if err := u.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil membuat user",
		"data":    storeResponse(user),
	})
}

func (u *UserController) Login(c *gin.Context) {
	var (
		request LoginRequest
		user    model.User
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	if err := u.Db.First(&user, "username = ? AND password = ?", request.Username, request.Password).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	accessToken, refreshToken, err := service.CreateToken(user)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "fail",
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
		return
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	if err := u.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"data":    loginResponse(user),
		"message": "Berhasil Login",
	})
}

func (u *UserController) Refresh(c *gin.Context) {
	var (
		request RefreshRequest
		user    model.User
	)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	token, err := service.ValidateToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "unauthorized",
			"error":   err.Error(),
		})
		return
	}

	userId, ok := token.Get("userId")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "unauthorized",
		})

		return
	}

	if err := u.Db.First(&user, "id = ?", userId).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "unauthorized",
		})
		return
	}

	accessToken, err := service.CreateAccessToken(user)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "fail",
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
		return
	}

	user.AccessToken = accessToken
	if err := u.Db.Model(&user).Update("accessToken = ?", user.AccessToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    refreshResponse(user),
		"message": "Berhasil refresh token",
	})
}
