package main

func minInt(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func createGenMinRateFunc(rates []int, cnts []int) func() int {
	rates2 := make([]int, 0)
	for i := 0; i < minInt(len(rates), len(cnts)); i++ {
		rate := rates[i]
		cnt := cnts[i]
		for j := 0; j < cnt; j++ {
			rates2 = append(rates2, rate)
		}
	}

	i := 0
	sz := len(rates2)
	return func() int {
		if i >= sz {
			i = 0
		}
		val := rates2[i]
		i++
		return val
	}
}

func createGenSecRateFunc(genMinRateFunc func() int) func() int {
	var secRates []int
	sec := 0
	return func() int {
		if sec >= 60 {
			sec = 0
		}
		if sec == 0 {
			secRates = make([]int, 60)
			minRate := genMinRateFunc()
			secRate := minRate / 60
			remain := minRate % 60
			for i := 0; i < 60; i++ {
				secRates[i] = secRate
			}
			for i := 0; i < remain; i++ {
				secRates[i]++
			}
		}
		secRate := secRates[sec]
		sec++
		return secRate
	}
}
