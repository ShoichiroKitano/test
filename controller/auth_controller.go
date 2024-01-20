package controller

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthController struct{}

func (controller *AuthController) Login(c echo.Context) error {
	// TODO: idとpassowrd 認証する。パスワードはhash化して比較する。
	sess, err := session.Get("session", c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["user_id"] = uint64(1)
	sess.Save(c.Request(), c.Response())
	return c.NoContent(http.StatusOK)
}

func UserID(c echo.Context) (uint64, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return 0, err
	}
	return sess.Values["user_id"].(uint64), nil
}
