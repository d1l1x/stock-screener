package main

import (
	"fmt"
	stock_screener "github.com/d1l1x/stock-screener"
	providers "github.com/d1l1x/stock-screener/providers"
)

func main() {

	assets := []stock_screener.Asset{
		{Symbol: "AAPL"},
	}

	scanner := stock_screener.NewScanner(assets, 10, nil, nil)

	err := scanner.Init(providers.Fmp, "../../envs/fmp.env")
	if err != nil {
		panic(err)
	}

	watchlist := scanner.Scan()
	fmt.Println(watchlist[0].Symbol)
}
