// pi/pi_test.go
package pi_test

import (
	"strings"
	"testing"

	"example.com/GO_PRAC/pi"
)

func TestCalculatePi(t *testing.T) {
	digits := 50 // тест с 50 знаками после запятой
	piVal := pi.CalculatePi(digits)
	piStr := piVal.Text('f', digits)

	frac := extractFraction(piStr)
	if len(frac) != digits {
		t.Errorf("Ожидалось %d цифр после запятой, получено %d", digits, len(frac))
	}
}

// extractFraction возвращает дробную часть числа в строке (после точки)
func extractFraction(s string) string {
	i := strings.Index(s, ".")
	if i == -1 {
		return ""
	}
	return s[i+1:]
}
