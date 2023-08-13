package indicators

import (
	"github.com/d1l1x/stock-screener/providers"
	"math"
)

func ATR(bars providers.BarHistory, period int) []float64 {
	if len(bars.Close) < period || len(bars.High) < period || len(bars.Low) < period {
		return nil
	}
	tr := TR(bars)
	atr, _ := MA(tr, period, WILDER)
	return atr
}

func TR(bars providers.BarHistory) []float64 {
	tr := make([]float64, len(bars.Close))
	tr[0] = bars.High[0] - bars.Low[0]
	for i := 1; i < len(bars.Close); i++ {
		highLow := bars.High[i] - bars.Low[i]
		highClose := math.Abs(bars.High[i] - bars.Close[i-1])
		lowClose := math.Abs(bars.Low[i] - bars.Close[i-1])
		tr[i] = math.Max(highLow, math.Max(highClose, lowClose))
	}
	return tr
}

func ATRP(bars providers.BarHistory, period int) []float64 {
	atr := ATR(bars, period)
	res := make([]float64, len(atr))
	for i, val := range atr {
		res[i] = val / bars.Close[i] * 100
	}
	return res
}
