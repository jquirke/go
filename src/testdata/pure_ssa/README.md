# Pure Function SSA Constant Folding Demo

This directory contains demonstrations of the `//go:pure` pragma with SSA-based constant folding.

## How It Works

When a function is marked `//go:pure` and called with constant arguments:

1. **Force Inline**: The compiler forces inlining (regardless of cost budget)
2. **SSA SCCP**: The SSA's Sparse Conditional Constant Propagation pass folds the constants
3. **Result**: Function call is completely eliminated, replaced with the constant result

## Architecture

This is a **scalable** approach that leverages existing compiler infrastructure:

- Uses SSA's proven constant folding (SCCP pass in `cmd/compile/internal/ssa/sccp.go`)
- No custom interpreter needed - reuses generic rewrite rules
- Handles all operations SCCP already supports (arithmetic, bitwise, comparisons, etc.)
- Only ~20 lines of new code in `cmd/compile/internal/inline/inl.go`

## Demo Files

### simple.go

Minimal example showing basic constant folding:

```go
//go:pure
func add(a, b int) int {
    return a + b
}

x := add(10, 20)  // Folded to: x = 30
```

**Build and inspect:**
```bash
go build -gcflags="-m" simple.go
# Shows: "inlining pure function add with constant arguments for SSA folding"

go tool compile -S simple.go | grep "MOVD.*\$30"
# Shows: MOVD $30, R0  (constant 30, no function call!)
```

### demo.go

Comprehensive example with multiple pure functions:
- Arithmetic: `add`, `multiply`, `subtract`, `divide`
- Bitwise: `bitwiseAnd`
- Control flow: `max` (with if statement)

Shows both constant and variable argument cases.

## Verification

### 1. Compiler Messages

```bash
$ go build -gcflags="-m" simple.go
./simple.go:10:10: inlining pure function add with constant arguments for SSA folding
```

### 2. Assembly Inspection

**Before optimization** (typical function call):
```assembly
MOVD a, R0
MOVD b, R1
CALL add
```

**After optimization** (constant folded):
```assembly
MOVD $30, R0    ; Just the constant!
```

### 3. Runtime Verification

```bash
$ go run simple.go
30    # add(10, 20) = 30
12    # add(5, 7) = 12
```

## Supported Operations

Since this leverages SSA SCCP, all operations in `possibleConst()` are supported:

**Arithmetic**: Add, Sub, Mul, Div, Mod (signed and unsigned)
**Bitwise**: And, Or, Xor, Lsh, Rsh
**Comparison**: Eq, Lt, Le, Gt, Ge
**Conversions**: Integer casts, float conversions
**Unary**: Neg, Com (complement), Not
**Math**: Floor, Ceil, Trunc, Sqrt

See `cmd/compile/internal/ssa/sccp.go:possibleConst()` for the complete list.

## Limitations

1. **Trust-based**: Does not verify purity - trusts the `//go:pure` pragma
2. **Constant arguments only**: At least one argument must be constant
3. **Inlining constraints**: Function must be inlineable (no recursion, etc.)

## Future Work

- Purity verification (static analysis or runtime checks)
- Memoization for pure functions with runtime-known arguments
- Cross-function constant propagation for pure function chains

## Implementation Details

### Code Location

- Pragma definition: `cmd/compile/internal/ir/node.go`
- Pragma parsing: `cmd/compile/internal/noder/lex.go`
- Inline forcing: `cmd/compile/internal/inline/inl.go:TryInlineCall()`
- Constant folding: `cmd/compile/internal/ssa/sccp.go` (existing infrastructure)

### Key Function

```go
// Force inline pure functions with constant arguments to enable SSA constant folding
if fn.Pragma&ir.Pure != 0 && allConstantArgs(call) {
    if base.Flag.LowerM > 0 {
        fmt.Printf("%v: inlining pure function %v with constant arguments for SSA folding\n",
            ir.Line(call), fn.Sym().Name)
    }
    return mkinlcall(callerfn, call, fn, bigCaller, closureCalledOnce)
}
```

## Comparison to POC

The previous POC (`testdata/pure_poc/`) used a custom interpreter:
- ❌ Unscalable - needed to implement every operation
- ❌ Error-prone - easy to introduce bugs in constant evaluation
- ❌ Limited - only handled simple expressions

This SSA approach:
- ✅ Scalable - reuses existing SCCP infrastructure
- ✅ Correct - proven through existing SSA optimizations
- ✅ Comprehensive - handles all operations SCCP supports
- ✅ Simple - minimal new code

## References

- SSA SCCP: Wegman & Zadeck, "Constant Propagation with Conditional Branches", TOPLAS 1991
- Go SSA: `cmd/compile/internal/ssa/`
- Inlining: `cmd/compile/internal/inline/`
