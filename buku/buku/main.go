package main

import (
    "buku/config"
    "buku/routes"

    "buku/middleware"
    "github.com/labstack/echo/v4"
)

func main() {
    config.ConnectDB()
    e := echo.New()
    if config.DB == nil {
        panic("Database connection failed. Please check your configuration.")
    }
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    routes.AuthRoutes(e)
    routes.UserRoutes(e)
    routes.GenreRoutes(e)
    routes.PenulisRoutes(e)
    routes.BookRoutes(e)
    routes.PeminjamanRoutes(e)


    e.Logger.Fatal(e.Start(":8080"))
}
