// package auth は、認証関連のHTTPリクエストを処理するハンドラーを提供します
package auth

import (
	"net/http"
	"voice-link/interface/handler/common"
	"voice-link/usecase"

	"github.com/labstack/echo/v4"
)

// AuthHandler は、認証関連のHTTPリクエストを処理するハンドラー構造体です
type AuthHandler struct {
	userUseCase usecase.UserUseCase
}

// NewAuthHandler は、AuthHandlerの新しいインスタンスを作成するファクトリ関数です
func NewAuthHandler(userUseCase usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{userUseCase}
}

// Register は、新規ユーザーを登録するためのハンドラー関数です
// POSTリクエストを受け取り、ユーザー情報をデータベースに保存します
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(common.RegisterUserRequest)
	if err := c.Bind(req); err != nil { // リクエストボディをバインド
		return common.SendBadRequestError(c, "Invalid request body")
	}

	// ユースケースレイヤーを呼び出してユーザー登録を実行
	user, err := h.userUseCase.Register(req.Name, req.Email, req.Password)

	// エラーが発生した場合
	if err != nil {
		return common.SendInternalServerError(c, err.Error())
	}

	return c.JSON(http.StatusCreated, user) // 201 Createdとユーザー情報を返却
}

// Login は、ユーザーログインを処理するハンドラー関数です
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(common.LoginRequest)
	if err := c.Bind(req); err != nil {
		return common.SendBadRequestError(c, "Invalid request body")
	}

	// ユースケースレイヤーを呼び出してログインを実行
	token, err := h.userUseCase.Login(req.Email, req.Password)
	if err != nil {
		return common.SendUnauthorizedError(c, err.Error())
	}

	return c.JSON(http.StatusOK, common.LoginResponse{Token: token})
}

// RequestPasswordReset は、パスワードリセットのリクエストを処理するハンドラー関数です
func (h *AuthHandler) RequestPasswordReset(c echo.Context) error {
	req := new(common.PasswordResetRequest)
	if err := c.Bind(req); err != nil {
		return common.SendBadRequestError(c, "Invalid request body")
	}

	// ユースケースレイヤーを呼び出してパスワードリセットリクエストを実行
	if err := h.userUseCase.RequestPasswordReset(req.Email); err != nil {
		return common.SendInternalServerError(c, err.Error())
	}

	// セキュリティ上の理由で、常に成功レスポンスを返す
	return common.SendMessageResponse(c, http.StatusOK, "If the email exists, a password reset link has been sent")
}

// ResetPassword は、パスワードリセットトークンを使用してパスワードをリセットするハンドラー関数です
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	req := new(common.PasswordResetConfirmRequest)
	if err := c.Bind(req); err != nil {
		return common.SendBadRequestError(c, "Invalid request body")
	}

	// ユースケースレイヤーを呼び出してパスワードリセットを実行
	if err := h.userUseCase.ResetPassword(req.Token, req.NewPassword); err != nil {
		return common.SendBadRequestError(c, err.Error())
	}

	return common.SendMessageResponse(c, http.StatusOK, "Password has been reset successfully")
}
