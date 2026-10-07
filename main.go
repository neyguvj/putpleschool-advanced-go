package main

import (
	"fmt"
	"math/rand"
)

func main() {
	const numbersCount = 10
	const maxNumber = 100

	numCh := make(chan int)
	squareCh := make(chan int)

	go func() {
		for range numbersCount {
			numCh <- rand.Intn(maxNumber + 1)
		}
		close(numCh)
	}()

	go func() {
		for num := range numCh {
			squareCh <- num * num
		}
		close(squareCh)
	}()

	for sq := range squareCh {
		fmt.Println(sq)
	}
}
