package models

type PaymentType string
type Currency string

const (
	Cash PaymentType = "наличные"
	Card PaymentType = "карта"
)

const (
	BYN Currency = "BYN"
)

type MonthlyReportOperation struct {
	Name        string
	Date        string
	PaymentType PaymentType
	Category    string
	Subcategory string
	Cost        float64
	Currency    Currency
	Last4       string
}
