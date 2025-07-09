package usecase

import (
	"errors"
	"os"
	"time"
	"voice-link/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Register(name, email, password string) (*model.User, error)
	Login(email, password string) (string, error)
	GetByID(id uint) (*model.User, error)
	UpdateUser(id uint, name, email string) (*model.User, error)
	DeleteUser(id uint) error
}

type userUseCase struct {
	userRepo model.UserRepository
}

func NewUserUseCase(userRepo model.UserRepository) UserUseCase {
	return &userUseCase{userRepo}
}

func (u *userUseCase) Register(name, email, password string) (*model.User, error) {
	// メールアドレスの重複チェック
	existingUser, err := u.userRepo.FindByEmail(email)

	// エラーがなく、既存のユーザーが存在する場合
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// ハッシュ化に失敗した場合
	if err != nil {
		return nil, err
	}

	// 新しいユーザーを作成
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	// ユーザーをデータベースに作成
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 作成したユーザーを返す
	return user, nil
}

func (u *userUseCase) Login(email, password string) (string, error) {
	// メールアドレスでユーザーを検索
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24時間有効
		"iat":     time.Now().Unix(),
	})

	// トークンの署名
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *userUseCase) GetByID(id uint) (*model.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUseCase) UpdateUser(id uint, name, email string) (*model.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Email = email

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) DeleteUser(id uint) error {
	return u.userRepo.Delete(id)
}
