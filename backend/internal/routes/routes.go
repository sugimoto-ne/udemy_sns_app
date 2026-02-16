package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/sns-backend/internal/handlers"
	"github.com/yourusername/sns-backend/internal/middleware"
)

// SetupRoutes - ルート設定
func SetupRoutes(e *echo.Echo) {
	// APIグループ
	api := e.Group("/api/v1")

	// 認証ルート
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", handlers.GetMe, middleware.JWTAuth())
		auth.POST("/refresh", handlers.RefreshToken)       // リフレッシュトークンでアクセストークンを再発行
		auth.POST("/logout", handlers.Logout)              // ログアウト（認証不要）
		auth.POST("/revoke-all", handlers.RevokeAllTokens, middleware.JWTAuth()) // 全デバイスログアウト
	}

	// ユーザールート
	users := api.Group("/users")
	{
		users.GET("/:username", handlers.GetUserByUsername, middleware.OptionalJWTAuth())
		users.PUT("/me", handlers.UpdateProfile, middleware.JWTAuth())
		users.GET("/:username/posts", handlers.GetUserPosts)
		users.GET("/:username/followers", handlers.GetFollowers)
		users.GET("/:username/following", handlers.GetFollowing)
		users.POST("/:username/follow", handlers.FollowUser, middleware.JWTAuth())
		users.DELETE("/:username/follow", handlers.UnfollowUser, middleware.JWTAuth())
	}

	// 投稿ルート
	posts := api.Group("/posts")
	{
		posts.GET("", handlers.GetTimeline, middleware.OptionalJWTAuth())
		posts.GET("/:id", handlers.GetPostByID, middleware.OptionalJWTAuth())
		posts.POST("", handlers.CreatePost, middleware.JWTAuth())
		posts.PUT("/:id", handlers.UpdatePost, middleware.JWTAuth())
		posts.DELETE("/:id", handlers.DeletePost, middleware.JWTAuth())

		// コメントルート
		posts.GET("/:id/comments", handlers.GetComments)
		posts.POST("/:id/comments", handlers.CreateComment, middleware.JWTAuth())

		// いいねルート
		posts.POST("/:id/like", handlers.LikePost, middleware.JWTAuth())
		posts.DELETE("/:id/like", handlers.UnlikePost, middleware.JWTAuth())
		posts.GET("/:id/likes", handlers.GetLikes)
	}

	// コメント削除ルート
	api.DELETE("/comments/:id", handlers.DeleteComment, middleware.JWTAuth())
}
