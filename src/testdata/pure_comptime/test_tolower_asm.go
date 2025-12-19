package main

import "strings"

const globalConst = "GLOBAL CONSTANT"

// Case 1: Constant argument - should be evaluated at compile time
func constCase() string {
	return strings.ToLower("HELLO WORLD")
}

// Case 2: Variable argument - normal runtime call
func varCase(s string) string {
	return strings.ToLower(s)
}

// Case 3: Constant declared at package level
func packageLevelConst() string {
	return strings.ToLower(globalConst)
}

// Case 4: Multiple constant calls
func multiConst() (string, string, string) {
	a := strings.ToLower("FIRST")
	b := strings.ToLower("SECOND")
	c := strings.ToLower("THIRD")
	return a, b, c
}

// Case 5: Mixed - one constant, one variable
func mixed(s string) (string, string) {
	constant := strings.ToLower("CONSTANT")
	variable := strings.ToLower(s)
	return constant, variable
}

func main() {
	// Just call the functions so they don't get optimized away
	_ = constCase()
	_ = varCase("TEST")
	_ = packageLevelConst()
	_, _, _ = multiConst()
	_, _ = mixed("VAR")
}
