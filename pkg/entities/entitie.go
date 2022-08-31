package entities

//Defines a payment schedule type
type PaymentSchedule uint8

const (
	AcceleratedBiWeekly PaymentSchedule = 26
	BiWeekly            PaymentSchedule = 24
	Montly              PaymentSchedule = 12
)

type MortageData struct {
	Price, DownPayment, InterestRate float64
	Period                           uint
	SchedulePeriods                  PaymentSchedule
}

type ResultData struct {
	Payment   float64
	Insurance float64
}

const (
	MinDownPaymentPercent = 0.05 // Minimun percentege for down payment percent for prices below 500000

	MaxDownPaymentPercent = 0.2 // CMHC not required if above this down payment percentege

	MaxCMHCPrice = 1000000 // Maximum price that allows CMHC

	MaxCMHCPeriod = 25 // Maximun period that allows CMHC

	MinDownPaymentPrice = 500000 // Minimum value that allows MinDownPaymentPercent

	ComplementDownPaymentPercent = 0.1 // Down payment percentage for the complement values for prices over 500000
)
