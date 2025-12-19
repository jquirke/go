# strings.ToLower SSA Constant Folding Optimization

## Summary

Successfully implemented compile-time constant folding for `strings.ToLower()` when called with ASCII string literals. The optimization eliminates the function call entirely and replaces it with the pre-computed lowercase constant.

## Implementation Details

**Location**: Go compiler SSA rewrite rules
**Files Modified**:
- `/Users/qjeremy/go/github.com/golang/go/src/cmd/compile/internal/ssa/_gen/generic.rules` (SSA rewrite rules)
- `/Users/qjeremy/go/github.com/golang/go/src/cmd/compile/internal/ssa/rewrite.go` (helper functions)

**Key Technique**: Pattern matching on `StaticLECall` in Generic SSA stage, before string decomposition and call expansion.

## Test Cases and Assembly Output

### Test Program
```go
package main

import (
	"fmt"
	"os"
	"strings"
)
const (
	_upperCase = "UPPERCASE"
)

func main() {
	// Case 1: String literal - should be folded at compile time
	s1 := strings.ToLower("HELLO WORLD")
	fmt.Printf("Constant: %q\n", s1)

	// Case 2: Variable from args - cannot be folded
	if len(os.Args) > 1 {
		s2 := strings.ToLower(os.Args[1])
		fmt.Printf("Arg: %q\n", s2)
	}

	// Case 3: Another constant - should be folded
	s3 := strings.ToLower(_upperCase)
	fmt.Printf("Constant 2: %q\n", s3)
}
```

### Assembly Output Analysis

#### Case 1: String Literal `"HELLO WORLD"` → **OPTIMIZED** ✅
```asm
0x001c 00028 (/tmp/test_tolower_final.go:12)	MOVD	$go:string."hello world"(SB), R0
0x0024 00036 (/tmp/test_tolower_final.go:12)	MOVD	$11, R1
```
**Result**: Direct constant load. **NO function call**. The compiler computed `"hello world"` at compile time.

#### Case 2: Variable `os.Args[1]` → **NOT OPTIMIZED** (correct behavior)
```asm
0x0078 00120 (/tmp/test_tolower_final.go:16)	LDP	16(R2), (R0, R1)
0x007c 00124 (/tmp/test_tolower_final.go:16)	CALL	strings.ToLower(SB)
```
**Result**: Runtime function call. **Cannot be folded** because the input is not known at compile time.

#### Case 3: String Literal `"UPPERCASE"` → **OPTIMIZED** ✅
```asm
0x00c4 00196 (/tmp/test_tolower_final.go:22)	MOVD	$go:string."uppercase"(SB), R0
0x00cc 00204 (/tmp/test_tolower_final.go:22)	MOVD	$9, R1
```
**Result**: Direct constant load. **NO function call**. The compiler computed `"uppercase"` at compile time.

## Performance Impact

### Before Optimization
- All 3 cases: **3 function calls** to `strings.ToLower(SB)`
- Each call involves stack frame setup, argument passing, loop execution, and return

### After Optimization
- Case 1 & 3: **0 function calls** (folded to constants)
- Case 2: **1 function call** (variable input, must be runtime)
- **Total: 67% reduction in function calls** for this example

### Instruction Count Comparison

**Case 1 (Literal) - Before:**
```asm
MOVD    $go:string."HELLO WORLD"(SB), R0  ; Load input pointer
MOVD    $11, R1                            ; Load input length
CALL    strings.ToLower(SB)                ; Function call (~20+ instructions)
; ... ToLower function body execution ...
MOVD    R0, result_ptr                     ; Store result pointer
MOVD    R1, result_len                     ; Store result length
```

**Case 1 (Literal) - After:**
```asm
MOVD    $go:string."hello world"(SB), R0  ; Load pre-computed result
MOVD    $11, R1                            ; Load length
```

**Savings**: ~20+ instructions eliminated per constant call.

## Limitations

The optimization only applies when:
1. ✅ Input is a **compile-time string constant**
2. ✅ Input contains **only ASCII characters** (no Unicode)
3. ❌ Does NOT apply to variables, function results, or non-ASCII strings

## Technical Implementation

The SSA rewrite rules match the following pattern:

```
(SelectN [0] call:(StaticLECall {callAux} (StringMake (Addr {s} _) _) mem))
  && isSameCall(callAux, "strings.ToLower")
  && isASCIIGoString(s)
  && clobber(call)
=> (ConstString <typ.String> {auxToString(goStringToLower(s))})
```

This pattern:
1. Matches `SelectN [0]` (the string result) from a `StaticLECall` to `strings.ToLower`
2. Extracts the constant string from the `go:string.XXX` symbol
3. Verifies it's ASCII-only
4. Marks the call as clobbered (dead)
5. Replaces the result with a pre-computed lowercase `ConstString`

The memory result (`SelectN [1]`) is similarly replaced with the original memory state, allowing dead code elimination to remove the entire call.

## Verification

### Runtime Output
```
$ go run test_tolower_final.go FOO
Constant: "hello world"
Arg: "foo"
Constant 2: "uppercase"
```

### Assembly Verification
```bash
$ go build -gcflags="-S" test_tolower_final.go 2>&1 | grep -c "CALL.*ToLower"
1
```
Only **1 call** to `ToLower` in the entire program (the variable case).

### Constant Verification
```bash
$ go build -gcflags="-S" test_tolower_final.go 2>&1 | grep "hello world\|uppercase"
MOVD    $go:string."hello world"(SB), R0
MOVD    $go:string."uppercase"(SB), R0
```
Both constants appear directly in the assembly as lowercase, with no function calls.

## Conclusion

The optimization successfully eliminates compile-time overhead for constant `strings.ToLower()` calls while preserving correct runtime behavior for variable inputs. This demonstrates the power of SSA-based constant folding in modern compilers.
