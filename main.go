package main

import (
	"encoding/csv"
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
	reportsMask      = "reports/*.csv"
	periodDateFormat = "01.2006"
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

func findReportByHeaderDate(period string) {
	files, err := filepath.Glob(reportsMask)

	if err != nil {
		log.Fatal("Reports not found")
	}

	for _, file := range files {
		// read report headers
		records := readSlicedCSVFile(file, 0, 15)
		for _, row := range records {
			if slices.Contains(row, "Период выписки:") {
				fmt.Println(row)
			}
		}
	}
}

func main() {
	var reportPeriod string

	fmt.Println("Please enter [MM.YYYY] report period:")
	fmt.Scanln(&reportPeriod)

	// validation
	_, err := time.Parse(periodDateFormat, reportPeriod)
	if err != nil {
		log.Fatal("Invalid period format")
	}

	findReportByHeaderDate(reportPeriod)

	// records := readSlicedCSVFile("/Users/vk/Projects/pf_reports/budget-collector/Vpsk_72482430.csv", 0, 15)
	// fmt.Println(records)
}
