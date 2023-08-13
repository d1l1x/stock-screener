package stock_screener

import (
	"github.com/d1l1x/stock-screener/providers"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log, _ = zap.NewProduction()

//var log = zap.Must(zap.NewDevelopment())

type Scanner struct {
	Filters        []func(bars *providers.BarHistory) bool
	Ranking        func([]Asset)
	Assets         []Asset
	LookBackPeriod int
	provider       providers.DataProvider
}

type Asset struct {
	Symbol string
	Name   string
	Id     string
}

func NewScanner(assets []Asset, lookBackPeriod int, filters []func(*providers.BarHistory) bool, ranking func([]Asset)) *Scanner {
	return &Scanner{
		Assets:         assets,
		Ranking:        ranking,
		LookBackPeriod: lookBackPeriod,
		Filters:        filters,
	}
}

func (s *Scanner) Init(provider providers.Provider, providerConfigPath string) error {

	log.Info("Init scanner")
	if s.Filters == nil {
		log.Warn("No filters defined")
	}

	viper.SetConfigFile(providerConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading provider config file", zap.Error(err))
	}

	var dp providers.DataProvider

	switch provider {
	case providers.Fmp:
		var config providers.FMPConfig
		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal("", zap.Error(err))
		}
		dp = providers.NewFmpProvider(config)
	default:
		log.Fatal("Unknown data provider")
	}
	s.provider = dp

	return nil
}

func (s *Scanner) Scan() []Asset {
	defer log.Sync()
	log.Info("Start scanning")

	errChan := make(chan error, len(s.Assets))

	var assetsToConsider []Asset

	// filter all assets
	for i, asset := range s.Assets {
		go func(idx int, a Asset) {
			bars, err := s.provider.GetHistBars(a.Symbol, s.LookBackPeriod)
			if err != nil {
				errChan <- err
				return
			}
			passed := true
			if s.Filters != nil {
				for _, filter := range s.Filters {
					if !filter(bars) {
						passed = false
						break
					}
				}
			}
			if passed {
				assetsToConsider = append(assetsToConsider, s.Assets[idx])
			}
			errChan <- nil
		}(i, asset)
	}
	// Check possible channel errors
	for range s.Assets {
		if err := <-errChan; err != nil {
			log.Error("Channel error", zap.Error(err))
		}
	}

	if s.Ranking != nil {
		log.Info("Apply ranking")
		s.Ranking(assetsToConsider)
	}

	return assetsToConsider
}
