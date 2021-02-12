package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"time"
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

var (
	fPrimeStr = "primes.txt"
	fAStr     = "A.txt"
	fBStr     = "B.txt"
)

func main() {
	fmt.Println("abGen - генерация номеров A- и B- простых чисел")
	const max = 100000000
	gt := GoldbachTables{}
	start := time.Now()
	gt.primeGen(max)
	end := time.Now().Sub(start)
	fmt.Printf("Простых чисел от 5 до %d: %d\n", max, len(gt.Primes))
	// fmt.Println(gt.Primes)
	fmt.Printf("Количество альфа-простых: %d\n", len(gt.A))
	fmt.Printf("Количество бета-простых: %d\n", len(gt.B))
	fmt.Printf("\nВыполнено за %v сек. Средняя скорость обработки %v чисел в миллисекунду", end.Seconds(), max/end.Milliseconds())

	writeData(fPrimeStr, gt.Primes)
	writeData(fAStr, gt.A)
	writeData(fBStr, gt.B)

}

func writeData(filename string, data []uint64) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	for i, b := range data {
		_, err := fmt.Fprintf(w, "%d", b)
		check(err)
		if i != len(data)-1 {
			_, err := w.WriteString(" ")
			check(err)
		}
	}
	log.Printf("Writing %v bytes to %s\n", w.Size(), filename)
	err = w.Flush()
	check(err)
	log.Println("Done")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: add comments!
func (gt *GoldbachTables) primeGen(max uint64) {
	const (
		isAlpha = 4
		isBeta  = 2
	)
	var addition uint64 = isAlpha
	var n uint64 = 1
	if len(gt.Primes) > 0 {
		n = gt.Primes[len(gt.Primes)-1] // n is the last known prime or 1
		if (n+1)%6 == 0 {               // if the last prime in table is an alpha-prime,
			addition = isBeta // select 'beta' addition for the next check
		}
	}
	for n+addition < max {
		n += addition // add to the last known prime 2 or 4
		fmt.Printf("n=%d\r", n)
		if gt.isPrime(n) {
			gt.Primes = append(gt.Primes, n)
			if addition == isAlpha {
				gt.A = append(gt.A, (n+1)/6) // new alpha-prime (6k-1), add k to table A
			} else {
				gt.B = append(gt.B, (n-1)/6) // new beta-prime (6k+1), add k to table B
			}
		}
		if addition == isAlpha {
			addition = isBeta
		} else {
			addition = isAlpha
		}
	}
}
func (gt *GoldbachTables) isPrime(n uint64) bool {
	if len(gt.Primes) == 0 {
		if n == 5 {
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
