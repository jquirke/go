# Assembly Comparison: ToLower with Constant vs Variable Arguments

## Summary of Compile-Time Evaluation

All constant calls to `strings.ToLower()` were evaluated at compile time:
- `strings.ToLower("HELLO WORLD")` → `"hello world"`
- `strings.ToLower(globalConst)` → `"global constant"` (package-level const)
- `strings.ToLower("FIRST")` → `"first"`
- `strings.ToLower("SECOND")` → `"second"`
- `strings.ToLower("THIRD")` → `"third"`
- `strings.ToLower("CONSTANT")` → `"constant"`

Variable calls remain as runtime function calls.

---

## Case 1: Constant Literal - `constCase()`

**Source:**
```go
func constCase() string {
    return strings.ToLower("HELLO WORLD")
}
```

**Assembly (ARM64):**
```asm
main.constCase STEXT size=16 args=0x0 locals=0x0 funcid=0x0 align=0x0 leaf
    TEXT   main.constCase(SB), LEAF|NOFRAME|ABIInternal, $0-0
    MOVD   $go:string."hello world"(SB), R0    // Load constant address
    MOVD   $11, R1                              // Length = 11
    RET    (R30)                                // Return
```

**Analysis:**
- **NO function call** to `strings.ToLower()`
- Just 3 instructions: load string pointer, load length, return
- `LEAF|NOFRAME` indicates no stack frame needed
- Total size: **16 bytes**

---

## Case 2: Variable Argument - `varCase(s string)`

**Source:**
```go
func varCase(s string) string {
    return strings.ToLower(s)
}
```

**Assembly (ARM64):**
```asm
main.varCase STEXT size=64 args=0x10 locals=0x18 funcid=0x0 align=0x0
    TEXT      main.varCase(SB), ABIInternal, $32-16
    MOVD      16(g), R16              // Stack check
    CMP       R16, RSP
    BLS       44
    MOVD.W    R30, -32(RSP)          // Save return address
    MOVD      R29, -8(RSP)           // Save frame pointer
    SUB       $8, RSP, R29
    MOVD      R0, main.s(FP)
    CALL      strings.ToLower(SB)     // **ACTUAL FUNCTION CALL**
    MOVD      -8(RSP), R29
    MOVD.P    32(RSP), R30
    RET       (R30)
    ...stack growth code...
```

**Analysis:**
- **Full function call** to `strings.ToLower()`
- Requires stack frame setup ($32 bytes)
- Stack overflow check
- Multiple memory operations
- Total size: **64 bytes** (4x larger than constant case)

---

## Case 3: Package-Level Const - `packageLevelConst()`

**Source:**
```go
const globalConst = "GLOBAL CONSTANT"

func packageLevelConst() string {
    return strings.ToLower(globalConst)
}
```

**Assembly (ARM64):**
```asm
main.packageLevelConst STEXT size=16 args=0x0 locals=0x0 funcid=0x0 align=0x0 leaf
    TEXT   main.packageLevelConst(SB), LEAF|NOFRAME|ABIInternal, $0-0
    MOVD   $go:string."global constant"(SB), R0   // Load constant
    MOVD   $15, R1                                 // Length = 15
    RET    (R30)
```

**Analysis:**
- Package-level const **also optimized** at compile time!
- Identical to Case 1: no function call, just load constant
- Proves compile-time evaluation works with `const` declarations
- Total size: **16 bytes**

---

## Case 4: Multiple Constants - `multiConst()`

**Source:**
```go
func multiConst() (string, string, string) {
    a := strings.ToLower("FIRST")
    b := strings.ToLower("SECOND")
    c := strings.ToLower("THIRD")
    return a, b, c
}
```

**Assembly (ARM64):**
```asm
main.multiConst STEXT size=48 args=0x0 locals=0x0 funcid=0x0 align=0x0 leaf
    TEXT   main.multiConst(SB), LEAF|NOFRAME|ABIInternal, $0-0
    MOVD   $go:string."first"(SB), R0      // Result 1
    MOVD   $5, R1
    MOVD   $go:string."second"(SB), R2     // Result 2
    MOVD   $6, R3
    MOVD   $go:string."third"(SB), R4      // Result 3
    MOVD   R1, R5
    RET    (R30)
```

**Analysis:**
- **All 3 calls** eliminated at compile time
- Just loads 3 constant string pointers and lengths
- No function calls, no stack frame
- Total size: **48 bytes** (still smaller than 1 runtime call!)

---

## Case 5: Mixed (Constant + Variable) - `mixed(s string)`

**Source:**
```go
func mixed(s string) (string, string) {
    constant := strings.ToLower("CONSTANT")
    variable := strings.ToLower(s)
    return constant, variable
}
```

**Assembly (ARM64):**
```asm
main.mixed STEXT size=96 args=0x10 locals=0x18 funcid=0x0 align=0x0
    ...stack setup...
    CALL      strings.ToLower(SB)        // Call for variable argument
    MOVD      R0, R2                     // Save variable result
    MOVD      R1, R3
    MOVD      $go:string."constant"(SB), R0   // Load constant (no call!)
    MOVD      $8, R1
    ...return...
```

**Analysis:**
- **Only 1 function call** for the variable argument
- Constant argument optimized away (directly loads `"constant"`)
- Demonstrates selective optimization
- Total size: **96 bytes**

---

## Performance Comparison

| Case | Function Calls | Instructions | Size | Stack Frame |
|------|---------------|--------------|------|-------------|
| constCase() | 0 | ~3 | 16 bytes | None (NOFRAME) |
| varCase(s) | 1 | ~16 | 64 bytes | 32 bytes |
| packageLevelConst() | 0 | ~3 | 16 bytes | None (NOFRAME) |
| multiConst() | 0 | ~7 | 48 bytes | None (NOFRAME) |
| mixed(s) | 1 | ~24 | 96 bytes | 32 bytes |

**Key Takeaways:**
1. Constant case is **4x smaller** than variable case (16 vs 64 bytes)
2. **Zero runtime overhead** for constant arguments - no function call at all
3. **Package-level const works** - fully evaluated at compile time
4. **Multiple constants** still much cheaper than single runtime call
5. **Mixed scenarios** optimize constants while keeping runtime calls for variables

---

## String Constants Generated

The compiler generated these string constants in the `.rodata` section:

```
go:string."hello world"      (11 bytes)
go:string."global constant"  (15 bytes)
go:string."first"             (5 bytes)
go:string."second"            (6 bytes)
go:string."third"             (5 bytes)
go:string."constant"          (8 bytes)
```

These are the **pre-computed results** that replaced the function calls!
