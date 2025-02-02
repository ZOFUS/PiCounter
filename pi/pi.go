// pi/pi.go
package pi

import (
	"sync"

	"github.com/ncw/gmp"
)

var (
	C1 = gmp.NewInt(640320)
	C3 = new(gmp.Int).Exp(C1, gmp.NewInt(3), nil)
)

const thresholdPar = 64

func bsPar(a, b int64) (P, Q, T *gmp.Int) {
	if b-a <= thresholdPar {
		return bs(a, b)
	}
	m := (a + b) / 2
	var wg sync.WaitGroup
	wg.Add(2)
	var P1, Q1, T1, P2, Q2, T2 *gmp.Int
	go func() {
		defer wg.Done()
		P1, Q1, T1 = bs(a, m)
	}()
	go func() {
		defer wg.Done()
		P2, Q2, T2 = bs(m, b)
	}()
	wg.Wait()

	P = new(gmp.Int).Mul(P1, P2)
	Q = new(gmp.Int).Mul(Q1, Q2)
	tmp1 := new(gmp.Int).Mul(T1, Q2)
	tmp2 := new(gmp.Int).Mul(P2, T2)
	T = new(gmp.Int).Add(tmp1, tmp2)
	return
}

func bs(a, b int64) (P, Q, T *gmp.Int) {
	if b-a == 1 {
		if a == 0 {
			return gmp.NewInt(1), gmp.NewInt(1), gmp.NewInt(13591409)
		}
		k := gmp.NewInt(a)
		sixK := new(gmp.Int).Mul(gmp.NewInt(6), k)
		t1 := new(gmp.Int).Sub(sixK, gmp.NewInt(5))
		t2 := new(gmp.Int).Sub(new(gmp.Int).Mul(gmp.NewInt(2), k), gmp.NewInt(1))
		t3 := new(gmp.Int).Sub(sixK, gmp.NewInt(1))
		P = new(gmp.Int).Mul(t1, t2)
		P.Mul(P, t3)
		Q = new(gmp.Int).Exp(k, gmp.NewInt(3), nil)
		Q.Mul(Q, C3)
		mult := new(gmp.Int).Mul(gmp.NewInt(545140134), k)
		sumTerm := new(gmp.Int).Add(gmp.NewInt(13591409), mult)
		T = new(gmp.Int).Mul(P, sumTerm)
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
