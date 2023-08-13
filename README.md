[![Build Status](https://github.com/d1l1x/stock-screener/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/d1l1x/stock-screener/actions/workflows/go.yml)
[![CodeQL](https://github.com/d1l1x/stock-screener/actions/workflows/codeql.yml/badge.svg)](https://github.com/d1l1x/stock-screener/actions/workflows/codeql.yml)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=d1l1x_stock-screener&metric=coverage)](https://sonarcloud.io/summary/new_code?id=d1l1x_stock-screener)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=d1l1x_stock-screener&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=d1l1x_stock-screener)

# Stock screener

Stock screener is a tool that scans through a list of assets and returns a list of assets that
meet a predefined set of filter criteria. It also allows to rank results.

## Description

While you might use the stock screener as a standalone tool (for example by extending one of the examples) it is
meant to be used as a package to facilitate automatic trading. The idea then is to use the stock screener to find
the most promising candidates for trading, according to strategy specific filter criteria.

Setting up a scanner is as easy as
```go
scanner := stock_screener.NewScanner(assets, 10, filters, ranking)
```

where 
* `assets` is a list of assets to scan, 
* `10` is the look-back period, i.e. the number of bars used to compute filters
* `filters` is a list of filter functions to apply to the assets
* `ranking` is a ranking function to apply to the filtered assets

## Filters
Filter functions must have the following signature
```go
func(*providers.BarHistory) bool
```
where `providers.BarHistory` is a struct that contains the historical data for a given asset.

## Ranking
Ranking functions must have the following signature
```go
func(*[]stock_screener.Asset)
```
where `[]stock_screener.Asset` is the list of assets that passed the filters.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Go >= `1.20.5`
- A valid subscription for any of the supported providers
  - [FMP](https://financialmodelingprep.com/developer/docs/)

## Installation

```shell
$> go get -u github.com/d1l1x/stock-screener
```

## Examples

In order to test any of the examples you require to have a valid subscriptions for any of the supported providers. It might
be necessary to adjust the code to use another provider.

### Simple example using FMP

This simple example sets up a scanner with a single asset and scans it using the FMP provider.
The scanner will return the same list of assets as there are not filters defined.
```go
package main

import (
	"fmt"
	stock_screener "github.com/d1l1x/stock-screener"
	providers "github.com/d1l1x/stock-screener/providers"
)

func main() {

	assets := []stock_screener.Asset{{Symbol: "AAPL"}}

	scanner := stock_screener.NewScanner(assets, 10, nil, nil)

	err := scanner.Init(providers.Fmp, "../../envs/fmp.env")
	if err != nil {
		panic(err)
	}

	watchlist := scanner.Scan()
	fmt.Println(watchlist[0].Symbol)
}
```

### Advanced example using filters and ranking

Like for the simple example, this example uses FMP, sets up a scanner but this time with a list of three
assets, scans it applying the provided filters and finally ranks it according to the given ranking
function.

```go
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
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[BSD 3-Clause License](LICENSE)
