package user

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

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "正常なユーザー情報取得",
			userID: "1",
			mockSetup: func(mockUC *common.MockUserUseCase) {
				user := &model.User{
					ID:    1,
					Name:  "テストユーザー",
					Email: "test@example.com",
				}
				mockUC.On("GetByID", uint(1)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "ユーザーが見つからない",
			userID: "999",
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("GetByID", uint(999)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "User not found",
		},
		{
			name:   "無効なユーザーID",
			userID: "invalid",
			mockSetup: func(mockUC *common.MockUserUseCase) {
				// モックの設定は不要（パースエラーで早期リターン）
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(common.MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewUserHandler(mockUC)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			// ハンドラーの実行
			err := handler.GetUser(c)

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

func TestUserHandler_GetCurrentUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "正常な現在のユーザー情報取得",
			userID: 1,
			mockSetup: func(mockUC *common.MockUserUseCase) {
				user := &model.User{
					ID:    1,
					Name:  "テストユーザー",
					Email: "test@example.com",
				}
				mockUC.On("GetByID", uint(1)).Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "認証されていないユーザー",
			userID: 0,
			mockSetup: func(mockUC *common.MockUserUseCase) {
				// モックの設定は不要（認証エラーで早期リターン）
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name:   "ユーザーが見つからない",
			userID: 1,
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("GetByID", uint(1)).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(common.MockUserUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewUserHandler(mockUC)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)

			// ユーザーIDをコンテキストに設定
			if tt.userID != 0 {
				c.Set("user_id", tt.userID)
			}

			// ハンドラーの実行
			err := handler.GetCurrentUser(c)

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

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    common.UpdateUserRequest
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "正常なユーザー情報更新",
			userID: "1",
			requestBody: common.UpdateUserRequest{
				Name:  "更新されたユーザー",
				Email: "updated@example.com",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				user := &model.User{
					ID:    1,
					Name:  "更新されたユーザー",
					Email: "updated@example.com",
				}
				mockUC.On("UpdateUser", uint(1), "更新されたユーザー", "updated@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "無効なユーザーID",
			userID: "invalid",
			requestBody: common.UpdateUserRequest{
				Name:  "更新されたユーザー",
				Email: "updated@example.com",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				// モックの設定は不要（パースエラーで早期リターン）
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID",
		},
		{
			name:   "更新エラー",
			userID: "1",
			requestBody: common.UpdateUserRequest{
				Name:  "更新されたユーザー",
				Email: "updated@example.com",
			},
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("UpdateUser", uint(1), "更新されたユーザー", "updated@example.com").Return(nil, assert.AnError)
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
			handler := NewUserHandler(mockUC)

			// リクエストボディの準備
			reqBody, _ := json.Marshal(tt.requestBody)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+tt.userID, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			// ハンドラーの実行
			err := handler.UpdateUser(c)

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

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*common.MockUserUseCase)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "正常なユーザー削除",
			userID: "1",
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("DeleteUser", uint(1)).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "無効なユーザーID",
			userID: "invalid",
			mockSetup: func(mockUC *common.MockUserUseCase) {
				// モックの設定は不要（パースエラーで早期リターン）
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID",
		},
		{
			name:   "削除エラー",
			userID: "1",
			mockSetup: func(mockUC *common.MockUserUseCase) {
				mockUC.On("DeleteUser", uint(1)).Return(assert.AnError)
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
			handler := NewUserHandler(mockUC)

			// テスト用のリクエストとレスポンスを作成
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()

			// Echoコンテキストの作成
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			// ハンドラーの実行
			err := handler.DeleteUser(c)

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
