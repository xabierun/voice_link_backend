// package common は、ハンドラー間で共有される共通の型を提供します
package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponse は、エラーレスポンスの構造を定義します
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse は、メッセージレスポンスの構造を定義します
type MessageResponse struct {
	Message string `json:"message"`
}

// SendErrorResponse は、エラーレスポンスを送信するヘルパー関数です
func SendErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, ErrorResponse{Error: message})
}

// SendMessageResponse は、メッセージレスポンスを送信するヘルパー関数です
func SendMessageResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, MessageResponse{Message: message})
}

// SendBadRequestError は、400 Bad Requestエラーを送信します
func SendBadRequestError(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusBadRequest, message)
}

// SendInternalServerError は、500 Internal Server Errorを送信します
func SendInternalServerError(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusInternalServerError, message)
}

// SendUnauthorizedError は、401 Unauthorizedエラーを送信します
func SendUnauthorizedError(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusUnauthorized, message)
}

// SendNotFoundError は、404 Not Foundエラーを送信します
func SendNotFoundError(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusNotFound, message)
}
