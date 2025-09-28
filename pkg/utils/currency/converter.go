// Package currency provides currency conversion functionality using the National Bank of Belarus API.
// It supports converting foreign currencies to Belarusian Rubles (BYN) based on official exchange rates.
package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	// baseNbRbApiURL is the base URL for the National Bank of Belarus exchange rates API
	baseNbRbApiURL = "https://api.nbrb.by/exrates"
	// baseNbRbApiDateFormat is the date format used by the NBRB API
	baseNbRbApiDateFormat = "2006-01-02"
)

// nbRbCurrencyCodeMap maps currency types to their corresponding NBRB API currency codes
var nbRbCurrencyCodeMap = map[Currency]int{
	USD: 431,
	RUB: 456,
}

// NbRbExchangeRate represents the exchange rate data structure returned by the NBRB API
type NbRbExchangeRate struct {
	Cur_Scale        int     // Currency scale (e.g., 1 for USD, 100 for RUB)
	Cur_OfficialRate float64 // Official exchange rate
}

// GetCurrencyRate retrieves the exchange rate for converting from one currency to another.
// Currently supports conversion to Belarusian Rubles (BYN) using NBRB official rates.
// Returns the exchange rate and any error that occurred during the API call.
func GetCurrencyRate(fromCurrency Currency, toCurrency Currency, date time.Time) (float64, error) {
	if toCurrency == BYN {
		rate := getBYNRateByDate(fromCurrency, date)
		return rate.Cur_OfficialRate / float64(rate.Cur_Scale), nil
	} else {
		return 0, errors.New("unsupported currency")
	}
}

// getBYNRateByDate fetches the exchange rate data from NBRB API for the specified currency and date.
// It makes an HTTP GET request to the NBRB API and parses the JSON response into NbRbExchangeRate struct.
// Returns the complete exchange rate data structure.
func getBYNRateByDate(toCurrency Currency, date time.Time) NbRbExchangeRate {
	params := url.Values{}
	params.Add("ondate", date.Format(baseNbRbApiDateFormat))
	currencyCode := nbRbCurrencyCodeMap[toCurrency]
	nbrbExchangeURL := baseNbRbApiURL + "/rates" + fmt.Sprintf("/%d", currencyCode) + "?" + params.Encode()

	resp, err := http.Get(nbrbExchangeURL)
	if err != nil {
		log.Fatalf("Error making GET nbrb rate request: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading nbrb rate response body: %v", err)
	}

	var rate NbRbExchangeRate

	err = json.Unmarshal(body, &rate)
	if err != nil {
		log.Fatalf("Error parsing nbrb rate response: %v", err)
	}

	return rate
}
