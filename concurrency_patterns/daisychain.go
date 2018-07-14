package main

import (
	"fmt"
	"math/rand"
	"time"
)

func w(l, r chan int) {
	l <- 1 + <-r
}

func main() {
	const n = 99999
	le := make(chan int)
	l := le
	r := le

	rand.Seed(time.Now().UnixNano())
	s := time.Now()

	for i := 0; i < n; i++ {
		r = make(chan int)
		go w(l, r)
		l = r
	}
	e := time.Since(s)
	go func(c chan int) { c <- 1 }(r)

	fmt.Println(e)
}
