package client

import (
	"testing"
)

func TestGetCoinPrice(t *testing.T) {
	price, err := GetCoinPrice("solusdt")
	if err != nil {
		t.Error(err)
	}
	t.Log(price)
}
