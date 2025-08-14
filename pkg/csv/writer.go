// Package csv provides utilities for reading and writing CSV files.
//
// This package handles CSV operations with specific configurations for
// European-style CSV formatting, including semicolon separators and
// proper file handling with error management.
package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// WriteDataToCSVFile writes a 2D string slice to a CSV file with European formatting.
//
// This function creates a new CSV file (or overwrites existing) and writes all records
// using semicolon (;) as the field separator, which is the standard for European CSV files.
// The function automatically handles file creation, writing, and cleanup.
//
// Args:
//
//	filePath: string path where the CSV file should be created/written
//	records:  [][]string slice containing the data to write (each inner slice is a row)
//
// Example:
//
//	records := [][]string{
//		{"Date", "Category", "Amount"},
//		{"01.01.2024", "Groceries", "50,00"},
//		{"02.01.2024", "Transport", "25,50"},
//	}
//	WriteDataToCSVFile("expenses.csv", records)
//	// Creates expenses.csv with semicolon-separated values
//
// Note: This function will call log.Fatal() if the file cannot be created,
// which will terminate the program. Consider using WriteDataToCSVFileSafe for error handling.
func WriteDataToCSVFile(filePath string, records [][]string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Unable to create file: "+filePath, err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	for _, record := range records {
		if err := csvWriter.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
	}

	fmt.Println("CSV file saved as", filePath)
}
