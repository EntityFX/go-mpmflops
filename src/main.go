package main

import (
	"runtime"

	"github.com/entityfx/mpmflops/mpmflops"
)

func main() {
	threads := runtime.NumCPU()

	mpmflops.MpmflopsRun(threads)
}
