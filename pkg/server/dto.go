package server

type GetMortageData struct {
	Price       float64 `schema:"price"`
	DownPayment float64 `schema:"down_payment"`
	InteresRate float64 `schema:"interest_rate"`
	Period      uint    `schema:"period"`
	Schedule    string  `schema:"payment_schedule"`
}

type ResponseMortageData struct {
	Payment   float64 `json:"payment"`
	Insurance float64 `json:"insurance"`
}
