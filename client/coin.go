package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const BINANCE_API_URL = "https://api.binance.com/api/v3/ticker/price?symbol=%s"

type CoinPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func GetCoinPrice(symbol string) (*CoinPrice, error) {
	symbolUpper := strings.ToUpper(symbol)
	resp, err := http.Get(fmt.Sprintf(BINANCE_API_URL, symbolUpper))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	coinPrice := new(CoinPrice)
	err = json.Unmarshal(bytes, coinPrice)
	if err != nil {
		return nil, err
	}
	return coinPrice, nil
}
