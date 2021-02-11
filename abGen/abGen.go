package main

import "fmt"

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
	const max = 1000
	var gt GoldbachTables

}

func (gt *GoldbachTables) primeGen(max uint64) {
	for i := len(gt.Primes); i < max; i++ {

	}
}

func (gt *GoldbachTables) isPrime(n uint64) bool {
	return false
}
