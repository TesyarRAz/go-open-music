package user

import (
	"strconv"

	"github.com/TesyarRAz/go-open-music/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

func storeResponse(user model.User) interface{} {
	return gin.H{
		"userId": strconv.Itoa(user.ID),
	}
}

func loginResponse(user model.User) interface{} {
	return gin.H{
		"accessToken":  user.AccessToken.String,
		"refreshToken": user.RefreshToken.String,
	}
}

func refreshResponse(user model.User) interface{} {
	return gin.H{
		"accessToken": user.AccessToken.String,
	}
}
