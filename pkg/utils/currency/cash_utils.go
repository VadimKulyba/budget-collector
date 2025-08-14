// Package currency provides utilities for handling currency operations.
//
// This file handles the conversion between string representations of money
// (with European comma decimal separators) and float64 values used for calculations.
package currency

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// floatSize defines the bit size for float parsing operations
const floatSize = 64

// StrToMoney converts a string representation of money to a float64 value.
//
// This function handles European-style number formatting where commas are used
// as decimal separators instead of periods. It also removes any spaces from
// the input string before parsing.
//
// Args:
//
//	str: string representation of money amount (e.g., "1 234,56" or "1234,56")
//
// Returns:
//
//	float64: parsed monetary value
//
// Example:
//
//	amount := StrToMoney("1 234,56")
//	// Returns: 1234.56
//
//	amount2 := StrToMoney("99,99")
//	// Returns: 99.99
//
// Note: This function will call log.Fatal() if the string cannot be parsed,
// which will terminate the program. Consider using StrToMoneySafe for error handling.
func StrToMoney(str string) float64 {
	res, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(str, ",", "."), " ", ""), floatSize)

	if err != nil {
		log.Fatal(err)
	}

	return res
}

// MoneyToStr converts a float64 monetary value to a formatted string.
//
// This function formats monetary values using European-style formatting with
// commas as decimal separators and exactly 2 decimal places.
//
// Args:
//
//	amount: float64 monetary value to format
//
// Returns:
//
//	string: formatted monetary string with comma decimal separator
//
// Example:
//
//	formatted := MoneyToStr(1234.56)
//	// Returns: "1234,56"
//
//	formatted2 := MoneyToStr(99.9)
//	// Returns: "99,90"
//
//	formatted3 := MoneyToStr(100)
//	// Returns: "100,00"
func MoneyToStr(amount float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", amount), ".", ",")
}
