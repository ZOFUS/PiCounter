// pi/pi.go
package pi

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
)

const (
	UpdateEvery = 10
	cacheSize   = 500 // Размер кольцевого буфера для кэширования факториалов
)

// cachedFactorial кэширует вычисленные факториалы в кольцевом буфере
var cachedFactorial []factorialCacheEntry
var factorialCacheIndex int64
var factorialMutex sync.Mutex

type factorialCacheEntry struct {
	n      int64
	value  *big.Int
	cached bool
}

func PrecomputeFactorials(maxN int64) {
	factorialMutex.Lock()
	defer factorialMutex.Unlock()

	cachedFactorial = make([]factorialCacheEntry, cacheSize) // Инициализируем кольцевой буфер

	cachedFactorial[0] = factorialCacheEntry{n: 0, value: big.NewInt(1), cached: true} // Кэшируем 0!
	cachedFactorial[1] = factorialCacheEntry{n: 1, value: big.NewInt(1), cached: true} // Кэшируем 1!
	factorialCacheIndex = 2                                                            // Начинаем с индекса 2 для добавления следующих факториалов

	for i := int64(2); i <= maxN; i++ {
		index := factorialCacheIndex % cacheSize // Рассчитываем индекс в кольцевом буфере
		previousIndex := (factorialCacheIndex - 1 + cacheSize) % cacheSize

		cachedFactorial[index] = factorialCacheEntry{
			n:      i,
			value:  new(big.Int).Mul(cachedFactorial[previousIndex].value, big.NewInt(i)),
			cached: true,
		}

		factorialCacheIndex++
	}
}

// Factorial вычисляет факториал числа n с использованием кэша (кольцевой буфер)
func Factorial(n int64) (*big.Int, error) {
	factorialMutex.Lock()
	defer factorialMutex.Unlock()

	// Проверяем, есть ли факториал в кэше
	for _, entry := range cachedFactorial {
		if entry.n == n && entry.cached {
			return new(big.Int).Set(entry.value), nil
		}
	}

	// Если факториал не найден, вычисляем его
	result := big.NewInt(1)
	for i := int64(1); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}

	// Добавляем результат в кэш, если размер позволяет
	if int64(len(cachedFactorial)) < cacheSize {
		cachedFactorial = append(cachedFactorial, factorialCacheEntry{
			n:      n,
			value:  new(big.Int).Set(result),
			cached: true,
		})
	} else {
		// Перезаписываем самый старый элемент в кольцевом буфере
		index := factorialCacheIndex % cacheSize
		cachedFactorial[index] = factorialCacheEntry{
			n:      n,
			value:  new(big.Int).Set(result),
			cached: true,
		}
		factorialCacheIndex++
	}

	return result, nil
}

// calculatePiPartial вычисляет часть ряда для заданного диапазона k
func calculatePiPartial(start, end int64, precision uint, sum *big.Float) error {
	const (
		multiplier  = 640320
		multiplier3 = multiplier * multiplier * multiplier
	)

	// Переиспользуемые объекты big.Int и big.Float, выделены вне цикла
	numerator := new(big.Int)
	denominator := new(big.Int)
	term := new(big.Float).SetPrec(precision)
	pow := new(big.Int)
	tmpFloat := new(big.Float).SetPrec(precision)
	const545140134 := big.NewInt(545140134)     // Предварительное выделение констант
	const13591409 := big.NewInt(13591409)       // Предварительное выделение констант
	constMultiplier3 := big.NewInt(multiplier3) // Предварительное выделение констант
	constantTerm := new(big.Int)                //  <---  Предварительное выделение для (13591409 + 545140134*k)

	for k := start; k <= end; k++ {
		numeratorInt, err := Factorial(6 * k)
		if err != nil {
			return fmt.Errorf("Factorial(6*%d) error: %w", k, err)
		}
		numerator.Set(numeratorInt)

		//Переиспользование объектов
		constantTerm.Mul(big.NewInt(k), const545140134) // Вычисляем (545140134*k) как big.Int
		constantTerm.Add(constantTerm, const13591409)   // Прибавляем 13591409 к результату (как big.Int)
		numerator.Mul(numerator, constantTerm)          // Умножаем numerator на constantTerm (оба big.Int)
		if k%2 != 0 {
			numerator.Neg(numerator)
		}

		denominatorInt, err := Factorial(3 * k)
		if err != nil {
			return fmt.Errorf("Factorial(3*%d) error: %w", k, err)
		}
		denominator.Set(denominatorInt)

		fkInt, err := Factorial(k)
		if err != nil {
			return fmt.Errorf("Factorial(%d) error: %w", k, err)
		}
		fk := new(big.Int).Set(fkInt)

		denominator.Mul(denominator, fk).Mul(denominator, fk).Mul(denominator, fk)
		pow.Set(big.NewInt(multiplier3))
		pow.Set(constMultiplier3)
		pow.Exp(pow, big.NewInt(k), nil)
		denominator.Mul(denominator, pow)

		term.SetInt(numerator)
		term.Quo(term, tmpFloat.SetInt(denominator))

		sum.Add(sum, term)
	}
	return nil
}

