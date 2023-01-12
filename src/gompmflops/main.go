package main

import (
	"gompmflops/mpmflops"
	"runtime"
)

func main() {
	threads := runtime.NumCPU()
	println(threads)
	mpmflops.MpmflopsRun(threads)
}
