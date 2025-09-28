// package user は、ユーザー情報管理のHTTPリクエストを処理するハンドラーを提供します
package user

import (
	"net/http"
	"strconv"
	"voice-link/interface/handler/common"
	"voice-link/interface/middleware"
	"voice-link/usecase"

	"github.com/labstack/echo/v4"
)

// UserHandler は、ユーザー情報管理のHTTPリクエストを処理するハンドラー構造体です
// userUseCaseフィールドには、ビジネスロジックを実行するためのインターフェースが格納されます
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler は、UserHandlerの新しいインスタンスを作成するファクトリ関数です
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

// GetUser は、指定されたIDのユーザー情報を取得するハンドラー関数です
// URLパラメータからユーザーIDを取得し、該当するユーザー情報を返します
func (h *UserHandler) GetUser(c echo.Context) error {
	// URLパラメータからIDを取得し、uint型に変換
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return common.SendBadRequestError(c, "Invalid user ID")
	}

	// ユースケースレイヤーを呼び出してユーザー情報を取得
	user, err := h.userUseCase.GetByID(uint(id))
	if err != nil {
		return common.SendNotFoundError(c, "User not found")
	}

	return c.JSON(http.StatusOK, user) // 200 OKとユーザー情報を返却
}

// GetCurrentUser は、現在ログインしているユーザーの情報を取得するハンドラー関数です
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return common.SendUnauthorizedError(c, "User not authenticated")
	}

	user, err := h.userUseCase.GetByID(userID)
	if err != nil {
		return common.SendNotFoundError(c, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser は、指定されたIDのユーザー情報を更新するハンドラー関数です
// URLパラメータからユーザーIDを取得し、リクエストボディの内容でユーザー情報を更新します
func (h *UserHandler) UpdateUser(c echo.Context) error {
	// URLパラメータからIDを取得し、uint型に変換
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return common.SendBadRequestError(c, "Invalid user ID")
	}

	req := new(common.UpdateUserRequest)
	if err := c.Bind(req); err != nil { // リクエストボディをバインド
		return common.SendBadRequestError(c, "Invalid request body")
	}

	// ユースケースレイヤーを呼び出してユーザー情報を更新
	user, err := h.userUseCase.UpdateUser(uint(id), req.Name, req.Email)
	if err != nil {
		return common.SendInternalServerError(c, err.Error())
	}

	return c.JSON(http.StatusOK, user) // 200 OKと更新後のユーザー情報を返却
}

// UpdateCurrentUser は、現在ログインしているユーザーの情報を更新するハンドラー関数です
func (h *UserHandler) UpdateCurrentUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return common.SendUnauthorizedError(c, "User not authenticated")
	}

	req := new(common.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return common.SendBadRequestError(c, "Invalid request body")
	}

	user, err := h.userUseCase.UpdateUser(userID, req.Name, req.Email)
	if err != nil {
		return common.SendInternalServerError(c, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser は、指定されたIDのユーザーを削除するハンドラー関数です
// URLパラメータからユーザーIDを取得し、該当するユーザーを削除します
func (h *UserHandler) DeleteUser(c echo.Context) error {
	// URLパラメータからIDを取得し、uint型に変換
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return common.SendBadRequestError(c, "Invalid user ID")
	}

	// ユースケースレイヤーを呼び出してユーザーを削除
	if err := h.userUseCase.DeleteUser(uint(id)); err != nil {
		return common.SendInternalServerError(c, err.Error())
	}

	return c.NoContent(http.StatusNoContent) // 204 No Contentを返却
}

// DeleteCurrentUser は、現在ログインしているユーザーを削除するハンドラー関数です
func (h *UserHandler) DeleteCurrentUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == 0 {
		return common.SendUnauthorizedError(c, "User not authenticated")
	}

	if err := h.userUseCase.DeleteUser(userID); err != nil {
		return common.SendInternalServerError(c, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
