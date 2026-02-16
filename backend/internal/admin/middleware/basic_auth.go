package middleware

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// BasicAuth - Basic認証ミドルウェア（管理画面用）
func BasicAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing authorization header")
			}

			if !strings.HasPrefix(auth, "Basic ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format")
			}

			payload, err := base64.StdEncoding.DecodeString(auth[6:])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid base64 encoding")
			}

			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != 2 {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials format")
			}

			username := pair[0]
			password := pair[1]

			expectedUser := os.Getenv("ADMIN_BASIC_USER")
			expectedPass := os.Getenv("ADMIN_BASIC_PASSWORD")

			if username != expectedUser || password != expectedPass {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
			}

			return next(c)
		}
	}
}
