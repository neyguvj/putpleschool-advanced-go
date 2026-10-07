package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	const numbersCount = 10
	const maxNumber = 100

	numCh := make(chan int, numbersCount)
	squareCh := make(chan int, numbersCount)
	wg := sync.WaitGroup{}
	wg.Go(func() {
		for range numbersCount {
			numCh <- rand.Intn(maxNumber)
		}
		close(numCh)
	})
	wg.Go(func() {
		for num := range numCh {
			squareCh <- num * num
		}
		close(squareCh)
	})
	wg.Wait()

	for sq := range squareCh {
		fmt.Println(sq)
	}
}
