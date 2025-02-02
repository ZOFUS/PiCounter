// pi/pi_C.go
package pi

/*
#cgo LDFLAGS: -lgmp
#include <gmp.h>
#include <stdlib.h>

char* compute_pi(int digits) {
    // Реализация вычисления π с использованием GMP
    // ... (ваш код на C с mpf_t) ...
    char *str;
    gmp_asprintf(&str, "%.*Ff", digits, pi); // Преобразование в строку
    return str;
}
*/
import "C"

import (
	"unsafe"
)

func CalculatePi(digits int) string {
	cStr := C.compute_pi(C.int(digits))
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr)
}
