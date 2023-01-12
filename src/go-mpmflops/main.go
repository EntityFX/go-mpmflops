package main

import (
	"runtime"

	"github.com/EntityFX/go-mpmflops/mpmflops"
)

func main() {
	threads := runtime.NumCPU()

	mpmflops.MpmflopsRun(threads)
}
