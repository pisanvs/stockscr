package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {

	// setup tick for every second
	ticker := time.NewTicker(time.Second)

	// initialize configuration
	var config Config
	var stockList []Stock

	ReadConfig(config)

	stockList = config.Stocks

	// initialize termui

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// get terminal width
	maxX, maxY := ui.TerminalDimensions()

	// model ui elements

	p := widgets.NewParagraph()
	p.Title = "COMMANDS"
	p.Text = "> "
	p.SetRect(0, 0, maxX, maxY/10)

	l := widgets.NewList()
	l.SetRect(0, maxY/10, maxX/3, maxY)
	l.Title = "STOCKS"
	for _, stock := range stockList {
		l.Rows = append(l.Rows, stock.symbol)
	}

	ui.Render(p, l)

	// UI Loop

	// MODES:
	// normal mode = 0
	// insert mode = 1
	var mode int = 0

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent {
				if mode == 0 {
					switch e.ID {
					case "n":
						p.Text = "New Stock Symbol: "
						ui.Render(p)
						mode = 1
					case "<C-c>":
						return
					case "q":
						return
					}
				}
				if mode == 1 {
					switch e.ID {
					case "<Escape>":
						p.Text = "> "
						ui.Render(p)
						mode = 0
					case "<Enter>":
						// enable input handling
						mode = 0
						sname := strings.Split(p.Text[18:], ":")
						stockList = append(stockList, Stock{sname[1], sname[0], 0})
						l.Rows = append(l.Rows, sname[1])
						p.Text = "> "
						ui.Render(p, l)
					case "<Backspace>":
						if len(p.Text) < 19 {
							break
						}
						p.Text = p.Text[:len(p.Text)-1]
						ui.Render(p)
					default:
						p.Text = p.Text + e.ID
						ui.Render(p)
					}
				}
			}
		case <-ticker.C:
			// update stock prices
			for i := range stockList {
				go updatePrice(&stockList[i])
				l.Rows[i] = stockList[i].symbol + ": " + strconv.FormatFloat(stockList[i].price, 'f', 2, 64)
				ui.Render(l)
			}
		}
	}
}
