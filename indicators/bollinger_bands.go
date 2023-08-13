package indicators

import (
	"fmt"
)

type BBands struct {
	upper      []float64
	mean       []float64
	lower      []float64
	band_width []float64
}

func BB(input []float64, period int, factor float64, method maType) (BBands, error) {
	err := CheckInput(input, period)
	if err != nil {
		return BBands{}, err
	}
	if factor < 0 {
		return BBands{}, fmt.Errorf("invalid factor: %v", factor)
	}
	res := BBands{mean: []float64{}, lower: []float64{}, upper: []float64{}, band_width: []float64{}}

	res.mean, err = MA(input, period, method)
	if err != nil {
		return BBands{}, err
	}

	res.upper = make([]float64, len(input))
	res.lower = make([]float64, len(input))
	res.band_width = make([]float64, len(input))
	stddev := 0.0
	for i := period - 1; i < len(input); i++ {
		stddev = factor * StdDev(input[i-period+1:i+1])
		res.upper[i] = res.mean[i] + stddev
		res.lower[i] = res.mean[i] - stddev
		res.band_width[i] = res.upper[i] - res.lower[i]
	}

	return res, nil
}
