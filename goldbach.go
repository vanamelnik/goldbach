package main

import (
	"fmt"
	"goldbach/files"
	"log"
	"sort"
	"time"
)

const (
	fPrimeBin = "tables/primes.gz"
	fABin     = "tables/A.gz"
	fBBin     = "tables/B.gz"

	maxR   = 1000
	firstR = 1
)

var (
	primes, A, B []uint64
)

func main() {
	// fmt.Printf("Reading the table of prime numbers from file %s...", fPrimeBin)
	// primes, err := unzipFile(fPrimeBin)
	// if err != nil {
	// 	log.Fatalf("\r\n%v\r\n", err)
	// }
	// fmt.Println("OK")
	fmt.Printf("Reading the table of alpha indexes from file %s...", fABin)
	var err error
	A, err = files.UnzipFile(fABin)
	if err != nil {
		log.Fatalf("\r\n%v\r\n", err)
	}
	fmt.Println("OK")
	fmt.Printf("Reading the table of beta indexes from file %s...", fBBin)
	B, err = files.UnzipFile(fBBin)
	if err != nil {
		log.Fatalf("\r\n%v\r\n", err)
	}
	fmt.Println("OK")

	// Bmax := B[len(B)-1]
	// _ = Bmax // TODO: Remove this stub

	// fmt.Printf("\r\nThere are %d prime numbers in the table\r\n", len(primes))
	// fmt.Printf("Last five primes is: %v\r\n", primes[len(primes)-5:])
	fmt.Printf("\r\nThere are %d alpha prime indexes in the table\r\n", len(A))
	fmt.Printf("Last five indexes is: %v\r\n", A[len(A)-5:])
	fmt.Printf("\r\nThere are %d beta prime indexes in the table\r\n", len(B))
	fmt.Printf("Last five indexes is: %v\r\n", B[len(B)-5:])
	fmt.Println("Calculating v(r) B + A")
	calculateN2(B, A)
}

func binarySearch(data []uint64, x uint64) bool {
	// max, min := len(data), 0
	// i := max / 2
	// for x != data[i] {
	// 	// fmt.Printf("min:%d, max:%d, i:%d\r\n", min, max, i)
	// 	if min == max-1 {
	// 		return false
	// 	} else if x < data[i] {
	// 		max = i
	// 		i = min + (i-min)/2
	// 	} else {
	// 		min = i
	// 		i += (max - i) / 2
	// 	}
	// }
	// return true
	k := sort.Search(len(data), func(i int) bool { return data[i] >= x })
	if k < len(data) && data[k] == x {
		return true
	} else {
		return false
	}
}

func calculateN(arr []uint64) {
	max := arr[len(arr)-1]
	N := uint64(firstR*2 + 1)
	for r := firstR; r <= maxR; r++ {
		// fmt.Printf("r = %d", r)
		i, maxI := 0, 0
		q := make(chan bool)
		tick := time.Tick(time.Second / 4)
		go func() {
			for {
				select {
				case <-tick:
					fmt.Printf("N = %d, iMax = %d, %d%% complete                \r", N, maxI, uint64(maxI)*100/uint64(r))
				case <-q:
					return
				}
			}
		}()
		for ; i < r; i++ {
			// fmt.Printf("\tA(i)=%d, N=%d, N-A(i) = %d-%d = %d ", A[i], N, N, A[i], uint64(N)-A[i])
			if i > maxI {
				maxI = i
			}
			if binarySearch(arr, N-arr[i]) {
				N++
				// fmt.Printf("N=%d\r", N)
				if N > max {
					q <- true
					time.Sleep(time.Millisecond)
					fmt.Println("\nReached max number in the table!")
					return
				}
				i = -1
				// fmt.Print("- white, N = N+1\r\n")
			}
		}
		q <- true
		time.Sleep(time.Millisecond)
		fmt.Printf("N(%d) = %d                                     \r\n", r, N)
	}
}

func calculateN2(arr1, arr2 []uint64) {
	max1 := arr1[len(arr1)-1]
	max2 := arr2[len(arr2)-1]
	N := uint64(firstR*2 + 1)
	for r := firstR; r <= maxR; r++ {
		// fmt.Printf("r = %d", r)
		i, maxI := 0, 0
		q := make(chan bool)
		tick := time.Tick(time.Second / 4)
		go func() {
			for {
				select {
				case <-tick:
					fmt.Printf("N = %d, iMax = %d, %d%% complete                \r", N, maxI, uint64(maxI)*100/uint64(r))
				case <-q:
					return
				}
			}
		}()
		for ; i < r; i++ {
			// fmt.Printf("\tA(i)=%d, N=%d, N-A(i) = %d-%d = %d ", A[i], N, N, A[i], uint64(N)-A[i])
			if i > maxI {
				maxI = i
			}
			if binarySearch(arr2, N-arr1[i]) {
				N++
				// fmt.Printf("N=%d\r", N)
				if N > max1 || N > max2 {
					q <- true
					time.Sleep(time.Millisecond)
					fmt.Println("\nReached max number in the table!")
					return
				}
				i = -1
				// fmt.Print("- white, N = N+1\r\n")
			}
		}
		q <- true
		time.Sleep(time.Millisecond)
		fmt.Printf("N(%d) = %d                                     \r\n", r, N)
	}
}
