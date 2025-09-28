// package common は、ハンドラー間で共有される共通の型を提供します
package common

import (
	"voice-link/domain/model"

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

func (m *MockUserUseCase) RequestPasswordReset(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserUseCase) ResetPassword(token, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}
