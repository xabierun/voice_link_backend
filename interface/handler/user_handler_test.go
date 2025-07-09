package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"voice-link/domain/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUseCase は、UserUseCaseのモック実装です
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(name, email, password string) (*model.User, error) {
	args := m.Called(name, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUseCase) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserUseCase) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(id uint, name, email string) (*model.User, error) {
	args := m.Called(id, name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserUseCase) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    RegisterUserRequest
		mockSetup      func(*MockUserUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常なユーザー登録",
			requestBody: RegisterUserRequest{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUC *MockUserUseCase) {
				expectedUser := &model.User{
					ID:    1,
					Name:  "テストユーザー",
					Email: "test@example.com",
				}
				mockUC.On("Register", "テストユーザー", "test@example.com", "password123").Return(expectedUser, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":    float64(1),
				"name":  "テストユーザー",
				"email": "test@example.com",
			},
		},
		{
			name: "無効なリクエストボディ",
			requestBody: RegisterUserRequest{
				Name:     "",
				Email:    "invalid-email",
				Password: "123",
			},
			mockSetup:      func(mockUC *MockUserUseCase) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewUserHandler(mockUC)

			// Echoの設定
			e := echo.New()
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// テスト実行
			err := handler.Register(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockUC.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    LoginRequest
		mockSetup      func(*MockUserUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常なログイン",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(mockUC *MockUserUseCase) {
				mockUC.On("Login", "test@example.com", "password123").Return("jwt-token-here", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"token": "jwt-token-here",
			},
		},
		{
			name: "認証失敗",
			requestBody: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(mockUC *MockUserUseCase) {
				mockUC.On("Login", "test@example.com", "wrongpassword").Return("", assert.AnError)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewUserHandler(mockUC)

			// Echoの設定
			e := echo.New()
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// テスト実行
			err := handler.Login(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockUC.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		mockSetup      func(*MockUserUseCase)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "正常なユーザー情報取得",
			userID: 1,
			mockSetup: func(mockUC *MockUserUseCase) {
				expectedUser := &model.User{
					ID:    1,
					Name:  "テストユーザー",
					Email: "test@example.com",
				}
				mockUC.On("GetByID", uint(1)).Return(expectedUser, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":    float64(1),
				"name":  "テストユーザー",
				"email": "test@example.com",
			},
		},
		{
			name:   "ユーザーが見つからない",
			userID: 999,
			mockSetup: func(mockUC *MockUserUseCase) {
				mockUC.On("GetByID", uint(999)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "User not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewUserHandler(mockUC)

			// Echoの設定
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", tt.userID)

			// テスト実行
			err := handler.GetCurrentUser(c)

			// アサーション
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)

			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockUC.AssertExpectations(t)
		})
	}
}
