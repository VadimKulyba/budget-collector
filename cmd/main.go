package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
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

	for _, reportPath := range reports {
		// read report headers
		records := readSlicedCSVFile(reportPath, 0, 15)
		for _, row := range records {
			if slices.Contains(row, periodLabel) && slices.Contains(row, periodValue) {
				return reportPath, nil
			}
		}
	}

	return "", errors.New("report not found")
}

func main() {
	var reportPeriod string

	fmt.Println("Please enter [MM.YYYY] report period:")
	fmt.Scanln(&reportPeriod)

	reportPath, err := findReportByHeaderPeriod(reportPeriod)

	if err != nil {
		log.Fatal(err)
	}

	// records := readSlicedCSVFile("/Users/vk/Projects/pf_reports/budget-collector/Vpsk_72482430.csv", 0, 15)
	fmt.Println(reportPath)
}
