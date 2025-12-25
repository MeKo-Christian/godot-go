package nativestructure

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/godot-go/godot-go/cmd/extensionapiparser"
)

//go:embed nativestructures.go.tmpl
var nativeStructuresText string

// Generate will generate Go wrappers for all Godot base types
func Generate(projectPath string, eapi extensionapiparser.ExtensionApi) error {
	if err := GenerateNativeStrucutres(projectPath, eapi); err != nil {
		return fmt.Errorf("native structures: %w", err)
	}
	return nil
}

func GenerateNativeStrucutres(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("nativestructures.gen.go").
		Funcs(template.FuncMap{
			"nativeStructureFormatToFields": nativeStructureFormatToFields,
			"hasPrefix":                     strings.HasPrefix,
		}).
		Parse(nativeStructuresText)
	if err != nil {
		return fmt.Errorf("parse template nativestructures.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template nativestructures.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "nativestructure", fmt.Sprintf("nativestructures.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func writeGeneratedFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
