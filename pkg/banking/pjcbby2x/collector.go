package pjcbby2x

import (
	"budget-collector/pkg/csv"
	"budget-collector/pkg/models"
	"budget-collector/pkg/utils/currency"
	"budget-collector/pkg/utils/datetime"
	"errors"
	"log"
	"math"
	"path/filepath"
	"slices"
	"strings"
)

const (
	reportCollectionMask = "reports/*.csv"
)

const (
	periodLabel = "Период выписки:"
)

const (
	operationNameKey     = "Операция"
	operationSumKey      = "Сумма"
	operationDateKey     = "Дата операции по счету"
	operationCategoryKey = "Категория операции"
)

// excluded operations
const (
	topUpOperation            = "CH Payment To Client Contract"
	moneyTransferOperation    = "BLR MINSK MOBILE BANK"
	serviceOperation          = "CH Payment BLR MINSK P2P SDBO NO FEE"
	internalTransferOperation = "CH Debit BLR MINSK P2P SDBO NO FEE"
)

func FindReportByHeaderPeriod(period string) (string, error) {
	periodRange := datetime.GetMonthRangeByPeriod(period)
	reports, err := filepath.Glob(reportCollectionMask)

	if err != nil {
		log.Fatal("Reports not found")
	}

	const startHeaderPosition = 0
	const endHeaderPosition = 15

	for _, reportPath := range reports {
		// read report headers
		records := csv.ReadSlicedCSVFile(reportPath, startHeaderPosition, endHeaderPosition)
		for _, row := range records {
			if slices.Contains(row, periodLabel) && slices.Contains(row, periodRange) {
				return reportPath, nil
			}
		}
	}

	return "", errors.New("report not found")
}

func CollectMonthlyReport(records [][]string) []models.MonthlyReportOperation {

	var paymentMethods []ReportPaymentMethodStat
	var currentPaymentMethod ReportPaymentMethodStat

	for index, value := range records {
		// Header marker
		if len(value) == 1 && strings.Contains(value[0], "Операции по ........") {
			parts := strings.Split(value[0], "........")
			currentPaymentMethod.last4 = parts[len(parts)-1]
			currentPaymentMethod.headerIndex = uint16(index + 1)
			// Last line marker
		} else if len(value) > 0 && strings.Contains(value[0], "Всего по контракту") {
			currentPaymentMethod.lastLineIndex = uint16(index - 1)
			paymentMethods = append(paymentMethods, currentPaymentMethod)
			currentPaymentMethod = ReportPaymentMethodStat{}
		}
	}

	var operations []models.MonthlyReportOperation

	for _, pm := range paymentMethods {
		headerMap := make(map[string]int)
		for index, header := range records[pm.headerIndex] {
			headerMap[header] = index
		}

		for i := pm.headerIndex + 1; i <= pm.lastLineIndex; i++ {
			operationName := records[i][headerMap[operationNameKey]]
			operationCost := currency.StrToMoney(records[i][headerMap[operationSumKey]])

			excludedOperations := strings.Contains(operationName, topUpOperation) ||
				strings.Contains(operationName, moneyTransferOperation) ||
				strings.Contains(operationName, serviceOperation) ||
				strings.Contains(operationName, internalTransferOperation)

			if operationCost < 0 && !excludedOperations {

				operation := models.MonthlyReportOperation{
					Name:        operationName,
					Date:        records[i][headerMap[operationDateKey]],
					PaymentType: models.Card,
					Category:    CategoryMap[records[i][headerMap[operationCategoryKey]]],
					Subcategory: "", // TODO
					Cost:        math.Abs(operationCost),
					Currency:    "BYN", // TODO
					Last4:       pm.last4,
				}

				operations = append(operations, operation)
			}
		}
	}

	return operations
}
