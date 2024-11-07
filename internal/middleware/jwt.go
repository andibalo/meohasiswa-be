package middleware

// TokenClaims : struct for validate token claims
//type TokenClaims struct {
//	ID       interface{} `json:"id"`
//	Sub      int         `json:"sub"`
//	Email    string      `json:"email"`
//	Name     string      `json:"name"`
//	UserName string      `json:"username"`
//	Token    string      `json:"token"`
//	jwt.RegisteredClaims
//}
//
//// contextClaimKey key value store/get token on context
//const ContextClaimKey = "ctx.mw.auth.claim"
//
//// JwtMiddleware : check jwt token header bearer scheme
//func JwtMiddleware(cfg *config.Config) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		ctx.Writer.Header().Set("Content-Type", "application/json")
//		secretKey := cfg.Token.JWTSecret
//		staticToken := cfg.Token.JWTStatic
//
//		// token claims
//		claims := &TokenClaims{}
//		headerToken, err := ParseTokenFromHeader(cfg, ctx)
//		if err != nil {
//			httpresp.HttpRespError(ctx, err)
//			return
//		}
//
//		if headerToken == staticToken {
//			ctx.Set(httpclient.XUserEmail, constant.EMAIL_ADMIN_EFISHERY)
//			ctx.Set(ContextClaimKey, &TokenClaims{
//				Email: constant.EMAIL_ADMIN_EFISHERY,
//				Name:  "superadmin",
//			})
//
//			ctx.Next()
//			return
//		}
//		token, err := jwt.ParseWithClaims(headerToken, claims, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // check signing method
//				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//			}
//			return []byte(secretKey), nil
//		})
//		// check parse token error
//		if err != nil {
//			httpresp.HttpRespError(ctx, apperr.NewWithCode(apperr.CodeHTTPUnauthorized, err.Error()))
//			return
//		}
//
//		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
//			claims.Token = headerToken
//			ctx.Set(httpclient.XUserEmail, claims.Email)
//			ctx.Set(ContextClaimKey, claims)
//			ctx.Next()
//		} else {
//			httpresp.HttpRespError(ctx, apperr.NewWithCode(apperr.CodeHTTPUnauthorized, err.Error()))
//			return
//		}
//	}
//}
//
//func ParseTokenFromHeader(cfg *config.Config, ctx *gin.Context) (string, error) {
//	var (
//		headerToken = ctx.Request.Header.Get("Authorization")
//		splitToken  []string
//	)
//
//	splitToken = strings.Split(headerToken, "Bearer ")
//
//	// check valid bearer token
//	if len(splitToken) <= 1 {
//		return "", apperr.NewWithCode(apperr.CodeHTTPUnauthorized, `Invalid Token`)
//	}
//
//	return splitToken[1], nil
//}
//
//func ParseToken(c *gin.Context) *TokenClaims {
//	log := cilog.Init(cilog.Options{})
//
//	v := c.Value(ContextClaimKey)
//	token := new(TokenClaims)
//	if v == nil {
//		return token
//	}
//	out, ok := v.(*TokenClaims)
//	if !ok {
//		return token
//	}
//	log.InfoWithContext(c.Request.Context(), "token parsed",
//		cilog.Field("user_email", out.Email),
//		cilog.Field("user_name", out.Name),
//	)
//
//	return out
//}
//
//func GetToken(c *gin.Context) string {
//	authorization := c.Request.Header.Get("Authorization")
//	tokens := strings.Split(authorization, "Bearer ")
//
//	return tokens[1]
//}
