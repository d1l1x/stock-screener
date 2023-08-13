package providers

import (
	"go.uber.org/zap"
)

var log, _ = zap.NewProduction()

type Provider int

const (
	Fmp = iota
)

type BarHistory struct {
	Open     []float64
	High     []float64
	Low      []float64
	Close    []float64
	AdjClose []float64
	Volume   []int64
}

type BarType int

const (
	Open = iota
	High
	Low
	Close
	AdjClose
	Volume
)

type DataProvider interface {
	GetHistBars(symbol string, period int) (*BarHistory, error)
}
