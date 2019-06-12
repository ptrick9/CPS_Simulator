package main

import (
	"../simulator/cps"
	"sync"
	"testing"
	"time"
)

func TestAddMeasurement(t *testing.T) {
	p := cps.Params{}
	p.NumGridSamples = 5000
	travelList := make([]bool,0)
	travelList = append(travelList, false)

	test_squares := make([]cps.Square, 2)

	test_squares[0] = cps.Square{0, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	test_squares[1] = cps.Square{0, 0, 0.0, 0, make([]float32, p.NumGridSamples),
		p.NumGridSamples, 0.0, 0, 0, false,
		0.0, 0.0, false, travelList, sync.Mutex{}}

	var wg sync.WaitGroup

	start := time.Now()
	wg.Add(5000)
	for i := 0; i < 5000; i++ {
		go func(i int) {
			defer wg.Done()
			test_squares[i%2].TakeMeasurement(float32(i))
		}(i)
	}

	wg.Wait()
	end := time.Now()
	elapsed := end.Sub(start)
	t.Logf("Time: %v", elapsed)
	if test_squares[0].Avg != 2499.0 {
		t.Errorf("Race condition, got: %f, want: %f",test_squares[0].Avg,2499.0)
	}

	if test_squares[1].Avg != 2500.0 {
		t.Errorf("Race condition, got: %f, want: %f",test_squares[1].Avg,2500.0)
	}
}