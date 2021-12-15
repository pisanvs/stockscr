package main

import (
	"encoding/json"
	"net/http"
)

func GetNews(fullsymbol string) ([]Headline, error) {
	var news []Headline

	req, err := http.NewRequest("GET", "https://news-headlines.tradingview.com/headlines/?client=web&lang=en&locale=en&proSymbol="+fullsymbol, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.tradingview.com/")
	req.Header.Set("Origin", "https://www.tradingview.com")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&news)
	if err != nil {
		return nil, err
	}

	return news, nil
}
