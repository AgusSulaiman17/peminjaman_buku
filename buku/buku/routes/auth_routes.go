package routes

import (
    "github.com/labstack/echo/v4"
    "buku/controllers"
)

func AuthRoutes(e *echo.Echo) {
    e.POST("/register", controllers.Register)
    e.POST("/login", controllers.Login)
}
