package main

import "fmt"

//go:pure
func add(a, b int) int {
	return a + b
}

//go:pure
func multiply(x, y int) int {
	return x * y
}

//go:pure
func subtract(a, b int) int {
	return a - b
}

//go:pure
func divide(a, b int) int {
	if b == 0 {
		return 0
	}
	return a / b
}

//go:pure
func bitwiseAnd(a, b int) int {
	return a & b
}

//go:pure
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== Pure Function Constant Folding Demo ===")
	fmt.Println()

	// These calls have constant arguments and will be:
	// 1. Force-inlined (because //go:pure + all args constant)
	// 2. Constant-folded by SSA SCCP pass
	// Result: No function call in assembly, just constants!

	fmt.Println("Constant argument cases (fully optimized):")
	fmt.Printf("  add(10, 20) = %d\n", add(10, 20))
	fmt.Printf("  multiply(7, 8) = %d\n", multiply(7, 8))
	fmt.Printf("  subtract(100, 42) = %d\n", subtract(100, 42))
	fmt.Printf("  divide(100, 5) = %d\n", divide(100, 5))
	fmt.Printf("  bitwiseAnd(0xFF, 0x0F) = %d\n", bitwiseAnd(0xFF, 0x0F))
	fmt.Printf("  max(42, 17) = %d\n", max(42, 17))
	fmt.Println()

	// These calls have variable arguments and will be:
	// 1. Inlined normally (small functions, within budget)
	// 2. NOT constant-folded (arguments not known at compile time)
	// Result: Inlined code in assembly

	x, y := 5, 3
	fmt.Println("Variable argument cases (inlined but not folded):")
	fmt.Printf("  add(x, y) = %d\n", add(x, y))
	fmt.Printf("  multiply(x, y) = %d\n", multiply(x, y))
	fmt.Printf("  subtract(x, y) = %d\n", subtract(x, y))
	fmt.Println()

	fmt.Println("To see the optimization in action:")
	fmt.Println("  go build -gcflags=\"-m\" demo.go")
	fmt.Println("  # Look for: 'inlining pure function ... with constant arguments for SSA folding'")
	fmt.Println()
	fmt.Println("  go tool compile -S demo.go | grep -A 5 'add(10, 20)'")
	fmt.Println("  # You'll see MOVD $30 instead of a CALL instruction")
}
