package providers

import (
	"context"
	"fmt"
	fmp "github.com/spacecodewor/fmpcloud-go"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

type FMPConfig struct {
	APIKey                string `mapstructure:"apikey"`
	RequestLimitPerMinute int    `mapstructure:"requestlimitperminute"`
}

type FmpProvider struct {
	Client  *fmp.APIClient
	Limiter *rate.Limiter
}

func NewFmpProvider(config FMPConfig) *FmpProvider {
	log.Debug("Initializing FMP API client")
	client, err := fmp.NewAPIClient(fmp.Config{
		APIKey:  config.APIKey, // Set Your API Key from site, default: demo
		Debug:   false,         // Set flag for debug request and response, default: false
		Timeout: 60,            // Set timeout for http client, default: 25
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("Initializing rate limiter")
	limiter := rate.NewLimiter(rate.Every(time.Minute/time.Duration(config.RequestLimitPerMinute)), 1)

	return &FmpProvider{Client: client, Limiter: limiter}
}

func (fmp *FmpProvider) GetHistBars(symbol string, period int) (*BarHistory, error) {

	// Wait for a token from the rate limiter
	if err := fmp.Limiter.Wait(context.Background()); err != nil {
		return nil, fmt.Errorf("Not allowed to proceed for symbol %s: %v\n", symbol, err)
	}

	log.Debug("Get historical bars", zap.String("symbol", symbol))
	bars, err := fmp.Client.Stock.DailyLastNDays(symbol, period)
	if err != nil {
		return nil, err
	}
	if len(bars.Historical) >= period {
		history := new(BarHistory)
		for _, bar := range bars.Historical {
			history.Open = append(history.Open, bar.Open)
			history.High = append(history.High, bar.High)
			history.Low = append(history.Low, bar.Low)
			//TODO: Decide whether to use adjusted close or close
			history.Close = append(history.Close, bar.Close)
			history.AdjClose = append(history.Close, bar.AdjClose)
			history.Volume = append(history.Volume, int64(bar.Volume))
		}
		return history, nil
	} else {
		return nil, fmt.Errorf("not enough bars for %v", symbol)
	}
}
