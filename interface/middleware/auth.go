package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTClaims は、JWTトークンに含まれるクレーム情報を定義します
type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthMiddleware は、JWTトークンによる認証を行うミドルウェアです
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーからトークンを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Bearerトークンの形式をチェック
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			tokenString := tokenParts[1]

			// JWTトークンを検証
			token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				// 署名アルゴリズムの検証
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// クレームの取得
			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				// コンテキストにユーザーIDを設定
				c.Set("user_id", claims.UserID)
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token claims",
			})
		}
	}
}

// GetUserIDFromContext は、コンテキストからユーザーIDを取得するヘルパー関数です
func GetUserIDFromContext(c echo.Context) uint {
	userID := c.Get("user_id")
	switch v := userID.(type) {
	case uint:
		return v
	case float64:
		return uint(v)
	default:
		return 0
	}
}
