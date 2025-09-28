package usecase

import (
	"errors"
	"os"
	"testing"
	"voice-link/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository は、UserRepositoryのモック実装です
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByPasswordResetToken(token string) (*model.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserUseCase_Register(t *testing.T) {
	// JWT_SECRETの設定
	os.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name          string
		nameInput     string
		emailInput    string
		passwordInput string
		mockSetup     func(*MockUserRepository)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:          "正常なユーザー登録",
			nameInput:     "テストユーザー",
			emailInput:    "test@example.com",
			passwordInput: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// FindByEmailでユーザーが見つからない場合
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, errors.New("user not found"))
				// Createでユーザー作成成功
				mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedUser: &model.User{
				Name:  "テストユーザー",
				Email: "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:          "メールアドレス重複",
			nameInput:     "テストユーザー",
			emailInput:    "existing@example.com",
			passwordInput: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				existingUser := &model.User{
					ID:    1,
					Name:  "既存ユーザー",
					Email: "existing@example.com",
				}
				mockRepo.On("FindByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedUser:  nil,
			expectedError: errors.New("email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// ユースケースの作成
			useCase := NewUserUseCase(mockRepo)

			// テスト実行
			user, err := useCase.Register(tt.nameInput, tt.emailInput, tt.passwordInput)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.NotEmpty(t, user.Password) // パスワードがハッシュ化されていることを確認
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_Login(t *testing.T) {
	// JWT_SECRETの設定
	os.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name          string
		emailInput    string
		passwordInput string
		mockSetup     func(*MockUserRepository)
		expectedToken string
		expectedError error
	}{
		{
			name:          "正常なログイン",
			emailInput:    "test@example.com",
			passwordInput: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				// ハッシュ化されたパスワードを作成
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					ID:       1,
					Name:     "テストユーザー",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectedToken: "", // 実際のトークンは動的に生成されるため空文字
			expectedError: nil,
		},
		{
			name:          "ユーザーが見つからない",
			emailInput:    "nonexistent@example.com",
			passwordInput: "password123",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))
			},
			expectedToken: "",
			expectedError: errors.New("invalid email or password"),
		},
		{
			name:          "パスワードが間違っている",
			emailInput:    "test@example.com",
			passwordInput: "wrongpassword",
			mockSetup: func(mockRepo *MockUserRepository) {
				// 正しいパスワードでハッシュ化
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					ID:       1,
					Name:     "テストユーザー",
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}
				mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
			},
			expectedToken: "",
			expectedError: errors.New("invalid email or password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// ユースケースの作成
			useCase := NewUserUseCase(mockRepo)

			// テスト実行
			token, err := useCase.Login(tt.emailInput, tt.passwordInput)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				// JWTトークンの形式を簡単にチェック（.で区切られている）
				assert.Contains(t, token, ".")
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		idInput       uint
		mockSetup     func(*MockUserRepository)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:    "正常なユーザー取得",
			idInput: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &model.User{
					ID:    1,
					Name:  "テストユーザー",
					Email: "test@example.com",
				}
				mockRepo.On("FindByID", uint(1)).Return(user, nil)
			},
			expectedUser: &model.User{
				ID:    1,
				Name:  "テストユーザー",
				Email: "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:    "ユーザーが見つからない",
			idInput: 999,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// ユースケースの作成
			useCase := NewUserUseCase(mockRepo)

			// テスト実行
			user, err := useCase.GetByID(tt.idInput)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		idInput       uint
		nameInput     string
		emailInput    string
		mockSetup     func(*MockUserRepository)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:       "正常なユーザー更新",
			idInput:    1,
			nameInput:  "更新されたユーザー",
			emailInput: "updated@example.com",
			mockSetup: func(mockRepo *MockUserRepository) {
				user := &model.User{
					ID:    1,
					Name:  "元のユーザー",
					Email: "original@example.com",
				}
				mockRepo.On("FindByID", uint(1)).Return(user, nil)
				mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil)
			},
			expectedUser: &model.User{
				ID:    1,
				Name:  "更新されたユーザー",
				Email: "updated@example.com",
			},
			expectedError: nil,
		},
		{
			name:       "ユーザーが見つからない",
			idInput:    999,
			nameInput:  "更新されたユーザー",
			emailInput: "updated@example.com",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// ユースケースの作成
			useCase := NewUserUseCase(mockRepo)

			// テスト実行
			user, err := useCase.UpdateUser(tt.idInput, tt.nameInput, tt.emailInput)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		idInput       uint
		mockSetup     func(*MockUserRepository)
		expectedError error
	}{
		{
			name:    "正常なユーザー削除",
			idInput: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Delete", uint(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:    "ユーザーが見つからない",
			idInput: 999,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Delete", uint(999)).Return(errors.New("user not found"))
			},
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			// ユースケースの作成
			useCase := NewUserUseCase(mockRepo)

			// テスト実行
			err := useCase.DeleteUser(tt.idInput)

			// アサーション
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
