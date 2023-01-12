package mpmflops

import (
	"fmt"
	"mpmflops"
	"sync"
	"time"
)

const (
	title      string = "Data in & out "
	titleConst string = "Data const "
	heading    string = "MP MFLOPS Benchmark"

	numberOfWords   int = 102400 // E Number of words in arrays
	numberOfRepeats int = 2500   // R Number of repeat passes
)

type RunTestsFunc func(words int, repeats int, offset int, part int, threads int, xCpu []float32)

func RunTests(words int, repeats int, offset int, part int, threads int, xCpu []float32) {
	wds := words / threads
	xcp := xCpu[offset*wds:]

	for i := 0; i < repeats; i++ {
		// calculations in CPU
		switch part {
		case 0:
			Triad(wds, aval, xval, xcp)
		case 1:
			TriadPlusMid(wds, aval, bval, cval, dval, eval, fval, xcp)
		case 2:
			TriadPlusLarge(wds, aval, bval, cval, dval, eval, fval, gval, hval,
				jval, kval, lval, mval, oval, pval, qval,
				rval, sval, tval, uval, vval, wval, yval, xcp)
		}
	}
}

func RunConstTests(words int, repeats int, offset int, part int, threads int, xCpu []float32) {
	wds := words / threads
	xcp := xCpu[offset*wds:]

	for i := 0; i < repeats; i++ {
		// calculations in CPU
		switch part {
		case 0:
			triadConst(wds, xcp)
		case 1:
			triadConstPlusMid(wds, xcp)
		case 2:
			triadConstPlusLarge(wds, xcp)
		}
	}
}

func RunParallelTests(runTestsFunc RunTestsFunc, words int, repeats int, initValue float32, part int, threads int, xCpu []float32) time.Duration {
	initXCpu(initValue, xCpu)

	waitGroup := sync.WaitGroup{}

	startTime := time.Now()
	for thread := 0; thread < threads; thread++ {
		waitGroup.Add(1)
		go func(idx int) {
			defer waitGroup.Done()
			runTestsFunc(words, repeats, idx, part, threads, xCpu)
		}(thread)
	}
	waitGroup.Wait()

	return time.Since(startTime)
}

func RunAllTests(runTestsFunc mpmflops.RunTestsFunc, threads int, numberOfRepeats int, startWords int, parts int, calibrate bool, title string) {
	sizeX := 0
	mflops := 0.0
	isTestsOk := true
	runSecs := time.Duration(0)
	startRepeats := numberOfRepeats * threads

	for part := 0; part < parts; part++ {
		words := startWords
		repeats := startRepeats
		for p := 0; p < parts; p++ {
			sizeX = words * 4

			// Allocate arrays for host CPU
			xCpu := make([]float32, sizeX)

			if calibrate {
				endTime := mpmflops.RunParallelTests(runTestsFunc, words, repeats, mpmflops.Newdata, part, threads, xCpu)
				repeats = int(float64(repeats) * 15.0 / endTime.Seconds())
				startRepeats = repeats

				calibrate = false
			}

			testDuration := mpmflops.RunParallelTests(runTestsFunc, words, repeats, mpmflops.Newdata, part, threads, xCpu)

			opwd := mpmflops.getOpwd(part)

			fpmops := float64(words * opwd)
			mflops = float64(repeats) * fpmops / 1000000.0 / testDuration.Seconds()
			runSecs += testDuration

			// Print results
			fmt.Printf("%15s %9d %5d %8d %10.6f %8.0f ", title, words, opwd, repeats, testDuration.Seconds(), mflops)

			isTestOk := true
			isTestOk, isTestsOk = Validate(xCpu, words)

			words = words * 10
			repeats = repeats / 10

			if repeats < 1 {
				repeats = 1
			}

			if isTestOk {
				fmt.Printf(" %10.6f   Yes\n", xCpu[0])
			} else {
				fmt.Printf("   See log     No\n")
			}
		}

		fmt.Printf("\n")
	}

	if !isTestsOk {
		fmt.Printf(" ERROR - At least one first result of 0.999999 - no calculations?\n\n")
	}
}

func MpmflopsRun(threads int) {
	fmt.Printf("%d CPUs Available\n", threads)
	fmt.Printf("\n")
	fmt.Printf("##############################################\n\n")
	fmt.Printf("  %s, %d Threads, %v\n", heading, threads, time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Printf("  Test             4 Byte  Ops/   Repeat    Seconds   MFLOPS       First   All\n")
	fmt.Printf("                    Words  Word   Passes                         Results  Same\n\n")

	runTestFuncs := map[string]mpmflops.RunTestsFunc{
		title: func(words int, repeats int, offset int, part int, threads int, xCpu []float32) {
			mpmflops.RunTests(words, repeats, offset, part, threads, xCpu)
		},
		titleConst: func(words int, repeats int, offset int, part int, threads int, xCpu []float32) {
			mpmflops.RunConstTests(words, repeats, offset, part, threads, xCpu)
		},
	}

	for k, v := range runTestFuncs {
		mpmflops.RunAllTests(v, threads, numberOfRepeats, numberOfWords, 3, false, k)
		fmt.Printf("##############################################\n\n")
	}

	fmt.Printf("               End of test %s\n", time.Now().Format("Mon Jan 2 15:04:05 2006"))
}