package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("=== Compile-Time Evaluation - Primitive Types ===")
	fmt.Println()

	// String constants
	fmt.Println("String constants:")
	s1 := strings.ToLower("HELLO")
	fmt.Printf("  strings.ToLower(\"HELLO\") = %q\n", s1)

	// Int constants
	fmt.Println("\nInt constants:")
	s2 := strconv.Itoa(42)
	fmt.Printf("  strconv.Itoa(42) = %q\n", s2)

	s3 := strconv.Itoa(-100)
	fmt.Printf("  strconv.Itoa(-100) = %q\n", s3)

	// Variable arguments (runtime calls)
	fmt.Println("\nRuntime calls (variables):")
	str := "VARIABLE"
	s4 := strings.ToLower(str)
	fmt.Printf("  strings.ToLower(str) = %q\n", s4)

	num := 99
	s5 := strconv.Itoa(num)
	fmt.Printf("  strconv.Itoa(num) = %q\n", s5)

	// Verification
	fmt.Println("\n=== Verification ===")
	allCorrect := true

	if s1 != "hello" {
		fmt.Printf("✗ String constant failed: expected \"hello\", got %q\n", s1)
		allCorrect = false
	} else {
		fmt.Println("✓ String constant: Correct")
	}

	if s2 != "42" {
		fmt.Printf("✗ Int constant (42) failed: expected \"42\", got %q\n", s2)
		allCorrect = false
	} else {
		fmt.Println("✓ Int constant (42): Correct")
	}

	if s3 != "-100" {
		fmt.Printf("✗ Int constant (-100) failed: expected \"-100\", got %q\n", s3)
		allCorrect = false
	} else {
		fmt.Println("✓ Int constant (-100): Correct")
	}

	if s4 != "variable" {
		fmt.Printf("✗ String variable failed: expected \"variable\", got %q\n", s4)
		allCorrect = false
	} else {
		fmt.Println("✓ String variable: Correct")
	}

	if s5 != "99" {
		fmt.Printf("✗ Int variable failed: expected \"99\", got %q\n", s5)
		allCorrect = false
	} else {
		fmt.Println("✓ Int variable: Correct")
	}

	if allCorrect {
		fmt.Println("\n✅ All tests passed!")
	} else {
		fmt.Println("\n❌ Some tests failed")
	}
}
