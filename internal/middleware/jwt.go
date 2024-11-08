package middleware

import (
	"fmt"
	"github.com/andibalo/meowhasiswa-be/internal/config"
	"github.com/andibalo/meowhasiswa-be/internal/constants"
	"github.com/andibalo/meowhasiswa-be/internal/response"
	"github.com/andibalo/meowhasiswa-be/pkg/httpclient"
	"github.com/andibalo/meowhasiswa-be/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/samber/oops"
	"net/http"
	"strings"
)

// TokenClaims : struct for validate token claims
type TokenClaims struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	jwt.RegisteredClaims
}

// contextClaimKey key value store/get token on context
const ContextClaimKey = "ctx.mw.auth.claim"

// JwtMiddleware : check jwt token header bearer scheme
func JwtMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "application/json")
		secretKey := cfg.GetAuthCfg().JWTSecret
		staticToken := cfg.GetAuthCfg().JWTStaticToken

		// token claims
		claims := &TokenClaims{}
		headerToken, err := ParseTokenFromHeader(ctx)
		if err != nil {
			httpresp.HttpRespError(ctx, err)
			return
		}

		if headerToken == staticToken {
			ctx.Set(httpclient.XUserEmail, constants.EMAIL_ADMIN_MEOWHASISWA)
			ctx.Set(ContextClaimKey, &TokenClaims{
				Email:    constants.EMAIL_ADMIN_MEOWHASISWA,
				UserName: "superadmin",
				Role:     constants.ADMIN_ROLE,
			})

			ctx.Next()
			return
		}

		token, err := jwt.ParseWithClaims(headerToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // check signing method
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		// check parse token error
		if err != nil {
			httpresp.HttpRespError(ctx, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(err.Error()))
			return
		}

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			claims.Token = headerToken
			ctx.Set(httpclient.XUserEmail, claims.Email)
			ctx.Set(ContextClaimKey, claims)
			ctx.Next()
		} else {
			httpresp.HttpRespError(ctx, oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf(err.Error()))
			return
		}
	}
}

func ParseTokenFromHeader(ctx *gin.Context) (string, error) {
	var (
		headerToken = ctx.Request.Header.Get("Authorization")
		splitToken  []string
	)

	splitToken = strings.Split(headerToken, "Bearer ")

	// check valid bearer token
	if len(splitToken) <= 1 {
		return "", oops.Code(response.Unauthorized.AsString()).With(httpresp.StatusCodeCtxKey, http.StatusUnauthorized).Errorf("Invalid Token")
	}

	return splitToken[1], nil
}

func ParseToken(c *gin.Context) *TokenClaims {

	v := c.Value(ContextClaimKey)
	token := new(TokenClaims)
	if v == nil {
		return token
	}
	out, ok := v.(*TokenClaims)
	if !ok {
		return token
	}

	return out
}

func GetToken(c *gin.Context) string {
	authorization := c.Request.Header.Get("Authorization")
	tokens := strings.Split(authorization, "Bearer ")

	return tokens[1]
}
