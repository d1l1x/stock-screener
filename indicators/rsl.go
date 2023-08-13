package indicators

func RSL(input []float64, period int) ([]float64, error) {
	err := CheckInput(input, period)
	if err != nil {
		return nil, err
	}
	res := make([]float64, len(input))
	for i := period; i < len(input); i++ {
		res[i] = input[i] / Mean(input[i-period:i+1])
	}
	return res, nil
}
