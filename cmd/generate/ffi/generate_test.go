package ffi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/godot-go/godot-go/cmd/gdextensionparser/clang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// findProjectRoot walks up from the current directory to find the project root
// by looking for go.mod file.
func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (no go.mod found)")
		}
		dir = parent
	}
}

func TestGenerate(t *testing.T) {
	projectPath := findProjectRoot(t)
	ast := clang.CHeaderFileAST{
		Expr: []clang.Expr{},
	}
	var panicFunc assert.PanicTestFunc = func() {
		Generate(projectPath, ast)
	}
	require.NotPanics(t, panicFunc)
}
