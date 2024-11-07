package middleware

import (
    "net/http"
    "github.com/labstack/echo/v4"
	"fmt"
)

// Logger middleware
func Logger() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Log request
            fmt.Printf("Request Method: %s, Request URI: %s\n", c.Request().Method, c.Request().RequestURI)
            return next(c)
        }
    }
}

// Recover middleware
func Recover() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            defer func() {
                if err := recover(); err != nil {
                    c.Logger().Error(err)
                    c.String(http.StatusInternalServerError, "Internal Server Error")
                }
            }()
            return next(c)
        }
    }
}
