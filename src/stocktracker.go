package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func updatePrice(stock *Stock, pastPrice float64) error {
	// get stock price
	var symbolfull string
	if stock.Exchange == "" {
		symbolfull = stock.Symbol
	} else {
		symbolfull = stock.Exchange + ":" + stock.Symbol
	}
	cmd := exec.Command("tvs", "-n", symbolfull, "-s", "-q")

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n") {
		price, err := strconv.ParseFloat(string(line), 64)
		if err != nil {
			return err
		}
		stock.Price = price
		if price > pastPrice {
			stock.Highlow = " ▲"
		} else if price < pastPrice {
			stock.Highlow = " ▼"
		} else {
			stock.Highlow = " -"
		}

	}
	return nil
}
