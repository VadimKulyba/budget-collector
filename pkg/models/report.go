package models

import (
	"budget-collector/pkg/utils/currency"
)

type PaymentType string

const (
	Cash PaymentType = "наличные"
	Card PaymentType = "карта"
)

type MonthlyReportOperation struct {
	Name        string
	Date        string
	PaymentType PaymentType
	Category    string
	Subcategory string
	Cost        float64
	Currency    currency.Currency
	Last4       string
}
