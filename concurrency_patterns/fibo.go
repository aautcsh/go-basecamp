package main

import (
	"fmt"
	"math/rand"
	"time"
)

func f(n int) int {
	if n <= 2 {
		return n
	} else {
		return f(n-2) + f(n-1)
	}
}

func main() {
	const n = 45

	rand.Seed(time.Now().UnixNano())
	s := time.Now()

	for i := 0; i < n; i++ {
		fmt.Println(f(i))
	}
	e := time.Since(s)
	fmt.Println(e)
}
