// package common は、ハンドラー間で共有される共通の型を提供します
package common

// RegisterUserRequest は、ユーザー登録APIのリクエストボディの構造を定義します
// バリデーションタグを使用して、各フィールドの制約を指定しています
type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required"`           // 名前（必須）
	Email    string `json:"email" validate:"required,email"`    // メールアドレス（必須、メール形式）
	Password string `json:"password" validate:"required,min=6"` // パスワード（必須、最小6文字）
}

// LoginRequest は、ログインAPIのリクエストボディの構造を定義します
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"` // メールアドレス（必須、メール形式）
	Password string `json:"password" validate:"required"`    // パスワード（必須）
}

// LoginResponse は、ログインAPIのレスポンスボディの構造を定義します
type LoginResponse struct {
	Token string `json:"token"`
}

// UpdateUserRequest は、ユーザー情報更新APIのリクエストボディの構造を定義します
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"required"`        // 名前（必須）
	Email string `json:"email" validate:"required,email"` // メールアドレス（必須、メール形式）
}

// PasswordResetRequest は、パスワードリセットリクエストAPIのリクエストボディの構造を定義します
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"` // メールアドレス（必須、メール形式）
}

// PasswordResetConfirmRequest は、パスワードリセット確認APIのリクエストボディの構造を定義します
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" validate:"required"`              // リセットトークン（必須）
	NewPassword string `json:"new_password" validate:"required,min=6"` // 新しいパスワード（必須、最小6文字）
}
