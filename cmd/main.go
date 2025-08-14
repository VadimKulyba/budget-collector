package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const (
	reportCollectionMask   = "reports/*.csv"
	periodDateFormat       = "01.2006"
	periodBorderDateFormat = "02.01.2006"
	periodLabel            = "Период выписки:"
)

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
	name        string
	date        string
	paymentType PaymentType
	category    string
	subcategory string
	cost        float64
	currency    Currency
	last4       string
}

var CategoryMap = map[string]string{
	"Магазины продуктовые":     "Продукты и заморозка",
	"Ресторация / бары / кафе": "Кафе и доставки",
	"Аптеки":                           "Здоровье",
	"Медицинский сервис":               "Здоровье",
	"Магазины одежды":                  "Шопинг",
	"Различные магазины":               "Шопинг",
	"Прочее":                           "Шопинг",
	"Поставщик  услуг":                 "Здоровье",
	"Транспорт - Такси":                "Транспорт",
	"Транспортировка":                  "Транспорт",
	"Развлечения - кино":               "Развлечения",
	"Развлечения":                      "Развлечения",
	"Коммунальные услуги":              "Подписки",
	"Индивидуальные сервис провайдеры": "Спорт",
}

func getBasicCSVReader(filePath string) (*csv.Reader, *os.File) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}

	transformReader := transform.NewReader(file, charmap.Windows1251.NewDecoder())

	csvReader := csv.NewReader(transformReader)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	return csvReader, file
}

func readAllCSVFile(filePath string) [][]string {
	csvReader, file := getBasicCSVReader(filePath)
	defer file.Close()

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func readSlicedCSVFile(filePath string, start uint16, end uint16) [][]string {
	csvReader, file := getBasicCSVReader(filePath)
	defer file.Close()

	var records [][]string
	for i := start; i < end; i++ {
		record, err := csvReader.Read()
		if err != nil {
			log.Fatal("Unsupported file length "+filePath, err)
		}

		records = append(records, record)
	}

	return records
}

func writeDataToCSVFile(records [][]string) {
	file, err := os.Create("output.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
	}

	fmt.Println("CSV file saved as output.csv")
}

func getMonthRangeByPeriod(period string) (time.Time, time.Time) {
	// validate period
	startOfPeriod, err := time.Parse(periodDateFormat, period)
	if err != nil {
		log.Fatal("Invalid period format")
	}

	nextMonth := time.Date(
		startOfPeriod.Year(),
		startOfPeriod.Month()+1,
		1,
		0, 0, 0, 0,
		startOfPeriod.Location(),
	)

	endOfPeriod := nextMonth.Add(-24 * time.Hour)

	return startOfPeriod, endOfPeriod
}

func findReportByHeaderPeriod(period string) (string, error) {
	startOfPeriod, endOfPeriod := getMonthRangeByPeriod(period)
	reports, err := filepath.Glob(reportCollectionMask)

	if err != nil {
		log.Fatal("Reports not found")
	}

	periodValue := startOfPeriod.Format(periodBorderDateFormat) + "-" + endOfPeriod.Format(periodBorderDateFormat)

	const startHeaderPosition = 0
	const endHeaderPosition = 15

	for _, reportPath := range reports {
		// read report headers
		records := readSlicedCSVFile(reportPath, startHeaderPosition, endHeaderPosition)
		for _, row := range records {
			if slices.Contains(row, periodLabel) && slices.Contains(row, periodValue) {
				return reportPath, nil
			}
		}
	}

	return "", errors.New("report not found")
}

type ReportPaymentMethodStat struct {
	last4         string
	headerIndex   uint16
	lastLineIndex uint16
}

func strToMoney(str string) float64 {
	res, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(str, ",", "."), " ", ""), 64)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

func moneyToStr(amount float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", amount), ".", ",")
}

func cleanUpQuotes(str string) string {
	return strings.ReplaceAll(str, `"`, "")
}

func collectMonthlyReport(records [][]string) []MonthlyReportOperation {

	var paymentMethods []ReportPaymentMethodStat
	var currentPaymentMethod ReportPaymentMethodStat

	for index, value := range records {
		if len(value) == 1 && strings.Contains(value[0], "Операции по ........") {
			parts := strings.Split(value[0], "........")
			currentPaymentMethod.last4 = parts[len(parts)-1]
			currentPaymentMethod.headerIndex = uint16(index + 1)
		} else if len(value) > 0 && strings.Contains(value[0], "Всего по контракту") {
			currentPaymentMethod.lastLineIndex = uint16(index - 1)
			paymentMethods = append(paymentMethods, currentPaymentMethod)
			currentPaymentMethod = ReportPaymentMethodStat{}
		}
	}

	var operations []MonthlyReportOperation

	for _, pm := range paymentMethods {
		headerMap := make(map[string]int)
		for index, header := range records[pm.headerIndex] {
			headerMap[header] = index
		}

		for i := pm.headerIndex + 1; i <= pm.lastLineIndex; i++ {
			operationName := records[i][headerMap["Операция"]]
			operationCost := strToMoney(records[i][headerMap["Сумма"]])

			if operationCost < 0 && !strings.Contains(operationName, "BLR MINSK MOBILE BANK") && !strings.Contains(operationName, "CH Payment To Client Contract") && !strings.Contains(operationName, "CH Debit BLR MINSK P2P SDBO NO FEE") {

				operation := MonthlyReportOperation{
					name:        operationName,
					date:        records[i][headerMap["Дата операции по счету"]],
					paymentType: Card,
					category:    CategoryMap[records[i][headerMap["Категория операции"]]],
					subcategory: "",
					cost:        math.Abs(operationCost),
					currency:    "BYN",
					last4:       pm.last4,
				}

				operations = append(operations, operation)
			}
		}
	}

	return operations
}

func main() {
	var reportPeriod string

	fmt.Println("Please enter [MM.YYYY] report period:")
	fmt.Scanln(&reportPeriod)

	reportPath, err := findReportByHeaderPeriod(reportPeriod)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Found report: " + reportPath)
	}

	records := readAllCSVFile(reportPath)
	monthlyOperations := collectMonthlyReport(records)

	var results [][]string
	for _, operation := range monthlyOperations {
		results = append(results, []string{
			operation.date,
			string(operation.paymentType),
			operation.category,
			operation.subcategory,
			moneyToStr(operation.cost),
			"", // empty column
			cleanUpQuotes(operation.name),
		})
	}

	writeDataToCSVFile(results)
}
