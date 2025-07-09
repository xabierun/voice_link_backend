// package handler は、HTTPリクエストを処理するハンドラーを提供します
package handler

import (
	"net/http"
	"strconv"
	"voice-link/interface/middleware"
	"voice-link/usecase" // ビジネスロジックを含むusecaseパッケージをインポート

	"github.com/labstack/echo/v4" // Webフレームワーク Echo を使用
)

// UserHandler は、ユーザー関連のHTTPリクエストを処理するハンドラー構造体です
// userUseCaseフィールドには、ビジネスロジックを実行するためのインターフェースが格納されます
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler は、UserHandlerの新しいインスタンスを作成するファクトリ関数です
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

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

// Register は、新規ユーザーを登録するためのハンドラー関数です
// POSTリクエストを受け取り、ユーザー情報をデータベースに保存します
func (h *UserHandler) Register(c echo.Context) error {
	req := new(RegisterUserRequest)
	if err := c.Bind(req); err != nil { // リクエストボディをバインド
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// ユースケースレイヤーを呼び出してユーザー登録を実行
	user, err := h.userUseCase.Register(req.Name, req.Email, req.Password)

	// エラーが発生した場合
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, user) // 201 Createdとユーザー情報を返却
}

// Login は、ユーザーログインを処理するハンドラー関数です
func (h *UserHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// ユースケースレイヤーを呼び出してログインを実行
	token, err := h.userUseCase.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// GetUser は、指定されたIDのユーザー情報を取得するハンドラー関数です
// URLパラメータからユーザーIDを取得し、該当するユーザー情報を返します
func (h *UserHandler) GetUser(c echo.Context) error {
	// URLパラメータからIDを取得し、uint型に変換
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// ユースケースレイヤーを呼び出してユーザー情報を取得
	user, err := h.userUseCase.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user) // 200 OKとユーザー情報を返却
}

// GetCurrentUser は、現在ログインしているユーザーの情報を取得するハンドラー関数です
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	user, err := h.userUseCase.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser は、指定されたIDのユーザー情報を更新するハンドラー関数です
// URLパラメータからユーザーIDを取得し、リクエストボディの内容でユーザー情報を更新します
func (h *UserHandler) UpdateUser(c echo.Context) error {
	// URLパラメータからIDを取得し、uint型に変換
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	req := new(UpdateUserRequest)
	if err := c.Bind(req); err != nil { // リクエストボディをバインド
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// ユースケースレイヤーを呼び出してユーザー情報を更新
	user, err := h.userUseCase.UpdateUser(uint(id), req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user) // 200 OKと更新後のユーザー情報を返却
}

// UpdateCurrentUser は、現在ログインしているユーザーの情報を更新するハンドラー関数です
func (h *UserHandler) UpdateCurrentUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	req := new(UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	user, err := h.userUseCase.UpdateUser(userID, req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser は、指定されたIDのユーザーを削除するハンドラー関数です
// URLパラメータからユーザーIDを取得し、該当するユーザーを削除します
func (h *UserHandler) DeleteUser(c echo.Context) error {
	// URLパラメータからIDを取得し、uint型に変換
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// ユースケースレイヤーを呼び出してユーザーを削除
	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent) // 204 No Contentを返却
}

// DeleteCurrentUser は、現在ログインしているユーザーを削除するハンドラー関数です
func (h *UserHandler) DeleteCurrentUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User not authenticated",
		})
	}

	if err := h.userUseCase.DeleteUser(userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
