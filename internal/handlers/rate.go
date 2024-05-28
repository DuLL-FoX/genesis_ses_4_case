package handlers

import (
	"awesomeProject/internal/email"
	"encoding/json"
	"fmt"
	"net/http"
)

const nbuAPI = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=USD&json"

type ExchangeRate struct {
	Rate float64 `json:"rate"`
}

var GetExchangeRate = func() (float64, error) {
	req, err := http.NewRequest("GET", nbuAPI, nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var exchangeRates []ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRates); err != nil {
		return 0, err
	}

	if len(exchangeRates) == 0 {
		return 0, fmt.Errorf("no exchange rate data received")
	}

	return exchangeRates[0].Rate, nil
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	rate, err := GetExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Current USD to UAH exchange rate: %.2f", rate)
	email.SendEmails(rate) // For testing purposes only
}
