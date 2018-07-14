package main

import (
	"fmt"
	"math/rand"
	"time"
)

func f(c chan int) {
	i, j := 0, 1
	for {
		c <- j
		i, j = j, i+j
	}
}

func main() {
	const n = 50

	c := make(chan int)
	go f(c)

	rand.Seed(time.Now().UnixNano())
	s := time.Now()

	for i := 0; i < n; i++ {
		fmt.Println(<-c)
	}
	e := time.Since(s)
	fmt.Println(e)
}
