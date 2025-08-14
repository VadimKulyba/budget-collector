// Package csv provides utilities for reading and writing CSV files.
//
// This package handles CSV operations with specific configurations for
// European-style CSV formatting, including semicolon separators and
// proper file handling with error management.
package csv

import (
	"encoding/csv"
	"log"
	"os"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// getBasicCSVReader creates a configured CSV reader for Windows-1251 encoded files.
//
// This internal function sets up a CSV reader with European formatting preferences:
// semicolon separators, flexible field count, and Windows-1251 character encoding
// support for Cyrillic text commonly found in Belarusian bank reports.
//
// Args:
//
//	filePath: string path to the CSV file to read
//
// Returns:
//
//	*csv.Reader: configured CSV reader with European settings
//	*os.File: opened file handle (caller must close)
//
// Note: This function will call log.Fatal() if the file cannot be opened,
// which will terminate the program.
func getBasicCSVReader(filePath string) (*csv.Reader, *os.File) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file: "+filePath, err)
	}

	transformReader := transform.NewReader(file, charmap.Windows1251.NewDecoder())

	csvReader := csv.NewReader(transformReader)
	csvReader.FieldsPerRecord = -1
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true
	return csvReader, file
}

// ReadAllCSVFile reads the entire contents of a CSV file into memory.
//
// This function reads all records from a CSV file and returns them as a 2D string slice.
// It's designed for processing complete bank reports where all data needs to be analyzed.
// The function automatically handles file opening, reading, and closing.
//
// Args:
//
//	filePath: string path to the CSV file to read
//
// Returns:
//
//	[][]string: slice containing all CSV records (each inner slice is a row)
//
// Example:
//
//	records := ReadAllCSVFile("bank_report.csv")
//	for _, row := range records {
//		fmt.Println("Row:", row)
//	}
//
// Note: This function will call log.Fatal() if the file cannot be parsed as CSV,
// which will terminate the program. Consider using ReadAllCSVFileSafe for error handling.
func ReadAllCSVFile(filePath string) [][]string {
	csvReader, file := getBasicCSVReader(filePath)
	defer file.Close()

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

// ReadSlicedCSVFile reads a specific range of records from a CSV file.
//
// This function reads only a subset of records from a CSV file, useful for
// processing large files in chunks or reading specific sections like headers.
// It's particularly helpful when you only need to examine file metadata or
// process files in streaming fashion.
//
// Args:
//
//	filePath: string path to the CSV file to read
//	start:    uint16 starting record index (0-based)
//	end:      uint16 ending record index (exclusive)
//
// Returns:
//
//	[][]string: slice containing the requested range of CSV records
//
// Example:
//
//	// Read only the first 20 rows (headers and first few data rows)
//	headerRecords := ReadSlicedCSVFile("bank_report.csv", 0, 20)
//
//	// Read rows 100-200 for batch processing
//	batchRecords := ReadSlicedCSVFile("large_report.csv", 100, 200)
//
// Note: This function will call log.Fatal() if the requested range exceeds
// the file length, which will terminate the program.
func ReadSlicedCSVFile(filePath string, start uint16, end uint16) [][]string {
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
