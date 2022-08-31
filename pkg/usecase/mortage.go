package usecase

import (
	"context"
	"fmt"
	"math"
	"varmijo/mortage-calculator/pkg/entities"
)

//Calculate the mortage value given the input data
func (u *interactor) GetMortage(ctx context.Context, data *entities.MortageData) (*entities.ResultData, error) {
	if data.DownPayment > data.Price {
		return nil, fmt.Errorf("unexpected down payment value")
	}

	if data.Price > entities.MinDownPaymentPrice {
		limit := entities.MinDownPaymentPercent * entities.MinDownPaymentPrice
		limit = limit + (data.Price-entities.MinDownPaymentPrice)*entities.ComplementDownPaymentPercent

		if data.DownPayment < limit {
			return nil, fmt.Errorf("for the given price the down payment is too low")
		}
	}

	if data.DownPayment < entities.MinDownPaymentPercent*data.Price {
		return nil, fmt.Errorf("down payment is too low")
	}

	if data.Price > entities.MaxCMHCPrice && data.DownPayment < entities.MaxDownPaymentPercent*data.Price {
		return nil, fmt.Errorf("for the given price, the down payment is too low")
	}

	if data.Period > entities.MaxCMHCPeriod && data.DownPayment < entities.MaxDownPaymentPercent*data.Price {
		return nil, fmt.Errorf("for the given period, the down payment is too low")
	}

	return u.getBaseMortage(data), nil
}

//Calculate mortage value
func (u *interactor) getBaseMortage(data *entities.MortageData) *entities.ResultData {
	u.logger.Debugf("received data: %+v", data)

	principal := data.Price - data.DownPayment

	CMHCRate := u.getCMHCRate(data)

	u.logger.Debugf("CHMRate to be used: %f", CMHCRate)

	CMHCValue := CMHCRate * principal

	principal = principal + CMHCValue

	number_payments := float64(data.Period * uint(data.SchedulePeriods))

	interes_rate := data.InterestRate / float64(data.SchedulePeriods) / 100

	power := math.Pow(1+interes_rate, number_payments)

	payment := principal * interes_rate * power / (power - 1)

	u.logger.Debugf("calculated mortage value: %f", payment)

	return &entities.ResultData{
		Payment:   payment,
		Insurance: CMHCValue,
	}
}

//Calculate CMHC
func (u *interactor) getCMHCRate(data *entities.MortageData) float64 {
	downPaymentRate := data.DownPayment / data.Price * 100

	u.logger.Debugf("down payment rate: %f", downPaymentRate)

	if downPaymentRate < 5 {
		return 0
	}

	if downPaymentRate >= 5 && downPaymentRate < 10 {
		return 0.04
	}

	if downPaymentRate >= 10 && downPaymentRate < 15 {
		return 0.031
	}

	if downPaymentRate >= 15 && downPaymentRate < 20 {
		return 0.028
	}

	return 0
}
