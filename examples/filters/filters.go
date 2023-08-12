package main

import (
	"fmt"
	stock_screener "github.com/d1l1x/stock-screener"
	"github.com/d1l1x/stock-screener/providers"
)

func SetupFilters() []func(*providers.BarHistory) bool {
	var filters []func(*providers.BarHistory) bool

	filters = append(filters, func(history *providers.BarHistory) bool {
		return (history.Close[len(history.Close)-1]) < 0
	})

	return filters
}

func main() {

	assets := []stock_screener.Asset{
		{Symbol: "AAPL"},
	}

	filters := SetupFilters()

	scanner := stock_screener.NewScanner(assets, 10, filters, nil)

	err := scanner.Init(providers.Fmp, "../../envs/fmp.env")
	if err != nil {
		panic(err)
	}

	watchlist := scanner.Scan()
	fmt.Println(watchlist)
}
