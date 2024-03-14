package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
	"trucode/finalproject/controllers"
	"trucode/finalproject/env"
)

func main() {
	env.LoadVars()
	// START PROFILING (CPU)
	cpu, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpu.Close()
	pprof.StartCPUProfile(cpu)
	defer pprof.StopCPUProfile()
	// END PROFILING

	log.Println("Starting indexer!")
	startTime := time.Now()

	res := controllers.CheckIfIndexExists()
	if res == nil {
		log.Fatal("response is nil")
	}
	if res.StatusCode == http.StatusOK {
		controllers.DeleteIndex(os.Getenv("INDEX_NAME"))
	}
	emails, err := controllers.GetEmailsDir()
	if err != nil {
		fmt.Println("failed to get emails directory: ", err)
	}

	fmt.Println("Result:")
	fmt.Println("Index:", os.Getenv("INDEX_NAME"))
	fmt.Println("Records:", len(emails))
		
	duration := time.Since(startTime)
	log.Printf("Finished indexing. Time taken: %.2f seconds", duration.Seconds())

	// START PROFILING (MEMORY)
	runtime.GC()
	mem, err := os.Create("memory.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer mem.Close()
	if err := pprof.WriteHeapProfile(mem); err != nil {
		log.Fatal(err)
	}
	// END PROFILING

}
