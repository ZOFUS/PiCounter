package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"example.com/GO_PRAC/pi"
)

const digits = 1000000

func main() {
	startTotal := time.Now()
	fmt.Printf("Вычисление числа π с точностью до %d знаков после запятой...\n", digits)

	// Профилирование CPU
	cpuProfFile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Ошибка создания CPU профиля:", err)
		return
	}
	defer cpuProfFile.Close()
	pprof.StartCPUProfile(cpuProfFile)
	defer pprof.StopCPUProfile()

	// Вычисление π
	startCalc := time.Now()
	piStr := pi.CalculatePi(digits) // Результат уже в виде строки
	elapsedCalc := time.Since(startCalc)
	fmt.Printf("Вычисление завершено за %s\n", elapsedCalc)

	// Преобразование в строку не требуется, так как оно выполнено в CalculatePi
	startStr := time.Now()
	elapsedStr := time.Since(startStr)
	fmt.Printf("Преобразование в строку завершено за %s\n", elapsedStr)

	// Запись в файл
	startWrite := time.Now()
	err = os.WriteFile("pi.txt", []byte(piStr), 0644)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
	}
	elapsedWrite := time.Since(startWrite)
	fmt.Printf("Запись в файл завершено за %s\n", elapsedWrite)

	// Профилирование памяти
	startMemProf := time.Now()
	memProfFile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Ошибка создания профиля памяти:", err)
		return
	}
	defer memProfFile.Close()
	runtime.GC()
	if err := pprof.WriteHeapProfile(memProfFile); err != nil {
		fmt.Println("Ошибка записи профиля памяти:", err)
	}
	elapsedMemProf := time.Since(startMemProf)
	fmt.Printf("Профилирование памяти завершено за %s\n", elapsedMemProf)

	totalElapsed := time.Since(startTotal)
	fmt.Printf("Общее время работы: %s\n", totalElapsed)
}
