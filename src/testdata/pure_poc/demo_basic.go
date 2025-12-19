package main

//go:pure
//go:noinline
func add(a, b int) int {
	return a + b
}

func main() {
	// This will be evaluated at compile time
	result := add(2, 3)
	println("add(2, 3) =", result)
}
