package main

import (
	"fmt"
	stock_screener "github.com/d1l1x/stock-screener"
	"github.com/d1l1x/stock-screener/providers"
)

func SetupFilters() []func(*providers.BarHistory) bool {
	var filters []func(*providers.BarHistory) bool

	// Stupid filter that checks if the last close is greater than 0
	filters = append(filters, func(history *providers.BarHistory) bool {
		return (history.Close[len(history.Close)-1]) > 0
	})

	return filters
}

func SetupRanking() func(*[]stock_screener.Asset) {
	// Stupid ranking that just changes the order of the assets
	return func(assets *[]stock_screener.Asset) {
		(*assets)[0].Symbol = "GOOGL"
		(*assets)[1].Symbol = "MSFT"
		(*assets)[2].Symbol = "AAPL"
	}
}

func main() {

	assets := []stock_screener.Asset{
		{Symbol: "AAPL"},
		{Symbol: "MSFT"},
		{Symbol: "GOOGL"},
	}

	filters := SetupFilters()
	ranking := SetupRanking()

	scanner := stock_screener.NewScanner(assets, 10, filters, ranking)

	err := scanner.Init(providers.Fmp, "../../envs/fmp.env")
	if err != nil {
		panic(err)
	}

	watchlist := scanner.Scan()
	fmt.Println(watchlist)
}
