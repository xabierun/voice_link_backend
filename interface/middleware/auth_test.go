package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// JWT_SECRETの設定
	os.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   map[string]interface{}
		shouldSetUser  bool
		expectedUserID uint
	}{
		{
			name:           "認証ヘッダーなし",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header is required",
			},
			shouldSetUser:  false,
			expectedUserID: 0,
		},
		{
			name:           "無効な認証ヘッダー形式",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid authorization header format",
			},
			shouldSetUser:  false,
			expectedUserID: 0,
		},
		{
			name:           "無効なトークン",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token",
			},
			shouldSetUser:  false,
			expectedUserID: 0,
		},
		{
			name:           "有効なトークン",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			shouldSetUser:  true,
			expectedUserID: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoの設定
			e := echo.New()

			// 有効なトークンの場合は事前に生成
			if tt.shouldSetUser {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": tt.expectedUserID,
					"exp":     time.Now().Add(time.Hour).Unix(),
					"iat":     time.Now().Unix(),
				})
				tokenString, _ := token.SignedString([]byte("test-secret"))
				tt.authHeader = "Bearer " + tokenString
			}

			// テスト用のハンドラー
			handler := func(c echo.Context) error {
				if tt.shouldSetUser {
					userID := GetUserIDFromContext(c)
					assert.Equal(t, tt.expectedUserID, userID)
				}
				return c.String(http.StatusOK, "success")
			}

			// ミドルウェアの適用
			middleware := AuthMiddleware()
			handlerWithMiddleware := middleware(handler)

			// リクエストの作成
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// テスト実行
			err := handlerWithMiddleware(c)

			// アサーション
			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, "success", rec.Body.String())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				// JSONレスポンスの検証
				var response map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name           string
		userID         interface{}
		expectedUserID uint
	}{
		{
			name:           "正常なユーザーID",
			userID:         uint(1),
			expectedUserID: 1,
		},
		{
			name:           "ユーザーIDなし",
			userID:         nil,
			expectedUserID: 0,
		},
		{
			name:           "無効な型",
			userID:         "invalid",
			expectedUserID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Echoの設定
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// コンテキストにユーザーIDを設定
			if tt.userID != nil {
				c.Set("user_id", tt.userID)
			}

			// テスト実行
			userID := GetUserIDFromContext(c)

			// アサーション
			assert.Equal(t, tt.expectedUserID, userID)
		})
	}
}

func TestJWTTokenValidation(t *testing.T) {
	// JWT_SECRETの設定
	os.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name           string
		userID         uint
		expired        bool
		expectedStatus int
	}{
		{
			name:           "有効なトークン",
			userID:         1,
			expired:        false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "期限切れトークン",
			userID:         1,
			expired:        true,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// トークンの生成
			exp := time.Now()
			if tt.expired {
				exp = exp.Add(-time.Hour) // 1時間前
			} else {
				exp = exp.Add(time.Hour) // 1時間後
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"user_id": tt.userID,
				"exp":     exp.Unix(),
				"iat":     time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("test-secret"))

			// Echoの設定
			e := echo.New()

			// テスト用のハンドラー
			handler := func(c echo.Context) error {
				userID := GetUserIDFromContext(c)
				assert.Equal(t, tt.userID, userID)
				return c.String(http.StatusOK, "success")
			}

			// ミドルウェアの適用
			middleware := AuthMiddleware()
			handlerWithMiddleware := middleware(handler)

			// リクエストの作成
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", "Bearer "+tokenString)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// テスト実行
			err := handlerWithMiddleware(c)

			// アサーション
			if tt.expectedStatus == http.StatusOK {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, "success", rec.Body.String())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}
