// pi/pi.go
package pi

import (
	"math/big"
	"sync"

	"github.com/ncw/gmp"
)

// В данном коде мы реализуем алгоритм Чудновского с использованием бинарного разбиения.
// Все операции выполняются через обёртку GMP для быстроты работы с большими числами.

// Глобальные константы (типы *gmp.Int)
var (
	C  = gmp.NewInt(640320)
	C3 = new(gmp.Int).Exp(C, gmp.NewInt(3), nil) // C^3 = 640320^3
)

// thresholdPar определяет минимальный размер диапазона для параллельного разбиения
const thresholdPar = 64

// bsPar выполняет параллельное бинарное разбиение для диапазона [a, b).
// Если диапазон достаточно большой, он разбивается на две части, вычисляемые в горутинах.
func bsPar(a, b int64) (P, Q, T *gmp.Int) {
	if b-a <= thresholdPar {
		return bs(a, b)
	}
	m := (a + b) / 2
	var (
		P1, Q1, T1 *gmp.Int
		P2, Q2, T2 *gmp.Int
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		P1, Q1, T1 = bs(a, m)
	}()
	go func() {
		defer wg.Done()
		P2, Q2, T2 = bs(m, b)
	}()
	wg.Wait()

	// P = P1 * P2, Q = Q1 * Q2, T = T1*Q2 + P2*T2
	P = new(gmp.Int).Mul(P1, P2)
	Q = new(gmp.Int).Mul(Q1, Q2)
	tmp1 := new(gmp.Int).Mul(T1, Q2)
	tmp2 := new(gmp.Int).Mul(P2, T2)
	T = new(gmp.Int).Add(tmp1, tmp2)
	return
}

// bs выполняет последовательное бинарное разбиение для диапазона [a, b).
// Базовый случай: диапазон длины 1 (b-a == 1)
func bs(a, b int64) (P, Q, T *gmp.Int) {
	if b-a == 1 {
		if a == 0 {
			return gmp.NewInt(1), gmp.NewInt(1), gmp.NewInt(13591409)
		}
		// k = a
		k := gmp.NewInt(a)
		// P = (6k-5)*(2k-1)*(6k-1)
		sixK := new(gmp.Int).Mul(gmp.NewInt(6), k)
		t1 := new(gmp.Int).Sub(sixK, gmp.NewInt(5))
		t2 := new(gmp.Int).Sub(new(gmp.Int).Mul(gmp.NewInt(2), k), gmp.NewInt(1))
		t3 := new(gmp.Int).Sub(sixK, gmp.NewInt(1))
		P = new(gmp.Int).Mul(t1, t2)
		P.Mul(P, t3)

		// Q = k^3 * C^3
		Q = new(gmp.Int).Exp(k, gmp.NewInt(3), nil)
		Q.Mul(Q, C3)

		// T = P * (13591409 + 545140134*k)
		mult := new(gmp.Int).Mul(gmp.NewInt(545140134), k)
		sumTerm := new(gmp.Int).Add(gmp.NewInt(13591409), mult)
		T = new(gmp.Int).Mul(P, sumTerm)
		// Если k нечетное, T = -T
		if a%2 != 0 {
			T.Neg(T)
		}
		return
	} else {
		m := (a + b) / 2
		P1, Q1, T1 := bs(a, m)
		P2, Q2, T2 := bs(m, b)
		P = new(gmp.Int).Mul(P1, P2)
		Q = new(gmp.Int).Mul(Q1, Q2)
		tmp1 := new(gmp.Int).Mul(T1, Q2)
		tmp2 := new(gmp.Int).Mul(P2, T2)
		T = new(gmp.Int).Add(tmp1, tmp2)
		return
	}
}

// Float is обёртка для *gmp.Float с дополнительным методом для форматированного вывода.
type Float struct {
	*big.Float
}

// StringWithPrecision возвращает строку с числом π с указанным числом цифр после запятой.
// Для этого используется строковое представление в десятичной системе.
func (f *Float) StringWithPrecision(digits int) string {
	// Получаем строку с экспоненциальным представлением, затем преобразуем к обычной десятичной записи.
	// Метод f.String() по умолчанию выдаёт экспоненциальное представление, поэтому для наших целей
	// можно использовать f.Text('f', digits), если библиотека реализует подобный метод.
	// Если такой функции в обёртке нет, можно преобразовать через f.RatString() и затем отформатировать.
	return f.Float.Text('f', digits)
}

func gmpToBigInt(g *gmp.Int) *big.Int {
	bigInt := new(big.Int)
	bigInt.SetBytes(g.Bytes())
	return bigInt
}

// CalculatePi вычисляет число π с указанной точностью (число знаков после запятой).
// Приблизительно каждый член ряда даёт ~14 цифр, поэтому N = digits/14 + 1.
func CalculatePi(digits int) *big.Float {
	// Определяем рабочую точность (примерно digits*4 бит)
	prec := uint(digits * 4)
	N := int64(digits/14 + 1)

	_, Q, T := bsPar(0, N)

	// Используем эту функцию для Q и T:
	Qbig := gmpToBigInt(Q)
	Tbig := gmpToBigInt(T)

	// Теперь можно передавать их в big.Float
	Qf := new(big.Float).SetPrec(prec).SetInt(Qbig)
	Tf := new(big.Float).SetPrec(prec).SetInt(Tbig)

	// Вычисляем C^(3/2) = 640320^(3/2)
	C := big.NewInt(640320)
	cFloat := new(big.Float).SetPrec(prec).SetInt(C)
	sqrtC := new(big.Float).Sqrt(cFloat)
	c3_2 := new(big.Float).Mul(cFloat, sqrtC)

	// π = (C^(3/2) * Q) / (12 * T)
	num := new(big.Float).Mul(c3_2, Qf)
	denom := new(big.Float).Mul(big.NewFloat(12), Tf)
	piValue := new(big.Float).Quo(num, denom)

	return piValue
}
