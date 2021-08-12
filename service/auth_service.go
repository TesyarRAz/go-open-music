package service

import (
	"strings"
	"time"

	"github.com/TesyarRAz/go-open-music/config"
	"github.com/TesyarRAz/go-open-music/model"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

const (
	ACCESS_TOKEN_EXP  = time.Minute * 30   // 30 Menit
	REFRESH_TOKEN_EXP = time.Hour * 24 * 7 // 1 Minggu
)

func CreateAccessToken(user model.User) (string, error) {
	token := jwt.New()

	token.Set("authorized", true)
	token.Set("userId", user.ID)
	token.Set(jwt.ExpirationKey, time.Now().Add(ACCESS_TOKEN_EXP).Unix())

	// serialized, err := jwt.Sign(token, jwa.HS256, )

	serialized, err := jwt.Sign(token, jwa.SignatureAlgorithm(config.AppConfig.JWT_ENCRYPT), config.AppConfig.JWT_SECRET)

	if err != nil {
		return "", err
	}

	return string(serialized), err
}

func CreateRefreshToken(user model.User) (string, error) {
	token := jwt.New()

	token.Set("userId", user.ID)
	token.Set(jwt.ExpirationKey, time.Now().Add(REFRESH_TOKEN_EXP).Unix())

	// serialized, err := jwt.Sign(token, jwa.HS256, )

	serialized, err := jwt.Sign(token, jwa.SignatureAlgorithm(config.AppConfig.JWT_ENCRYPT), config.AppConfig.JWT_SECRET)

	if err != nil {
		return "", err
	}

	return string(serialized), err
}

func CreateToken(user model.User) (string, string, error) {
	accessToken, err := CreateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := CreateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func ValidateAuthorization(authorization string) (jwt.Token, error) {
	token := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer"))

	return ValidateToken(token)
}

func ValidateToken(token string) (jwt.Token, error) {
	return jwt.ParseString(token,
		jwt.WithVerify(jwa.SignatureAlgorithm(config.AppConfig.JWT_ENCRYPT), config.AppConfig.JWT_SECRET),
		jwt.WithValidate(true),
		jwt.WithAcceptableSkew(time.Duration(time.Now().Unix())),
	)
}
