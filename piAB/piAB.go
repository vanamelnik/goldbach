package main

import (
	"bufio"
	"fmt"
	"goldbach/files"
	"log"
	"os"
	"time"
)

const (
	max      = 1000000000
	fABin    = "tables/A.gz"
	fBBin    = "tables/B.gz"
	filename = "piB.txt"
)

func main() {
	fmt.Println("Загрузка файла", fBBin)
	B, err := files.UnzipFile(fBBin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Вычисление PiB(x) для x от 1 до ", max)
	aPrimes := getAPrimes(B, max)
	fmt.Println("Запись данных в файл ", filename)
	saveAPrimes(aPrimes, filename)
}

func getAPrimes(indexes []uint64, max int) []uint64 {
	aPrimes := make([]uint64, max+1)
	index := uint64(0)
	i := uint64(5)
	prime := indexes[index]*6 - 1
	q := make(chan bool)
	t := time.Tick(time.Second * 2)
	go func() {
		for {
			select {
			case <-t:
				fmt.Printf("i=%d, %d%% complete              \r", i, int(float64(i)/float64(max)*100))
			case <-q:
				return
			}
		}
	}()
	for ; i <= uint64(max); i++ {
		aPrimes[i] = index
		// fmt.Printf("i=%v\tindex=%v\tnext prime=%v\n", i, index, prime)
		if i == prime {
			index++
			prime = indexes[index]*6 - 1
		}
	}
	q <- true
	time.Sleep(time.Millisecond)
	fmt.Print("                                                 \r")
	return aPrimes
}

func saveAPrimes(aPrimes []uint64, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, num := range aPrimes {
		_, err = w.WriteString(fmt.Sprintf("%d\r\n", num))
	}
	if err = w.Flush(); err != nil {
		log.Fatal(err)
	}
}
