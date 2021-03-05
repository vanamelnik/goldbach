package main

import (
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	fPrimeBin = "tables/primes.gz"
	fABin     = "tables/A.gz"
	fBBin     = "tables/B.gz"
)

var (
	primes, A, B []uint64
)

func main() {
	fmt.Printf("Reading the table of prime numbers from file %s...", fPrimeBin)
	primes, err := unzipFile(fPrimeBin)
	if err != nil {
		log.Fatalf("\n%v\n", err)
	}
	fmt.Println("OK")
	fmt.Printf("Reading the table of prime numbers from file %s...", fABin)
	A, err = unzipFile(fABin)
	if err != nil {
		log.Fatalf("\n%v\n", err)
	}
	fmt.Println("OK")
	fmt.Printf("Reading the table of prime numbers from file %s...", fBBin)
	B, err = unzipFile(fBBin)
	if err != nil {
		log.Fatalf("\n%v\n", err)
	}
	fmt.Println("OK")

	fmt.Printf("\nThere are %d prime numbers in the table\n", len(primes))
	// fmt.Printf("Last five primes is: %v\n", primes[len(primes)-5:])
	fmt.Printf("\nThere are %d alpha prime indexes in the table\n", len(A))
	// fmt.Printf("Last five indexes is: %v\n", A[len(A)-5:])
	fmt.Printf("\nThere are %d beta prime indexes in the table\n", len(B))
	// fmt.Printf("Last five indexes is: %v\n", B[len(B)-5:])

	const maxN = 100
	for n := 5; n <= maxN; n++ {
		p := (n - 1) / 2
		fmt.Printf("N=%d", n)
		for i := 1; i <= p; i++ {
			// fmt.Printf("\t%d\t%d\t", i, n-i)
			if !(binarySearch(A, uint64(i)) && binarySearch(A, uint64(n-i))) {
				p--
				// fmt.Println()
			} else {
				// fmt.Println("OK")
			}
		}
		fmt.Printf("OK = %d\n", p)
	}
}

func unzipFile(name string) ([]uint64, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	buf := make([]byte, 1048576*100)
	var data []uint64
	var count int
	for {
		n, err := zr.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		count += n
		for i := 0; i < n; i += 8 {
			num := binary.BigEndian.Uint64(buf[i:])
			data = append(data, num)
		}
		if err == io.EOF {
			break
		}
	}

	return data, nil
}

func binarySearch(data []uint64, x uint64) bool {
	max, min := len(data), 0
	i := max / 2
	for x != data[i] {
		// fmt.Printf("min:%d, max:%d, i:%d\n", min, max, i)
		if min == max-1 {
			return false
		} else if x < data[i] {
			max = i
			i = min + (i-min)/2
		} else {
			min = i
			i += (max - i) / 2
		}
	}
	return true
}
