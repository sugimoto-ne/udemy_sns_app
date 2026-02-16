package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/admin/handlers"
	adminMiddleware "github.com/yourusername/sns-backend/internal/admin/middleware"
	"github.com/yourusername/sns-backend/internal/admin/renderer"
)

// SetupAdminRoutes - 管理画面ルート設定
func SetupAdminRoutes(e *echo.Echo) {
	// テンプレートレンダラー設定
	templateRenderer, err := renderer.NewTemplateRenderer("internal/admin/templates")
	if err != nil {
		panic("Failed to initialize template renderer: " + err.Error())
	}
	e.Renderer = templateRenderer

	// 静的ファイル配信
	e.Static("/static", "static")

	// Basic認証を全体に適用（一時的に無効化）
	admin := e.Group("/admin") // adminMiddleware.BasicAuth()

	// ハンドラー初期化
	authHandler := handlers.NewAuthHandler()
	dashboardHandler := handlers.NewDashboardHandler()
	userHandler := handlers.NewUserHandler()
	passwordResetHandler := handlers.NewPasswordResetAdminHandler()
	logHandler := handlers.NewLogHandler()

	// 認証不要のルート
	admin.GET("/login", authHandler.ShowLoginPage)
	admin.POST("/login", authHandler.Login)

	// 管理者JWT認証必須のルート（admin_tokenを使用）
	adminAuth := admin.Group("", adminMiddleware.AdminJWTAuth())

	// 認証後のルート
	adminAuth.POST("/logout", authHandler.Logout)

	// 画面ルート
	adminAuth.GET("/dashboard", dashboardHandler.ShowDashboard)
	adminAuth.GET("/users", userHandler.ShowUserList)
	adminAuth.GET("/users/:id", userHandler.ShowUserDetail)
	adminAuth.GET("/password-resets", passwordResetHandler.ShowPasswordResetList)
	adminAuth.GET("/logs", logHandler.ShowLogList)

	// APIルート
	api := adminAuth.Group("/api")

	// ダッシュボードAPI
	api.GET("/dashboard/stats", dashboardHandler.GetDashboardStats)
	api.GET("/dashboard/charts/posts", dashboardHandler.GetPostsChartData)
	api.GET("/dashboard/charts/users", dashboardHandler.GetUsersChartData)

	// ユーザー管理API
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/:id", userHandler.GetUser)
	api.PATCH("/users/:id/status", userHandler.UpdateUserStatus)
	api.POST("/users/batch-update-status", userHandler.BatchUpdateUserStatus)

	// パスワードリセットAPI
	api.GET("/password-resets", passwordResetHandler.GetPasswordResets)
	api.POST("/password-resets/:id/approve", passwordResetHandler.ApproveResetRequest)

	// 操作ログAPI
	api.GET("/logs", logHandler.GetLogs)
}
