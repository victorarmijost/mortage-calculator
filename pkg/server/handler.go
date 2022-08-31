package server

import (
	"fmt"
	"net/http"
	"strings"
	"varmijo/mortage-calculator/pkg/entities"
	"varmijo/mortage-calculator/pkg/usecase"

	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
)

//Defines the allowed inputs for the payment schedule values
const (
	acceleratedBiWeekly = "accelerated-bi-weekly"
	biWeekly            = "bi-weekly"
	montly              = "monthly"
)

//Get mortage http handler
func getMortage(controller usecase.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		//Parse query parameters
		query := &GetMortageData{}
		err := schema.NewDecoder().Decode(query, c.QueryParams())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		//Get entity from input values
		data, err := getEntity(query)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		//Execute bussines logic
		res, err := controller.GetMortage(ctx, data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, ResponseMortageData{
			Payment:   res.Payment,
			Insurance: res.Insurance,
		})
	}
}

//Try to get an entity object with the user data
func getEntity(data *GetMortageData) (*entities.MortageData, error) {
	//Check inputs
	if err := validateData(data); err != nil {
		return nil, err
	}

	//Normalize to lower case, to allow the user to send it with upper case mixes
	norm_schedule := strings.ToLower(data.Schedule)

	//Maps the payment schedule sent by the user to its equivalent entity value
	var entity_schedule entities.PaymentSchedule
	switch norm_schedule {
	case acceleratedBiWeekly:
		entity_schedule = entities.AcceleratedBiWeekly
	case biWeekly:
		entity_schedule = entities.BiWeekly
	case montly:
		entity_schedule = entities.Montly
	default:
		return nil, fmt.Errorf("wrong payment schedule value")
	}

	return &entities.MortageData{
		Price:           data.Price,
		DownPayment:     data.DownPayment,
		InterestRate:    data.InteresRate,
		Period:          data.Period,
		SchedulePeriods: entity_schedule,
	}, nil
}

//Checks that the input values are correct
func validateData(data *GetMortageData) error {
	if data.Price <= 0 {
		return fmt.Errorf("wrong price value")
	}

	if data.DownPayment <= 0 || data.DownPayment > data.Price {
		return fmt.Errorf("wrong down payment value")
	}

	if data.InteresRate <= 0 || data.InteresRate > 100 {
		return fmt.Errorf("wrong interest rate value")
	}

	if data.Period == 0 || data.Period%5 != 0 || data.Period > 30 {
		return fmt.Errorf("wrong period value")
	}

	return nil
}
