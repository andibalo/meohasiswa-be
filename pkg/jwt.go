package pkg

import (
	"github.com/andibalo/meowhasiswa-be/internal/model"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/oops"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func GenerateToken(user *model.User) (tokenString string, err error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})

	tokenString, err = token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		log.Println(err)

		return "", oops.Code(response.ServerError.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusInternalServerError).Errorf("Failed to sign JWT")
	}

	return tokenString, nil
}
