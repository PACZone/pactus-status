package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type PriceExbitronAPI struct {
	LastPrice string `json:"last_price"`
}

func GetPACPrice(priceEndPoint string) float64 {
	prices := make(map[string]map[string]PriceExbitronAPI)

	resp, err := http.Get(priceEndPoint)
	if err != nil {
		log.Println(err)
		return 0
	}

	log.Println(prices)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0
	}

	log.Println(prices)

	err = json.Unmarshal(data, &prices)
	if err != nil {
		log.Println(err)
		return 0
	}

	log.Println(prices)
	price, ok := prices["ticker_name"]["PAC_USDT"]
	if !ok {
		return 0
	}

	num, err := strconv.ParseFloat(price.LastPrice, 64)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return 0
	}

	return num
}

func FormatNumber(num int64) string {
	numStr := strconv.FormatInt(num, 10)

	var formattedNum string
	for i, c := range numStr {
		if (i > 0) && (len(numStr)-i)%3 == 0 {
			formattedNum += ","
		}
		formattedNum += string(c)
	}

	return formattedNum
}
