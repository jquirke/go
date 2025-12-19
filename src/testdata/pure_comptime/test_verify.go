package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== Compile-Time Evaluation Demo ===")
	fmt.Println()

	// Case 1: Constant - evaluated at compile time
	fmt.Println("Case 1: ToLower(\"HELLO WORLD\"):")
	result1 := strings.ToLower("HELLO WORLD")
	fmt.Printf("  Result: %q\n", result1)
	fmt.Println("  [Compile-time evaluated - no function call!]")
	fmt.Println()

	// Case 2: Variable - runtime call
	fmt.Println("Case 2: ToLower(variable):")
	input := "RUNTIME CALL"
	result2 := strings.ToLower(input)
	fmt.Printf("  Input:  %q\n", input)
	fmt.Printf("  Result: %q\n", result2)
	fmt.Println("  [Runtime function call]")
	fmt.Println()

	// Case 3: Package-level const
	const packageConst = "PACKAGE LEVEL"
	fmt.Println("Case 3: ToLower(packageConst):")
	result3 := strings.ToLower(packageConst)
	fmt.Printf("  Result: %q\n", result3)
	fmt.Println("  [Compile-time evaluated - const works!]")
	fmt.Println()

	// Verify correctness
	fmt.Println("=== Verification ===")
	if result1 != "hello world" {
		fmt.Printf("ERROR: Expected 'hello world', got %q\n", result1)
	} else {
		fmt.Println("✓ Case 1: Correct")
	}

	if result2 != "runtime call" {
		fmt.Printf("ERROR: Expected 'runtime call', got %q\n", result2)
	} else {
		fmt.Println("✓ Case 2: Correct")
	}

	if result3 != "package level" {
		fmt.Printf("ERROR: Expected 'package level', got %q\n", result3)
	} else {
		fmt.Println("✓ Case 3: Correct")
	}
}
