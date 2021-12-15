package main

type Stock struct {
	Symbol   string  `json:"symbol"`
	Exchange string  `json:"exchange"`
	Price    float64 `json:"price"`
	Highlow  string  `json:"highlow"`
}

type Config []Stock

type Headline struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}
