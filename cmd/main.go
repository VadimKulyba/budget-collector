package main

import (
	"fmt"
	"log"
	"strings"

	"budget-collector/pkg/banking/pjcbby2x"
	"budget-collector/pkg/csv"
	"budget-collector/pkg/models"
	"budget-collector/pkg/utils/currency"
)

func cleanUpQuotes(str string) string {
	return strings.ReplaceAll(str, `"`, "")
}

func main() {
	var reportPeriod string

	fmt.Println("Please enter [MM.YYYY] report period:")
	fmt.Scanln(&reportPeriod)

	reportPaths, reportPathsErr := pjcbby2x.FindReportByHeaderPeriod(reportPeriod)

	if reportPathsErr != nil {
		log.Fatal(reportPathsErr)
	}

	var monthlyOperations []models.MonthlyReportOperation
	for _, reportPath := range reportPaths {
		fmt.Println("Found report: " + reportPath)

		records := csv.ReadAllCSVFile(reportPath)
		monthlyReport := pjcbby2x.CollectMonthlyReport(records)

		monthlyOperations = append(monthlyOperations, monthlyReport...)
	}

	var results [][]string
	for _, operation := range monthlyOperations {
		results = append(results, []string{
			operation.Date,
			string(operation.PaymentType),
			operation.Category,
			operation.Subcategory,
			currency.MoneyToStr(operation.Cost),
			"", // empty column by default
			cleanUpQuotes(operation.Name),
		})
	}

	csv.WriteDataToCSVFile("output.csv", results)
}
