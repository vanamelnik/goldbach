package main

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const max = 1000

type dataEncoderFunc func(io.Writer, []uint64) (int, error)

var (
	Primes []uint64
	A      []uint64
	B      []uint64
)

var (
	fPrimeStr = "tables/primes.txt"
	fAStr     = "tables/A.txt"
	fBStr     = "tables/B.txt"

	fPrimeBin = "tables/primes.bin"
	fABin     = "tables/A.bin"
	fBBin     = "tables/B.bin"
)

var start time.Time

func main() {
	fmt.Println("abGen - генерация номеров A- и B- простых чисел")
	capacity := max / int(math.Log(float64(max))) // pi(x) ~ x/log(x)
	fmt.Printf("Поиск простых чисел в интервале 5 - %d. Pi(x) ~ %d\n", max, capacity)
	Primes = make([]uint64, 0, capacity)
	A = make([]uint64, 0, capacity/2)
	B = make([]uint64, 0, capacity/2)

	start = time.Now()
	primeGen(max)
	end := time.Now().Sub(start)
	fmt.Printf("Простых чисел от 5 до %d: %d\n", max, len(Primes))
	fmt.Printf("Количество альфа-простых: %d\n", len(A))
	fmt.Printf("Количество бета-простых: %d\n", len(B))
	if end.Microseconds() != 0 {
		fmt.Printf("\nВыполнено за %v Средняя скорость обработки %v чисел в миллисекунду\n", end.Round(time.Second), max/end.Milliseconds())
	}

	// writeZipData(fPrimeStr, Primes, encodeText)
	// writeZipData(fAStr, A, encodeText)
	// writeZipData(fBStr, B, encodeText)
	writeZipData(fPrimeBin, Primes, encodeBin)
	writeZipData(fABin, A, encodeBin)
	writeZipData(fBBin, B, encodeBin)

}

func writeZipData(filename string, data []uint64, encode dataEncoderFunc) {
	zipName := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".gz"
	f, err := os.Create(zipName)
	check(err)
	defer f.Close()
	zw := gzip.NewWriter(f)
	zw.Name = filepath.Base(filename)
	n, err := encode(zw, data)
	check(err)
	zw.Close()
	err = zw.Flush()
	check(err)
	log.Printf("%d bytes written to %s", n, filename)
}

func writeData(filename string, data []uint64, encode dataEncoderFunc) {
	f, err := os.Create(filename)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	n, err := encode(w, data)
	err = w.Flush()
	check(err)
	log.Printf("%d bytes written to %s", n, filename)
}

func encodeText(w io.Writer, data []uint64) (count int, err error) {
	for i, b := range data {
		n, err := fmt.Fprintf(w, "%d", b)
		count += n
		if err != nil {
			return count, err
		}
		if i != len(data)-1 {
			n, err := fmt.Fprint(w, " ")
			count += n
			if err != nil {
				return count, err
			}
		}
	}
	return count, err
}

func encodeBin(w io.Writer, data []uint64) (count int, err error) {
	b := make([]byte, 8)
	for _, num := range data {
		binary.BigEndian.PutUint64(b, num)
		n, err := w.Write(b)
		count += n
		if err != nil {
			return count, err
		}
	}
	return count, err
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: add comments!
func primeGen(max uint64) {
	const (
		isAlpha = 4
		isBeta  = 2
	)
	var addition uint64 = isAlpha
	var n uint64 = 1
	if len(Primes) > 0 {
		n = Primes[len(Primes)-1] // n is the last known prime or 1
		if (n+1)%6 == 0 {         // if the last prime in table is an alpha-prime,
			addition = isBeta // select 'beta' addition for the next check
		}
	}

	interval := time.Second
	tick := time.Tick(interval)

	// This goroutine calculates time left and prints progress of calculating and estimated time
	q := make(chan int)
	go func(quit chan int) {
		var nOld uint64
		for {
			select {
			case <-quit:
				return
			case <-tick:
				if nOld != n {
					elapsed := time.Now().Sub(start)
					timeLeft := time.Duration(int64((max-n)/(n-nOld)) * time.Second.Nanoseconds())
					fmt.Printf("%d%% complete, elapsed time %v sec., time left %v, speed %d nums/s    \r",
						uint64(float64(n)/float64(max)*100), elapsed.Truncate(time.Second), timeLeft.Truncate(time.Second), n-nOld)
					nOld = n
				}
			}
		}
	}(q)

	for n+addition < max {
		n += addition // add to the last known prime 2 or 4
		if isPrime(n) {
			Primes = append(Primes, n)
			if addition == isAlpha {
				A = append(A, (n+1)/6) // new alpha-prime (6k-1), add k to table A
			} else {
				B = append(B, (n-1)/6) // new beta-prime (6k+1), add k to table B
			}
		}
		if addition == isAlpha {
			addition = isBeta
		} else {
			addition = isAlpha
		}
	}

	q <- 0
	time.Sleep(time.Millisecond)
	fmt.Printf("%s\r", string(80*' '))
}

func isPrime(n uint64) bool {
	if len(Primes) == 0 {
		if n == 5 {
			return true
		}
		return false
	}
	for i := 0; Primes[i] <= uint64(math.Sqrt(float64(n))); i++ {
		if n%Primes[i] == 0 {
			return false
		}
	}
	return true
}
