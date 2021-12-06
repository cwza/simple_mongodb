package main

import (
	"fmt"
	"os"
	"testing"
)

func sumInts(as []int) int {
	s := 0
	for _, a := range as {
		s += a
	}
	return s
}

func TestInitConfig(t *testing.T) {
	os.Setenv("CONSUMERURL", "xxxx")
	config, err := initConfig("./config.toml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("config: %+v\n", config)
}

func TestCreateGenMinRateFunc(t *testing.T) {
	// rates := []int{6000, 12000, 18000, 24000, 30000, 24000, 18000, 12000, 6000, 0}
	// cnts := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	rates := []int{10, 20}
	cnts := []int{5, 2}
	genMinRateFunc := createGenMinRateFunc(rates, cnts)

	cnt := 100
	rates2 := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		rates2[i] = genMinRateFunc()
	}
	fmt.Printf("rates2: %v\n", rates2)
}

func TestCreateGenSecRateFunc2(t *testing.T) {
	// rates := []int{6000, 12000, 18000, 24000, 30000, 24000, 18000, 12000, 6000, 0}
	// cnts := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	rates := []int{100, 15000}
	cnts := []int{10, 10}
	genMinRateFunc := createGenMinRateFunc(rates, cnts)

	genSecRateFunc := createGenSecRateFunc(genMinRateFunc)
	cnt := 10 * 60 * 2
	rates2 := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		rates2[i] = genSecRateFunc()
	}
	for i := 0; i < cnt/60; i++ {
		tmp := rates2[i*60 : (i+1)*60]
		fmt.Printf("%v ", tmp)
		fmt.Printf("%d ", sumInts(tmp))
		fmt.Println()
	}
}
