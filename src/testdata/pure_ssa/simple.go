package main

//go:pure
func add(a, b int) int {
	return a + b
}

func main() {
	// Constant arguments - should be folded to 30
	x := add(10, 20)
	println(x)

	// Variable arguments - still a function call
	y := 5
	z := add(y, 7)
	println(z)
}
