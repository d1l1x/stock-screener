package main

import (
	"fmt"
	stock_screener "github.com/d1l1x/stock-screener"
	"github.com/d1l1x/stock-screener/indicators"
	"github.com/d1l1x/stock-screener/providers"
)

func reverse(slice []float64) {
	for i := len(slice)/2 - 1; i >= 0; i-- {
		opp := len(slice) - 1 - i
		slice[i], slice[opp] = slice[opp], slice[i]
	}
}

func SetupFilters() []func(*providers.BarHistory) bool {
	var filters []func(*providers.BarHistory) bool

	filters = append(filters,
		// Stock dropped more than 5%
		func(bars *providers.BarHistory) bool {
			// Indicators assume the most recent value to be the last,
			// so we have to reverse fmp data
			reverse(bars.Close)
			roc := indicators.ROC(bars.Close, 10)
			return roc[len(roc)-1] < 0.0
		})

	return filters
}

func main() {

	assets := []stock_screener.Asset{
		{Symbol: "AAPL"},
		{Symbol: "GOOGL"},
		{Symbol: "MSFT"},
		{Symbol: "PG"},
	}

	filters := SetupFilters()

	scanner := stock_screener.NewScanner(assets, 100, filters, nil)

	err := scanner.Init(providers.Fmp, "./envs/fmp.env")
	if err != nil {
		panic(err)
	}

	watchlist := scanner.Scan()
	fmt.Println(watchlist)
}
