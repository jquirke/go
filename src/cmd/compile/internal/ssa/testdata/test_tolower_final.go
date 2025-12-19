package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Constant - should be folded at compile time
	s1 := strings.ToLower("HELLO WORLD")
	fmt.Printf("Constant: %q\n", s1)

	// Variable from args - cannot be folded
	if len(os.Args) > 1 {
		s2 := strings.ToLower(os.Args[1])
		fmt.Printf("Arg: %q\n", s2)
	}

	// Another constant - should be folded
	s3 := strings.ToLower("UPPERCASE")
	fmt.Printf("Constant 2: %q\n", s3)
}
