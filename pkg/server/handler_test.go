package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"varmijo/mortage-calculator/pkg/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type testQueryData struct {
	Price        string
	DownPayment  string
	InterestRate string
	Period       string
	Schedule     string
}

type testCase struct {
	name          string
	data          *testQueryData
	status        int
	error_message string
	result        float64
}

func TestMortageHandler(t *testing.T) {
	result_prefix := "error running get mortage handler test, running case:"

	logger := logrus.New()

	e := echo.New()
	controller := usecase.New(logger)

	Register(e, controller, logger)

	tt := []testCase{
		{
			name: "happy path - accelerated-bi-weekly",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "500",
				InterestRate: "10",
				Period:       "5",
				Schedule:     "accelerated-bi-weekly",
			},
			status: http.StatusOK,
			result: 4.894728537748993,
		},
		{
			name: "happy path - bi-weekly",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "500",
				InterestRate: "10",
				Period:       "5",
				Schedule:     "bi-weekly",
			},
			status: http.StatusOK,
			result: 5.303275761953747,
		},
		{
			name: "happy path - monthly",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "500",
				InterestRate: "10",
				Period:       "5",
				Schedule:     "monthly",
			},
			status: http.StatusOK,
			result: 10.623522355634174,
		},
		{
			name: "period non multiple of five",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "500",
				InterestRate: "10",
				Period:       "2",
				Schedule:     "monthly",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong period value",
		},
		{
			name: "missing price",
			data: &testQueryData{
				DownPayment:  "500",
				InterestRate: "10",
				Period:       "5",
				Schedule:     "monthly",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong price value",
		},
		{
			name: "missing down payment",
			data: &testQueryData{
				Price:        "1000",
				InterestRate: "10",
				Period:       "5",
				Schedule:     "monthly",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong down payment value",
		},
		{
			name: "missing interest rate",
			data: &testQueryData{
				Price:       "1000",
				DownPayment: "500",
				Period:      "5",
				Schedule:    "monthly",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong interest rate value",
		},
		{
			name: "missing period",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "500",
				InterestRate: "10",
				Schedule:     "monthly",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong period value",
		},
		{
			name: "missing schedule",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "500",
				InterestRate: "10",
				Period:       "5",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong payment schedule value",
		},
		{
			name: "down payment too big",
			data: &testQueryData{
				Price:        "1000",
				DownPayment:  "1500",
				InterestRate: "10",
				Period:       "5",
				Schedule:     "monthly",
			},
			status:        http.StatusBadRequest,
			error_message: "wrong down payment value",
		},
	}

	for _, test := range tt {
		req := httptest.NewRequest(http.MethodGet, getUrl(test.data), nil)

		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		//Check that the correct http code is returned
		if rec.Code != test.status {
			t.Errorf("%s %s, expected http code %d, got %d", result_prefix, test.name, test.status, rec.Code)
			continue
		}

		//Parse successfull value and validates it
		if test.status == http.StatusOK {
			res := &ResponseMortageData{}
			err := json.Unmarshal(rec.Body.Bytes(), res)
			if err != nil {
				t.Errorf("%s %s, error decoding response, %v", result_prefix, test.name, err)
			}

			if res.Payment != test.result {
				t.Errorf("%s %s, expected result %f, got %f", result_prefix, test.name, test.result, res.Payment)
			}
		} else { //Parse error response and validates it
			res := &ErrorMessage{}
			err := json.Unmarshal(rec.Body.Bytes(), res)
			if err != nil {
				t.Errorf("%s %s, error decoding error message, %v", result_prefix, test.name, err)
			}

			if res.Message != test.error_message {
				t.Errorf("%s %s, expected error message %s, got %s", result_prefix, test.name, test.error_message, res.Message)
			}
		}
	}

}

func getUrl(data *testQueryData) string {
	uri := "/api/v1/mortage?payment-schedule=%s&price=%s&down-payment=%s&period=%s&interest-rate=%s"
	return fmt.Sprintf(uri, data.Schedule, data.Price, data.DownPayment, data.Period, data.InterestRate)
}
