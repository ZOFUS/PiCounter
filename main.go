// main.go
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"example.com/GO_PRAC/pi"
	"github.com/schollz/progressbar/v3"
)

const (
	precision = 10000
)

func main() {
	fmt.Printf("Вычисление числа π с точностью до %d знаков...\n", precision)

	// Профилирование CPU
	cpuProfFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile file:", err)
		return
	}
	defer cpuProfFile.Close()
	if err := pprof.StartCPUProfile(cpuProfFile); err != nil {
		fmt.Println("Error starting CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	maxN := int64(6 * precision)
	pi.PrecomputeFactorials(maxN)

	bar := progressbar.NewOptions(precision,
		progressbar.OptionSetDescription("Прогресс"),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: "-", BarStart: "[", BarEnd: "]"}),
		progressbar.OptionThrottle(100*time.Millisecond),
	)

	piValue, err := pi.CalculatePi(precision, func(iteration int) {
		fmt.Printf("\nВычислено итераций: %d\n", iteration)
	}, maxN)
	if err != nil {
		fmt.Println("Error calculating Pi:", err) 
		return                                    
	}

	fmt.Println("\nРезультат вычисления числа π:")
	fmt.Println(piValue.Text('f', precision))

	// Memory Profiling
	memProfFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Error creating memory profile file:", err)
		return
	}
	defer memProfFile.Close()
	runtime.GC() 
	if err := pprof.WriteHeapProfile(memProfFile); err != nil {
		fmt.Println("Error writing memory profile:", err)
		return
	}
}
