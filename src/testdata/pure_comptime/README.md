# Compile-Time Execution Demos for Pure Functions

This directory contains demonstration programs showing compile-time execution
of pure functions with constant arguments.

## Files

### test_tolower_asm.go
Comprehensive test showing assembly comparison across 5 different cases:
- Case 1: Constant literal `"HELLO WORLD"` (optimized)
- Case 2: Variable argument (runtime call)
- Case 3: Package-level const (optimized)
- Case 4: Multiple constant calls (all optimized)
- Case 5: Mixed constant + variable (selective optimization)

**Usage:**
```bash
# Build and see compile-time evaluation messages
go build -gcflags="-m" test_tolower_asm.go

# Generate assembly to verify optimization
go build -gcflags="-S" test_tolower_asm.go 2>&1 | less

# Run the program
go run test_tolower_asm.go
```

**Expected output with -gcflags="-m":**
```
compile-time evaluated strings.ToLower to constant "hello world"
compile-time evaluated strings.ToLower to constant "global constant"
compile-time evaluated strings.ToLower to constant "first"
compile-time evaluated strings.ToLower to constant "second"
compile-time evaluated strings.ToLower to constant "third"
compile-time evaluated strings.ToLower to constant "constant"
```

### test_verify.go
Simple verification program that demonstrates and validates compile-time
execution with user-friendly output.

**Usage:**
```bash
go run test_verify.go
```

**Expected output:**
```
=== Compile-Time Evaluation Demo ===

Case 1: ToLower("HELLO WORLD"):
  Result: "hello world"
  [Compile-time evaluated - no function call!]

Case 2: ToLower(variable):
  Input:  "RUNTIME CALL"
  Result: "runtime call"
  [Runtime function call]

Case 3: ToLower(packageConst):
  Result: "package level"
  [Compile-time evaluated - const works!]

=== Verification ===
✓ Case 1: Correct
✓ Case 2: Correct
✓ Case 3: Correct
```

## Assembly Comparison

See `cmd/compile/internal/walk/asm_comparison.md` for detailed assembly
analysis showing:

- **Constant case**: 16 bytes, 3 instructions, no function call
- **Variable case**: 64 bytes, 16 instructions, full function call
- **Performance**: 4x smaller code, zero runtime overhead

## Implementation Details

See `cmd/compile/internal/walk/COMPILE_TIME_EXEC_SUMMARY.md` for complete
implementation documentation.

## How It Works

1. Compiler detects `strings.ToLower("HELLO")` with constant argument
2. Generates temporary Go program calling the function
3. Executes: `go run /tmp/go-pure-eval-*/main.go`
4. Captures output: `"hello"`
5. Replaces call with constant in AST
6. Emits assembly: `MOVD $go:string."hello"(SB), R0`

## Requirements

- Built Go toolchain (compile-time execution uses `GOROOT/bin/go`)
- `//go:pure` pragma on function
