package routes

import (
	"github.com/jim-ww/nms-go/internal/features/auth/handler"
	"github.com/jim-ww/nms-go/internal/features/auth/middleware"
	"github.com/jim-ww/nms-go/internal/lib/api/handlers/health"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, authHandler *handler.AuthHandler, middleware middleware.AuthMiddleware) {
	e.Static("/web", "./web")
	e.File("/favicon.ico", "web/favicon.ico")

	api := e.Group("/api")
	api.GET("/health", health.Health)
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)

	unauthorizedOnly := e.Group("", middleware.OnlyUnauthorized)
	unauthorizedOnly.GET("/login", authHandler.LoginForm)
	unauthorizedOnly.GET("/register", authHandler.RegisterForm)

	user := e.Group("", middleware.OnlyUser)
	user.GET("/logout", handler.Logout)
	user.GET("", func(c echo.Context) error {
		return c.String(200, "Logged in as user")
	})

	admin := e.Group("/admin", middleware.OnlyAdmins)
	admin.GET("/admin", func(c echo.Context) error {
		return c.String(200, "Logged in as admin")
	})
}
