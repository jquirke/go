package main

//go:pure
//go:noinline
func add(a, b int) int {
	return a + b
}

//go:pure
//go:noinline
func multiply(x, y int) int {
	return x * y
}

//go:pure
//go:noinline
func subtract(a, b int) int {
	return a - b
}

func main() {
	// These should be evaluated at compile time
	r1 := add(2, 3)
	r2 := multiply(4, 5)
	r3 := subtract(10, 3)

	println("add(2, 3) =", r1)
	println("multiply(4, 5) =", r2)
	println("subtract(10, 3) =", r3)

	// This should NOT be evaluated (non-constant)
	x := 10
	r4 := add(x, 5)
	println("add(x, 5) =", r4)
}
