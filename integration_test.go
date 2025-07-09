package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"voice-link/domain/model"
	"voice-link/infrastructure/persistence"
	"voice-link/interface/handler"
	"voice-link/interface/router"
	"voice-link/usecase"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB は、テスト用のデータベースを設定します
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// マイグレーション
	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	return db
}

// setupTestApp は、テスト用のアプリケーションを設定します
func setupTestApp(t *testing.T) *echo.Echo {
	// JWT_SECRETの設定
	os.Setenv("JWT_SECRET", "test-secret")

	// テスト用データベースの設定
	db := setupTestDB(t)

	// 依存関係の注入
	userRepo := persistence.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	// Echoのインスタンスを作成
	e := echo.New()

	// ルーティングの設定
	r := router.NewRouter(e, userHandler)
	r.Setup()

	return e
}

func TestIntegration_UserRegistrationAndLogin(t *testing.T) {
	// テスト用アプリケーションの設定
	app := setupTestApp(t)

	// 1. ユーザー登録のテスト
	t.Run("ユーザー登録", func(t *testing.T) {
		registerData := map[string]interface{}{
			"name":     "テストユーザー",
			"email":    "test@example.com",
			"password": "password123",
		}

		jsonData, _ := json.Marshal(registerData)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "テストユーザー", response["name"])
		assert.Equal(t, "test@example.com", response["email"])
		assert.NotNil(t, response["id"])
	})

	// 2. ログインのテスト
	t.Run("ログイン", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "test@example.com",
			"password": "password123",
		}

		jsonData, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
	})

	// 3. 重複登録のテスト
	t.Run("重複登録エラー", func(t *testing.T) {
		registerData := map[string]interface{}{
			"name":     "重複ユーザー",
			"email":    "test@example.com", // 同じメールアドレス
			"password": "password123",
		}

		jsonData, _ := json.Marshal(registerData)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "email already exists", response["error"])
	})
}

func TestIntegration_ProtectedEndpoints(t *testing.T) {
	// テスト用アプリケーションの設定
	app := setupTestApp(t)

	// 1. ユーザー登録
	registerData := map[string]interface{}{
		"name":     "テストユーザー",
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(registerData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 2. ログインしてトークンを取得
	loginData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, _ = json.Marshal(loginData)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// 3. 保護されたエンドポイントのテスト
	t.Run("認証なしでアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "Authorization header is required", response["error"])
	})

	t.Run("有効なトークンでアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "テストユーザー", response["name"])
		assert.Equal(t, "test@example.com", response["email"])
	})

	t.Run("無効なトークンでアクセス", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "Invalid token", response["error"])
	})
}

func TestIntegration_UserCRUD(t *testing.T) {
	// テスト用アプリケーションの設定
	app := setupTestApp(t)

	// 1. ユーザー登録
	registerData := map[string]interface{}{
		"name":     "テストユーザー",
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(registerData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 2. ログインしてトークンを取得
	loginData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, _ = json.Marshal(loginData)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var loginResponse map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &loginResponse)
	token := loginResponse["token"].(string)

	// 3. ユーザー情報更新のテスト
	t.Run("ユーザー情報更新", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name":  "更新されたユーザー",
			"email": "updated@example.com",
		}

		jsonData, _ := json.Marshal(updateData)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "更新されたユーザー", response["name"])
		assert.Equal(t, "updated@example.com", response["email"])
	})

	// 4. 更新後のユーザー情報取得のテスト
	t.Run("更新後のユーザー情報取得", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rec := httptest.NewRecorder()

		app.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "更新されたユーザー", response["name"])
		assert.Equal(t, "updated@example.com", response["email"])
	})
}
