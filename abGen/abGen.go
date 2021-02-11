package main

import (
	"fmt"
	"math"
)

// var (
// 	// Primes - slice of prime numbers >= 5
// 	Primes []uint64
// 	// A - slice of k-numbers of alpha-primes (6k-1)
// 	A []uint64
// 	// B - slice of k-numbers of beta-primes (6k+1)
// 	B []uint64
// )

type GoldbachTables struct {
	Primes []uint64
	A      []uint64
	B      []uint64
}

func main() {
	fmt.Println("abGen - генерация номеров A- и B- простых чисел")
	const max = 100
	var gt GoldbachTables
	gt.primeGen(max)
	fmt.Printf("Простые числа до %d:\n", max)
	fmt.Println(gt.Primes)
	fmt.Println("\nАльфа-числа, номера:")
	fmt.Println(gt.A)
	fmt.Println("\nБета-числа, номера:")
	fmt.Println(gt.B)
}

// TODO: addition comments!
func (gt *GoldbachTables) primeGen(max uint64) {
	const (
		isAlpha = 4
		isBeta  = 2
	)
	var addition uint64 = isAlpha
	var n uint64
	if len(gt.Primes) > 0 {
		n = gt.Primes[len(gt.Primes)-1] // n is the last known prime or zero
		if (gt.Primes[n-1]+1)%6 == 0 {  // if the last prime in table is an alpha-prime,
			addition = isBeta // select 'beta' addition for the next check
		}
	}

	for n < max {
		n += addition // add to the last known prime 2 or 4
		if gt.isPrime(n) {
			gt.Primes = append(gt.Primes, n)
			if addition == isAlpha {
				gt.A = append(gt.A, (n+1)/6) // new alpha-prime (6k-1), add k to table A
			} else {
				gt.B = append(gt.B, (n-1)/6) // new beta-prime (6k+1), add k to table B
			}
			if addition == isAlpha {
				addition = isBeta
			} else {
				addition = isAlpha
			}
		}
	}
}
func (gt *GoldbachTables) isPrime(n uint64) bool {
	if len(gt.Primes) == 0 {
		if n == 5 {
			gt.Primes = append(gt.Primes, 5)
			return true
		}
		return false
	}
	for i := 0; gt.Primes[i] <= uint64(math.Sqrt(float64(n))); i++ {
		if n%gt.Primes[i] == 0 {
			return false
		}
	}
	return true
}
