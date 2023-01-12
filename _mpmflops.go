package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type myCalcs struct {
	x    []float32
	xlen int
}

var xCalcs myCalcs

const heading string = "MP MFLOPS Benchmark"

var (
	threads int = 1
	words   int = 102400 // E Number of words in arrays
	repeats int = 2500   // R Number of repeat passes
)

const (
	xval float32 = 0.999950
	aval float32 = 0.000020
	bval float32 = 0.999980
	cval float32 = 0.000011
	dval float32 = 1.000011
	eval float32 = 0.000012
	fval float32 = 0.999992
	gval float32 = 0.000013
	hval float32 = 1.000013
	jval float32 = 0.000014
	kval float32 = 0.999994
	lval float32 = 0.000015
	mval float32 = 1.000015
	oval float32 = 0.000016
	pval float32 = 0.999996
	qval float32 = 0.000017
	rval float32 = 1.000017
	sval float32 = 0.000018
	tval float32 = 1.000018
	uval float32 = 0.000019
	vval float32 = 1.000019
	wval float32 = 0.000021
	yval float32 = 1.000021

	newdata float32 = 0.999999
	title   string  = "Data in & out "
)

func triadPlus2(n int, a, b, c, d, e, f, g, h, j, k,
	l, m, o, p, q, r, s, t, u, v, w, y float32, x []float32) {

	for i := 0; i < n; i++ {
		x[i] = (x[i]+a)*b - (x[i]+c)*d + (x[i]+e)*f -
			(x[i]+g)*h + (x[i]+j)*k - (x[i]+l)*m +
			(x[i]+o)*p - (x[i]+q)*r + (x[i]+s)*t -
			(x[i]+u)*v + (x[i]+w)*y
	}
}

func triadPlus(n int, a, b, c, d, e, f float32, x []float32) {

	for i := 0; i < n; i++ {
		x[i] = (x[i]+a)*b - (x[i]+c)*d + (x[i]+e)*f
	}
}

func triad(n int, a, b float32, x []float32) {

	for i := 0; i < n; i++ {
		x[i] = (x[i] + a) * b
	}
}

func runTests(offset int, part int) {
	var wds int = xCalcs.xlen
	var xcp []float32 = xCalcs.x[offset*wds:]

	for i := 0; i < repeats; i++ {
		// calculations in CPU
		switch part {
		case 0:
			triad(wds, aval, xval, xcp)
		case 1:
			triadPlus(wds, aval, bval, cval, dval, eval, fval, xcp)
		case 2:
			triadPlus2(wds, aval, bval, cval, dval, eval, fval, gval, hval,
				jval, kval, lval, mval, oval, pval, qval,
				rval, sval, tval, uval, vval, wval, yval, xcp)
		}
	}
}

func getOpwd(part int) int {
	switch part {
	case 0:
		return 2
	case 1:
		return 8
	case 2:
		return 32
	default:
		return 0
	}
}

func initXCpu(value float32, xCpu []float32) {
	for i := 0; i < words; i++ {
		xCpu[i] = value
	}
}

func runTest(initValue float32, part int, threads int, xCpu []float32) time.Duration {
	initXCpu(initValue, xCpu)

	waitGroup := sync.WaitGroup{}

	startTime := time.Now()
	for thread := 0; thread < threads; thread++ {
		waitGroup.Add(1)
		go func(idx int) {
			defer waitGroup.Done()
			runTests(idx, part)
		}(thread)
	}
	waitGroup.Wait()

	return time.Since(startTime)
}

func main() {
	threads := runtime.NumCPU()
	var xCpu []float32
	var size_x int

	var fpmops float64
	var mflops float64

	var isok1 bool = false
	var isok2 bool = false
	var count1 int32 = 0
	var errors [2]string
	var erdata [5][10]int

	var pStart int = 0
	var pEnd int = 3
	var calibrate bool = false
	var runSecs float64 = 0.0

	repeats = repeats * threads

	startWords := words
	startRepeats := repeats

	fmt.Printf("%d CPUs Available\n", threads)
	fmt.Printf("\n")
	fmt.Printf("##############################################\n\n")
	fmt.Printf("  %s, %d Threads, %v\n", heading, threads, time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Printf("  Test             4 Byte  Ops/   Repeat    Seconds   MFLOPS       First   All\n")
	fmt.Printf("                    Words  Word   Passes                         Results  Same\n\n")

	for part := pStart; part < 3; part++ {
		isok1 = false
		words = startWords
		repeats = startRepeats
		for p := 0; p < pEnd; p++ {
			size_x = words * 4

			// Allocate arrays for host CPU
			xCpu = make([]float32, size_x)
			xCalcs.x = xCpu
			xCalcs.xlen = words / threads

			if calibrate {
				endTime := runTest(newdata, part, threads, xCpu)
				repeats = int(float64(repeats) * 15.0 / endTime.Seconds())
				startRepeats = repeats

				calibrate = false
			}

			endTime := runTest(newdata, part, threads, xCpu)

			opwd := getOpwd(part)

			fpmops = float64(words * opwd)
			mflops = float64(repeats) * fpmops / 1000000.0 / endTime.Seconds()
			runSecs = runSecs + endTime.Seconds()

			// Print results
			fmt.Printf("%15s %9d %5d %8d %10.6f %8.0f ", title, words, opwd, repeats, endTime.Seconds(), mflops)
			isok1 = false
			one := xCpu[0]

			if one == newdata {
				isok2, isok1 = true, true
			}

			for i := 1; i < words; i++ {
				if one != xCpu[i] {
					isok1 = true
					if count1 < 10 {
						errors[0] = fmt.Sprintf("%f", xCpu[i])
						errors[1] = fmt.Sprintf("%f", one)
						erdata[0][count1] = i
						erdata[1][count1] = words
						erdata[2][count1] = opwd
						erdata[3][count1] = repeats

						count1 = count1 + 1
					}
				}
			}
			if !isok1 {
				fmt.Printf(" %10.6f   Yes\n", xCpu[0])
			} else {
				fmt.Printf("   See log     No\n")
			}

			words = words * 10
			repeats = repeats / 10

			if repeats < 1 {
				repeats = 1
			}
		}

		fmt.Printf("\n")
	}

	if isok2 {
		fmt.Printf(" ERROR - At least one first result of 0.999999 - no calculations?\n\n")
	}

	fmt.Printf("               End of test %s\n", time.Now().Format("Mon Jan 2 15:04:05 2006"))
}
