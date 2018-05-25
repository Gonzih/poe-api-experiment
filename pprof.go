package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// Usage:
// run this before your code runs
// stopFn := enableCPUProf()
// defer stopFn()
func enableCPUProf() func() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	return pprof.StopCPUProfile
}

// Usage:
// run this after your code runs
// enableMEMProf()
func enableMEMProf() {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}
