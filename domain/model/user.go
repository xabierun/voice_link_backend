package model

import (
	"time"
)

type User struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	Name                 string     `json:"name" gorm:"not null"`
	Email                string     `json:"email" gorm:"unique;not null"`
	Password             string     `json:"-" gorm:"not null"`
	PasswordResetToken   *string    `json:"-" gorm:"unique"`
	PasswordResetExpires *time.Time `json:"-"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByPasswordResetToken(token string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}
