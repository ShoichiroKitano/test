package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"super-shiharai-kun/controller"
	"super-shiharai-kun/driver"
	"super-shiharai-kun/repository"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret")))) // TODO: 仮実装

	rdbDriver, err := driver.NewRDBDriver()
	if err != nil {
		e.Logger.Fatal(err)
	}
	repo := &repository.RDBRepository{Driver: rdbDriver}

	requireLogin := e.Group("/api")
	auth := &controller.AuthController{}
	requireLogin.POST("/login", auth.Login)

	invoice := &controller.InvoiceController{RDB: repo}
	requireLogin.POST("/invoices", invoice.Create)

	e.Logger.Fatal(e.Start("localhost:3000"))
}
