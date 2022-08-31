package server

import (
	"varmijo/mortage-calculator/pkg/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

//Register routes and middlewers
func Register(e *echo.Echo, controller usecase.Handler, log *logrus.Logger) error {
	//Set error handler middlewar
	setErrorHandler(e, log)

	//Enable CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	//Creates api routes set
	apiv1 := e.Group("api/v1")

	//Register route
	apiv1.GET("/mortage", getMortage(controller))

	return nil
}
