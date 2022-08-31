package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func getErrorHanlder(log *logrus.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {

		httpcode := http.StatusInternalServerError
		msg := "internal error"

		if he, ok := err.(*echo.HTTPError); ok {
			httpcode = he.Code
			msg = he.Message.(string)
		}

		if httpcode >= 500 {
			log.WithError(err).Error("request error")
		} else {
			log.WithError(err).Warnf("request warning")
		}

		_ = c.JSON(httpcode, ErrorMessage{
			Message: msg,
		})
	}
}

func setErrorHandler(e *echo.Echo, log *logrus.Logger) {
	e.HTTPErrorHandler = getErrorHanlder(log)
}
