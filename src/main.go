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
	updTick := time.NewTicker(time.Second * 5)

	// initialize configuration
	var stockList []Stock

	ReadConfig(&stockList)

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
	l.SetRect(0, maxY/10, maxX/4, maxY)
	l.Title = "STOCKS"
	for _, stock := range stockList {
		l.Rows = append(l.Rows, stock.Symbol)
	}
	l.SelectedRowStyle.Bg = ui.ColorGreen

	hl := widgets.NewList()
	hl.SetRect(maxX/4, maxY/10, maxX, maxY/4)
	hl.Title = "HEADLINES"
	hl.SelectedRowStyle.Bg = ui.ColorGreen

	ui.Render(p, l, hl)

	// UI Loop

	// MODES:
	// normal mode = 0
	// insert mode = 1
	var mode int = 0
	var selected int = 0
	uiloop(mode, selected, p, l, hl, stockList, updTick)
}

func uiloop(mode int, selected int, p *widgets.Paragraph, l *widgets.List, hl *widgets.List, stockList []Stock, updTick *time.Ticker) bool {
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
						// save config
						for i := range stockList {
							stockList[i].Headlines = nil
						}
						SaveConfig("./config.json", stockList)
						return true
					case "q":
						// save config
						for i := range stockList {
							stockList[i].Headlines = nil
						}
						SaveConfig("./config.json", stockList)
						return true
					case "j":
						if len(l.Rows) != 0 {
							if selected == 0 {
								hl.Rows = []string{}
								// only increment selected row if we are not at the start of the list
								if l.SelectedRow < len(l.Rows)-1 {
									l.SelectedRow++
								} else {
									l.SelectedRow = 0
								}
								if len(stockList[l.SelectedRow].Headlines) == 0 {
									res, err := GetNews(stockList[l.SelectedRow].Exchange + ":" + stockList[l.SelectedRow].Symbol)
									if err != nil {
										log.Println(err)
									}
									stockList[l.SelectedRow].Headlines = res
								}
								for _, headline := range stockList[l.SelectedRow].Headlines {
									hl.Rows = append(hl.Rows, headline.Title)
								}
							} else {
								if hl.SelectedRow > len(hl.Rows)-1 {
									hl.SelectedRow = 0
								} else {
									hl.SelectedRow++
								}
							}
							ui.Render(l, hl)
						}
					case "k":
						if len(l.Rows) != 0 {
							if selected == 0 {
								hl.Rows = []string{}
								if l.SelectedRow > 0 {
									l.SelectedRow--
								} else {
									l.SelectedRow = len(l.Rows) - 1
								}
								if len(stockList[l.SelectedRow].Headlines) == 0 {
									res, err := GetNews(stockList[l.SelectedRow].Exchange + ":" + stockList[l.SelectedRow].Symbol)
									if err != nil {
										log.Println(err)
									}
									stockList[l.SelectedRow].Headlines = res
								}
								for _, headline := range stockList[l.SelectedRow].Headlines {
									hl.Rows = append(hl.Rows, headline.Title)
								}
							} else {
								if hl.SelectedRow > 0 {
									hl.SelectedRow--
								}
							}
							ui.Render(l, hl)
						}
					case "l":
						if selected == 0 {
							selected = 1
						}
						ui.Render(l, hl)
					case "h":
						if selected == 1 {
							selected = 0
						}
						ui.Render(l, hl)
					}
				}
				if mode == 1 {
					switch e.ID {
					case "<Escape>":
						p.Text = "> "
						ui.Render(p)
						mode = 0
					case "<Enter>":
						mode = 0
						sname := make([]string, 2)
						// check if exchange was specified
						if strings.Contains(p.Text[18:], ":") {
							// split string
							sname = strings.Split(p.Text[18:], ":")
							stockList = append(stockList, Stock{sname[1], sname[0], 0, " -", []Headline{}})
						} else {
							sname[0] = p.Text[18:]
							stockList = append(stockList, Stock{sname[0], "", 0, " -", []Headline{}})
						}
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
		case <-updTick.C:
			for i := range stockList {
				go updatePrice(&stockList[i], stockList[i].Price)
				newPrice := strconv.FormatFloat(stockList[i].Price, 'f', 2, 64)

				l.Rows[i] = stockList[i].Symbol + ": " + newPrice + stockList[i].Highlow
				ui.Render(l)
			}
		}
	}

}