// CalculatePi вычисляет число π с использованием пула горутин и возвращает ошибку
func CalculatePi(digits int, update func(iteration int), maxN int64) (*big.Float, error) {
	const multiplier = 640320
	numCPU := runtime.NumCPU()
	prec := uint(digits * 10 / 3)
	sum := new(big.Float).SetPrec(prec)
	var wg sync.WaitGroup
	var mu sync.Mutex
	k := int64(0)
	batchSize := int64(digits / (numCPU / 2))

	termThreshold := new(big.Float).Quo(big.NewFloat(1), new(big.Float).SetPrec(prec).SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil))) // Вычисляем termThreshold здесь
	fmt.Println("Term threshold:", termThreshold.Text('e', 2))                                                                                                   // Выводим termThreshold для отладки
	const545140134 := big.NewInt(545140134)                                                                                                                      // Предварительное выделение константы для горутин
	const13591409 := big.NewInt(13591409)                                                                                                                        // Предварительное выделение константы для горутин

	taskChan := make(chan struct {
		start int64
		end   int64
		k     int64
	}, numCPU)

	// Запускаем worker-горутины
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			term := new(big.Float).SetPrec(prec)
			tmpFloat := new(big.Float).SetPrec(prec)
			constantTerm := new(big.Int)

			for task := range taskChan {
				partialSum := new(big.Float).SetPrec(prec)
				err := calculatePiPartial(task.start, task.end, prec, partialSum)
				if err != nil {
					fmt.Printf("Error in calculatePiPartial (worker %d, k=%d): %v\n", workerID, task.k, err)
					return
				}

				mu.Lock()
				sum.Add(sum, partialSum)
				mu.Unlock()

				//  Вычисляем term *только для *текущей* итерации k
				numInt, err := Factorial(6 * task.k)
				if err != nil {
					fmt.Printf("Error in Factorial (worker %d, k=%d): %v\n", workerID, task.k, err)
					return
				}
				num := new(big.Int).Set(numInt)

				// Переиспользование объектов и big.Int вычисления
				constantTerm.Mul(big.NewInt(task.k), const545140134)
				constantTerm.Add(constantTerm, const13591409)
				num.Mul(num, constantTerm)
				if task.k%2 != 0 {
					num.Neg(num)
				}
				denInt, err := Factorial(3 * task.k)
				if err != nil {
					fmt.Printf("Error in Factorial (worker %d, k=%d): %v\n", workerID, task.k, err)
					return
				}
				den := new(big.Int).Set(denInt)

				fkInt, err := Factorial(int64(task.k))
				if err != nil {
					fmt.Printf("Error in Factorial (worker %d, k=%d): %v\n", workerID, task.k, err)
					return
				}
				fk := new(big.Int).Set(fkInt)

				den.Mul(den, fk).Mul(den, fk).Mul(den, fk)
				pow := new(big.Int).SetInt64(multiplier * multiplier * multiplier)
				pow.Exp(pow, big.NewInt(int64(task.k)), nil)
				den.Mul(den, pow)

				term.SetInt(num)
				tmpFloat.SetInt(den)
				term.Quo(term, tmpFloat)

				if task.end%UpdateEvery == 0 {
					if update != nil {
						update(int(task.end))
					}
				}
			}
			fmt.Printf("Worker %d завершил работу\n", workerID)
		}(i)
	}

	for {
		start := k
		end := k + batchSize
		k = end + 1

		if end*6 > maxN {
			end = maxN / 6
		}

		taskChan <- struct {
			start int64
			end   int64
			k     int64
		}{start: start, end: end, k: k - 1}

		term := new(big.Float).SetPrec(prec)

		// Вычисляем term для проверки условия выхода (вне горутины, в основном потоке)
		numInt, err := Factorial(6 * (k - 1))
		if err != nil {
			return nil, fmt.Errorf("Factorial(6*%d) error: %w", k-1, err)
		}
		num := new(big.Int).Set(numInt)

		num.Mul(num, big.NewInt(13591409+545140134*(k-1)))
		if (k-1)%2 != 0 {
			num.Neg(num)
		}
		denInt, err := Factorial(3 * (k - 1))
		if err != nil {
			return nil, fmt.Errorf("Factorial(3*%d) error: %w", k-1, err)
		}
		den := new(big.Int).Set(denInt)

		fkInt, err := Factorial(int64(k - 1))
		if err != nil {
			return nil, fmt.Errorf("Factorial(%d) error: %w", k-1, err)
		}
		fk := new(big.Int).Set(fkInt)

		den.Mul(den, fk).Mul(den, fk).Mul(den, fk)
		pow := new(big.Int).SetInt64(multiplier * multiplier * multiplier)
		pow.Exp(pow, big.NewInt(int64(k-1)), nil)
		den.Mul(den, pow)

		term.SetInt(num)
		tmpFloat := new(big.Float).SetPrec(prec).SetInt(den)
		term.Quo(term, tmpFloat)

		if k > 0 {
			absTerm := new(big.Float).SetPrec(prec)
			absTerm.Abs(term)
			if absTerm.Cmp(termThreshold) < 0 {
				fmt.Println("Loop condition met, should break at k =", k-1)
				break
			}
		}
		if end*6 >= maxN && k > int64(digits) {
			fmt.Println("Max iterations reached, should break at k =", k-1)
			break
		}
	}
	close(taskChan)
	wg.Wait()

	multiplierSqrt3 := new(big.Float).SetPrec(prec).SetFloat64(multiplier)
	multiplierSqrt3.Sqrt(multiplierSqrt3)
	multiplierSqrt3.Mul(multiplierSqrt3, new(big.Float).SetPrec(prec).SetFloat64(multiplier))

	pi := new(big.Float).SetPrec(prec)
	pi.SetInt64(12)
	pi.Mul(pi, sum)
	pi.Quo(multiplierSqrt3, pi)

	return pi, nil
}
