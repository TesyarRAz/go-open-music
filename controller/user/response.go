package user

import (
	"strconv"

	"github.com/TesyarRAz/go-open-music/model"
	"github.com/gin-gonic/gin"
)

func storeResponse(user model.User) interface{} {
	return gin.H{
		"userId": strconv.Itoa(user.ID),
	}
}

func loginResponse(user model.User) interface{} {
	return gin.H{
		"accessToken":  user.AccessToken,
		"refreshToken": user.RefreshToken,
	}
}

func refreshResponse(user model.User) interface{} {
	return gin.H{
		"accessToken": user.AccessToken,
	}
}
