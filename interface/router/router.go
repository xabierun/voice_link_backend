package router

import (
	"voice-link/interface/handler/auth"
	"voice-link/interface/handler/user"
	authMiddleware "voice-link/interface/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	echo        *echo.Echo
	authHandler *auth.AuthHandler
	userHandler *user.UserHandler
}

func NewRouter(e *echo.Echo, authHandler *auth.AuthHandler, userHandler *user.UserHandler) *Router {
	return &Router{
		echo:        e,
		authHandler: authHandler,
		userHandler: userHandler,
	}
}

func (r *Router) Setup() {
	// ミドルウェアの設定
	r.echo.Use(echoMiddleware.Logger())
	r.echo.Use(echoMiddleware.Recover())
	r.echo.Use(echoMiddleware.CORS())

	// APIバージョン1のグループ
	v1 := r.echo.Group("/api/v1")

	// 認証不要なルーティング
	r.setupPublicRoutes(v1)

	// 認証が必要なルーティング
	r.setupProtectedRoutes(v1)
}

func (r *Router) setupPublicRoutes(api *echo.Group) {
	// 認証関連のルーティング
	auth := api.Group("/auth")
	{
		// ユーザー登録
		auth.POST("/register", r.authHandler.Register)
		// ログイン
		auth.POST("/login", r.authHandler.Login)
		// パスワードリセットリクエスト
		auth.POST("/password-reset", r.authHandler.RequestPasswordReset)
		// パスワードリセット確認
		auth.POST("/password-reset/confirm", r.authHandler.ResetPassword)
	}
}

func (r *Router) setupProtectedRoutes(api *echo.Group) {
	// 認証ミドルウェアを適用
	protected := api.Group("")
	protected.Use(authMiddleware.AuthMiddleware())

	// ユーザー関連のルーティング
	users := protected.Group("/users")
	{
		// 現在のユーザー情報の取得
		users.GET("/me", r.userHandler.GetCurrentUser)
		// 現在のユーザー情報の更新
		users.PUT("/me", r.userHandler.UpdateCurrentUser)
		// 現在のユーザーの削除
		users.DELETE("/me", r.userHandler.DeleteCurrentUser)

		// 管理者用のルーティング（特定のユーザーIDを指定）
		users.GET("/:id", r.userHandler.GetUser)
		users.PUT("/:id", r.userHandler.UpdateUser)
		users.DELETE("/:id", r.userHandler.DeleteUser)
	}
}
