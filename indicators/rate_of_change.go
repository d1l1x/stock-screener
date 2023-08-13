package indicators

func ROC(input []float64, period int) []float64 {
	roc := make([]float64, len(input))

	for i := period; i < len(input); i++ {
		roc[i] = ((input[i] - input[i-period]) / input[i-period]) * 100
	}

	return roc
}
