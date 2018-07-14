package main

func before(a, b int) (c int) {
	if a > b {
		c = a
	} else {
		c = b
	}
	return
}

func after(a, b int) (c int) {
	c = map[bool]int{true: a, false: b}[a > b]
	return
}

func main() {
	println(before(10, 10))
	println(before(10, 20))
	println(before(20, 10))
	println(after(10, 10))
	println(after(10, 20))
	println(after(20, 10))
}
