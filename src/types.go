package main

type Stock struct {
	symbol   string
	exchange string
	price    float64
}

type Config struct {
	Stocks []Stock
}
