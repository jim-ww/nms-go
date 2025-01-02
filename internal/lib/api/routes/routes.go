package routes

import (
	"github.com/jim-ww/nms-go/internal/features/auth/handler"
	"github.com/jim-ww/nms-go/internal/features/auth/middleware"
	"github.com/labstack/echo/v4"
)

func AddRoutes(e *echo.Echo, authHandler *handler.AuthHandler, middleware middleware.AuthMiddleware) {
	e.Static("/web", "./web")
	e.File("/favicon.ico", "web/favicon.ico")

	e.GET("/login", authHandler.LoginForm)
	e.GET("/register", authHandler.RegisterForm)

	api := e.Group("/api")
	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)
	api.POST("/logout", handler.Logout)

	user := e.Group("/", middleware.OnlyUser)
	admin := e.Group("/admin", middleware.OnlyAdmins)

	user.GET("", func(c echo.Context) error {
		return c.String(200, "Logged in as user!")
	})

	admin.GET("/admin", func(c echo.Context) error {
		return c.String(200, "Logged in as admin!")
	})
}
