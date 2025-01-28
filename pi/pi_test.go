// pi/pi_test.go
package pi_test

import (
	"fmt"
	"math"
	"math/big"
	"runtime"
	"strings"
	"testing"

	"example.com/GO_PRAC/pi"
)

func TestCalculatePiLoopCondition(t *testing.T) {
	digits := 3
	prec := uint(digits * 10 / 3)
	maxN := int64(6 * digits)
	pi.PrecomputeFactorials(maxN)

	termThreshold := new(big.Float).Quo(big.NewFloat(1), new(big.Float).SetPrec(prec).SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)))

	k := int64(0)
	batchSize := int64(digits / runtime.NumCPU())
	loopBroke := false
	outputBuffer := new(strings.Builder)

	var calculatedPi *big.Float

	for i := 0; i < 20; i++ {
		end := k + batchSize
		k = end + 1

		if end*6 > maxN {
			end = maxN / 6
		}
		if end <= 0 {
			end = 1
		}

		// Declare missing variables HERE - for unit test loop scope
		numerator := new(big.Int)                // <--- Declare numerator
		denominator := new(big.Int)              // <--- Declare denominator
		_ = denominator                          // <--- Фиктивное использование denominator, чтобы убрать предупреждение
		term := new(big.Float).SetPrec(prec)     // <--- Declare term
		pow := new(big.Int)                      // <--- Declare pow
		tmpFloat := new(big.Float).SetPrec(prec) // <--- Declare tmpFloat

		numInt, err := pi.Factorial(6 * k)
		if err != nil {
			t.Fatalf("Factorial error: %v", err)
		}
		num := new(big.Int).Set(numInt)

		num.Mul(num, big.NewInt(13591409+545140134*k))
		if k%2 != 0 {
			numerator.Neg(numerator)
		}
		denInt, err := pi.Factorial(3 * k)
		if err != nil {
			t.Fatalf("Factorial error: %v", err)
		}
		den := new(big.Int).Set(denInt)

		fkInt, err := pi.Factorial(int64(k))
		if err != nil {
			t.Fatalf("Factorial error: %v", err)
		}
		fk := new(big.Int).Set(fkInt)

		den.Mul(den, fk).Mul(den, fk).Mul(den, fk)
		pow = new(big.Int).SetInt64(640320 * 640320 * 640320)
		pow.Exp(pow, big.NewInt(int64(k)), nil)
		den.Mul(den, pow)

		term.SetInt(numerator)
		term.Quo(term, tmpFloat.SetInt(den))

		absTerm := new(big.Float).SetPrec(prec)
		absTerm.Abs(term)
		cmpResult := absTerm.Cmp(termThreshold)

		fmt.Fprintf(outputBuffer, "Iteration %d: end=%d, termThreshold=%s, term=%s, Comparison=%d (abs(term) < threshold?)\n",
			i, end, termThreshold.Text('e', 2), term.Text('e', 2), cmpResult)

		if k > 0 && cmpResult < 0 {
			fmt.Fprintln(outputBuffer, "Loop condition met, should break at k =", k)
			fmt.Println("Loop condition met, should break at k =", k)
			loopBroke = true
			calculatedPi, err = pi.CalculatePi(digits, nil, maxN)
			if err != nil {
				t.Fatalf("CalculatePi error: %v", err)
			}
			break
		}

	}

	testOutput := outputBuffer.String()

	fmt.Println(testOutput)

	if !loopBroke {
		t.Errorf("Loop did not terminate within %d iterations, condition not met. Output:\n%s", 20, testOutput)
	} else {
		if !strings.Contains(testOutput, "Loop condition met, should break at k =") {
			t.Errorf("Loop terminated, but without 'Loop condition met' message. Output:\n%s", testOutput)
		}

		expectedPi := new(big.Float).SetFloat64(math.Pi).SetPrec(prec)
		tolerance := new(big.Float).SetPrec(prec).SetFloat64(0.001)
		diff := new(big.Float).SetPrec(prec).Sub(calculatedPi, expectedPi)
		absDiff := new(big.Float).SetPrec(prec).Abs(diff)

		if absDiff.Cmp(tolerance) > 0 {
			t.Errorf("Calculated Pi value is not within tolerance. Expected Pi: %s, Calculated Pi: %s, Difference: %s, Tolerance: %s. Output:\n%s",
				expectedPi.Text('f', digits), calculatedPi.Text('f', digits), diff.Text('f', digits), tolerance.Text('f', digits), testOutput)
		} else {
			fmt.Printf("Calculated Pi value is within tolerance. Expected Pi: %s, Calculated Pi: %s, Difference: %s, Tolerance: %s\n",
				expectedPi.Text('f', digits), calculatedPi.Text('f', digits), diff.Text('f', digits), tolerance.Text('f', digits))
		}
	}
}
