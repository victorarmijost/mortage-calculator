package usecase

import (
	"context"
	"varmijo/mortage-calculator/pkg/entities"

	"github.com/sirupsen/logrus"
)

//Defines the methods for our usecases
type Handler interface {
	GetMortage(ctx context.Context, data *entities.MortageData) (*entities.ResultData, error)
}

//Creates a new use case handler
func New(logger *logrus.Logger) Handler {
	return &interactor{
		logger: logger,
	}
}

//Implements the Handler interface
type interactor struct {
	logger *logrus.Logger
}
