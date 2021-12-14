package main

import (
	"os/exec"
	"strconv"
	"strings"
)

func updatePrice(stock *Stock) error {
	cmd := exec.Command("tvs", "-n", stock.exchange+":"+stock.symbol, "-s", "-q")

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n") {
		price, err := strconv.ParseFloat(string(line), 64)
		if err != nil {
			return err
		}

		stock.price = price
	}
	return nil
}
