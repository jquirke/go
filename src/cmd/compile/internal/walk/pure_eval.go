// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"bytes"
	"fmt"
	"go/constant"
	"internal/buildcfg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cmd/compile/internal/base"
	"cmd/compile/internal/ir"
	"cmd/compile/internal/types"
)

func evaluatePureFunctionAtCompileTime(fn *ir.Func, args []ir.Node) ir.Node {
	// Don't recursively evaluate if we're inside a helper program execution
	if os.Getenv("GO_PURE_EVAL_HELPER") != "" {
		return nil
	}

	if base.Flag.LowerM > 1 {
		fmt.Printf("attempting compile-time evaluation of %s\n", ir.PkgFuncName(fn))
	}

	programSource := generateHelperProgram(fn, args)
	if programSource == "" {
		if base.Flag.LowerM > 1 {
			fmt.Printf("failed to generate helper program\n")
		}
		return nil
	}

	if base.Flag.LowerM > 1 {
		fmt.Printf("generated helper program:\n%s\n", programSource)
	}

	result := executeHelperProgram(programSource)
	if result == "" {
		if base.Flag.LowerM > 1 {
			fmt.Printf("failed to execute helper program\n")
		}
		return nil
	}

	if base.Flag.LowerM > 0 {
		fmt.Printf("compile-time evaluated %s to constant %q\n", ir.PkgFuncName(fn), result)
	}

	// Create a typed string constant
	s := ir.NewString(base.Pos, result)
	s.SetType(types.Types[types.TSTRING])
	s.SetTypecheck(1)
	return s
}

func generateHelperProgram(fn *ir.Func, args []ir.Node) string {
	pkgPath := fn.Sym().Pkg.Path
	funcName := fn.Sym().Name

	var argExprs []string
	for _, arg := range args {
		val := arg.Val()
		switch val.Kind() {
		case constant.String:
			argExprs = append(argExprs, fmt.Sprintf("%q", constant.StringVal(val)))
		case constant.Int:
			argExprs = append(argExprs, val.String())
		case constant.Float:
			argExprs = append(argExprs, val.String())
		case constant.Bool:
			argExprs = append(argExprs, fmt.Sprintf("%t", constant.BoolVal(val)))
		default:
			return ""
		}
	}

	var buf bytes.Buffer
	buf.WriteString("package main\n\n")
	buf.WriteString(fmt.Sprintf("import \"%s\"\n", pkgPath))
	buf.WriteString("import \"fmt\"\n\n")
	buf.WriteString("func main() {\n")
	buf.WriteString(fmt.Sprintf("\tresult := %s.%s(%s)\n", filepath.Base(pkgPath), funcName, strings.Join(argExprs, ", ")))
	buf.WriteString("\tfmt.Print(result)\n")
	buf.WriteString("}\n")

	return buf.String()
}

func executeHelperProgram(source string) string {
	tmpDir, err := os.MkdirTemp("", "go-pure-eval-*")
	if err != nil {
		return ""
	}
	defer os.RemoveAll(tmpDir)

	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(source), 0644); err != nil {
		return ""
	}

	// Use the go binary from the build's GOROOT
	goBin := "go"
	if buildcfg.GOROOT != "" {
		// Use the freshly built go binary from GOROOT/bin
		goBin = filepath.Join(buildcfg.GOROOT, "bin", "go")
		if _, err := os.Stat(goBin); err != nil {
			// Fall back to PATH if the GOROOT binary doesn't exist
			goBin = "go"
		}
	}

	cmd := exec.Command(goBin, "run", mainFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment with the correct GOROOT and recursion guard
	env := os.Environ()
	if buildcfg.GOROOT != "" {
		env = append(env, "GOROOT="+buildcfg.GOROOT)
	}
	// Prevent recursive evaluation
	env = append(env, "GO_PURE_EVAL_HELPER=1")
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		if base.Flag.LowerM > 1 {
			fmt.Printf("helper program error: %v\nstderr: %s\n", err, stderr.String())
		}
		return ""
	}

	return stdout.String()
}

func allArgsConstant(args []ir.Node, paramTypes []*types.Field) bool {
	if len(args) != len(paramTypes) {
		return false
	}

	for i, arg := range args {
		param := paramTypes[i]

		// Check if argument is a constant
		if arg.Op() != ir.OLITERAL {
			return false
		}

		val := arg.Val()
		kind := val.Kind()

		// Only support basic constant types
		switch kind {
		case constant.String, constant.Int, constant.Float, constant.Bool:
			// Supported
		default:
			return false
		}

		// Verify type compatibility
		switch {
		case param.Type.IsString() && kind != constant.String:
			return false
		case param.Type.IsInteger() && kind != constant.Int:
			return false
		case param.Type.IsFloat() && kind != constant.Float:
			return false
		case param.Type.IsBoolean() && kind != constant.Bool:
			return false
		}
	}
	return true
}
