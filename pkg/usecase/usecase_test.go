package usecase

import (
	"context"
	"testing"
	"varmijo/mortage-calculator/pkg/entities"

	"github.com/sirupsen/logrus"
)

type testCase struct {
	name    string
	Data    *entities.MortageData
	isError bool
	mortage float64
}

func TestMortageUsecase(t *testing.T) {
	result_prefix := "error running get mortage usecase test, running case:"

	logger := logrus.New()
	ctx := context.Background()

	controller := New(logger)

	tt := []testCase{
		{
			name: "happy path",
			Data: &entities.MortageData{
				Price:           1000,
				DownPayment:     500,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			mortage: 4.894728537748993,
		},
		{
			name: "down payment too big",
			Data: &entities.MortageData{
				Price:           500,
				DownPayment:     1000,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			isError: true,
		},
		{
			name: "down payment less than minimum",
			Data: &entities.MortageData{
				Price:           1000,
				DownPayment:     40,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			isError: true,
		},
		{
			name: "down payment less than minimum for big prices",
			Data: &entities.MortageData{
				Price:           1000,
				DownPayment:     40,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			isError: true,
		},
		{
			name: "price over max price value",
			Data: &entities.MortageData{
				Price:           entities.MaxCMHCPrice + 1,
				DownPayment:     entities.MaxCMHCPrice * 0.1,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			isError: true,
		},
		{
			name: "period to big for the given down payment",
			Data: &entities.MortageData{
				Price:           10000,
				DownPayment:     1000,
				InterestRate:    10,
				Period:          30,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			isError: true,
		},
		{
			name: "price over min down payment price",
			Data: &entities.MortageData{
				Price:           entities.MinDownPaymentPrice + 1,
				DownPayment:     (entities.MinDownPaymentPrice + 1) * entities.MinDownPaymentPercent,
				InterestRate:    10,
				Period:          25,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			isError: true,
		},
		{
			name: "happy path - with CMHC - range 1",
			Data: &entities.MortageData{
				Price:           1000,
				DownPayment:     60,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			mortage: 9.57017323700683,
		},
		{
			name: "happy path - with CMHC - range 2",
			Data: &entities.MortageData{
				Price:           1000,
				DownPayment:     110,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			mortage: 8.982707917906197,
		},
		{
			name: "happy path - with CMHC - range 3",
			Data: &entities.MortageData{
				Price:           1000,
				DownPayment:     160,
				InterestRate:    10,
				Period:          5,
				SchedulePeriods: entities.AcceleratedBiWeekly,
			},
			mortage: 8.45339197383402,
		},
	}

	for _, test := range tt {
		response, err := controller.GetMortage(ctx, test.Data)
		if test.isError {
			if err == nil {
				t.Errorf("%s %s, expected error, got nil", result_prefix, test.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s %s, got error, %v", result_prefix, test.name, err)
				continue
			}

			if response.Payment != test.mortage {
				t.Errorf("%s %s, expected %f, got %f", result_prefix, test.name, test.mortage, response.Payment)
			}
		}
	}

}
