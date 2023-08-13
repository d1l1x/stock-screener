package indicators

import "github.com/d1l1x/stock-screener/providers"

func Last(bars providers.BarHistory, barType providers.BarType) float64 {
	switch barType {
	case providers.Open:
		return bars.Open[len(bars.Open)-1]
	case providers.High:
		return bars.High[len(bars.High)-1]
	case providers.Low:
		return bars.Low[len(bars.Low)-1]
	case providers.Close:
		return bars.Close[len(bars.Close)-1]
	case providers.AdjClose:
		return bars.AdjClose[len(bars.AdjClose)-1]
	case providers.Volume:
		return float64(bars.Volume[len(bars.Volume)-1])
	}
	return 0
}
