package main

import (
	"flag"
	"runtime"

	"github.com/EntityFX/go-mpmflops/mpmflops"
)

func main() {
	threadsPtr := flag.Int("t", runtime.NumCPU(), "Threads count")
	calibratePtr := flag.Bool("c", false, "Calibrate")

	flag.Parse()

	if *threadsPtr < 1 || *threadsPtr > runtime.NumCPU() {
		*threadsPtr = runtime.NumCPU()
	}

	mpmflops.MpmflopsRun(*threadsPtr, *calibratePtr)
}
