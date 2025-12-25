package gdextensionparser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func TestGenerateGDExtensionInterfaceAST(t *testing.T) {
	projectPath := findProjectRoot(t)
	f, err := GenerateGDExtensionInterfaceAST(projectPath, "")
	require.NoError(t, err)
	spew.Dump(f)
}
