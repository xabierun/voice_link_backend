package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"voice-link/domain/model"
	"voice-link/interface/handler/common"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    common.RegisterUserRequest
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なユーザー登録",
			requestBody: common.RegisterUserRequest{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				user := &model.User{
					ID:    1,
					Name:  "テストユーザー",
					Email: "test@example.com",
				}
				mockUC.On("Register", "テストユーザー", "test@example.com", "password123").Return(user, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "メールアドレス重複エラー",
			requestBody: common.RegisterUserRequest{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("Register", "テストユーザー", "test@example.com", "password123").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(common.MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewAuthHandler(mockUC)

			// リクエストボディの準備
			reqBody, _ := json.Marshal(tt.requestBody)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)

			// ハンドラーの実行
			err := handler.Register(c)

			// アサーション
			if tt.expectedError != "" {
				assert.Error(t, err)
				var response common.ErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			// モックの検証
			mockUC.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    common.LoginRequest
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なログイン",
			requestBody: common.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("Login", "test@example.com", "password123").Return("jwt-token", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "認証失敗",
			requestBody: common.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("Login", "test@example.com", "wrongpassword").Return("", assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(common.MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewAuthHandler(mockUC)

			// リクエストボディの準備
			reqBody, _ := json.Marshal(tt.requestBody)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)

			// ハンドラーの実行
			err := handler.Login(c)

			// アサーション
			if tt.expectedError != "" {
				assert.Error(t, err)
				var response common.ErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			// モックの検証
			mockUC.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    common.PasswordResetRequest
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なパスワードリセットリクエスト",
			requestBody: common.PasswordResetRequest{
				Email: "test@example.com",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("RequestPasswordReset", "test@example.com").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "パスワードリセットリクエストエラー",
			requestBody: common.PasswordResetRequest{
				Email: "test@example.com",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("RequestPasswordReset", "test@example.com").Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(common.MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewAuthHandler(mockUC)

			// リクエストボディの準備
			reqBody, _ := json.Marshal(tt.requestBody)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password-reset", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)

			// ハンドラーの実行
			err := handler.RequestPasswordReset(c)

			// アサーション
			if tt.expectedError != "" {
				assert.Error(t, err)
				var response common.ErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			// モックの検証
			mockUC.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    common.PasswordResetConfirmRequest
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なパスワードリセット",
			requestBody: common.PasswordResetConfirmRequest{
				Token:       "valid-token",
				NewPassword: "newpassword123",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("ResetPassword", "valid-token", "newpassword123").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "無効なトークン",
			requestBody: common.PasswordResetConfirmRequest{
				Token:       "invalid-token",
				NewPassword: "newpassword123",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("ResetPassword", "invalid-token", "newpassword123").Return(assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(common.MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewAuthHandler(mockUC)

			// リクエストボディの準備
			reqBody, _ := json.Marshal(tt.requestBody)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/password-reset/confirm", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)

			// ハンドラーの実行
			err := handler.ResetPassword(c)

			// アサーション
			if tt.expectedError != "" {
				assert.Error(t, err)
				var response common.ErrorResponse
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			// モックの検証
			mockUC.AssertExpectations(t)
		})
	}
}
