package indicators

import (
	"github.com/d1l1x/stock-screener/providers"
	"math"
)

func ADX(bars providers.BarHistory, period int) []float64 {
	if len(bars.Close) < period {
		return nil
	}

	adip, adim := DiPlusMinus(bars, period)

	dx := 0.0
	for i := period; i < 2*period; i++ {
		dx += 100 * math.Abs(adip[i]-adim[i]) / (adip[i] + adim[i])
	}
	dx /= float64(period)

	adx := make([]float64, len(bars.Close))
	adx[2*period-1] = dx
	for i := 2 * period; i < len(bars.Close); i++ {
		adx[i] = (adx[i-1]*float64(period-1) + 100*math.Abs(adip[i]-adim[i])/(adip[i]+adim[i])) / float64(period)
	}

	return adx
}

func DiPlusMinus(bars providers.BarHistory, period int) ([]float64, []float64) {
	if len(bars.Close) < period {
		return nil, nil
	}

	dp := make([]float64, len(bars.Close))
	dm := make([]float64, len(bars.Close))

	tr := TR(bars)

	for i := 1; i < len(bars.Close); i++ {

		upMove := bars.High[i] - bars.High[i-1]
		downMove := bars.Low[i-1] - bars.Low[i]

		if upMove > downMove && upMove > 0 {
			dp[i] = upMove
			dm[i] = 0
		}
		if downMove > upMove && downMove > 0 {
			dp[i] = 0
			dm[i] = downMove
		}
	}

	dip := 0.0
	dim := 0.0
	trp := 0.0

	for i := 0; i < period; i++ {
		dip += dp[i]
		dim += dm[i]
		trp += tr[i]
	}

	pdp := dip
	pdm := dim
	ptr := trp

	for i := period; i < len(bars.Close); i++ {
		adp := pdp - pdp/float64(period) + dp[i]
		adm := pdm - pdm/float64(period) + dm[i]
		atr := ptr - ptr/float64(period) + tr[i]

		dp[i] = math.Round(100.0 * adp / atr)
		dm[i] = math.Round(100.0 * adm / atr)

		pdp = adp
		pdm = adm
		ptr = atr
	}

	return dp, dm

}
