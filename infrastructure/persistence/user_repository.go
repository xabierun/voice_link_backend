// package persistence は、データベースとの永続化層を提供します
package persistence

import (
	"voice-link/domain/model"

	"gorm.io/gorm"
)

// userRepository は、ユーザー情報のデータベース操作を担当する構造体です
type userRepository struct {
	db *gorm.DB // データベースコネクション
}

// NewUserRepository は、UserRepositoryインターフェースの新しいインスタンスを作成します
func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{db}
}

// Create は、新しいユーザーをデータベースに作成します
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID は、指定されたIDのユーザーをデータベースから検索します
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail は、指定されたメールアドレスのユーザーをデータベースから検索します
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update は、既存のユーザー情報をデータベースで更新します
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete は、指定されたIDのユーザーをデータベースから削除します
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
