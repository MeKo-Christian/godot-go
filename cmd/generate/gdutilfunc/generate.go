package gdutilfunc

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/godot-go/godot-go/cmd/extensionapiparser"
	"github.com/iancoleman/strcase"
)

//go:embed utilityfunctions.go.tmpl
var utilityFunctionsText string

// Generate will generate Go wrappers for all Godot base types
func Generate(projectPath string, eapi extensionapiparser.ExtensionApi) error {
	if err := GenerateUtilityFunctions(projectPath, eapi); err != nil {
		return fmt.Errorf("utility functions: %w", err)
	}
	return nil
}

func GenerateUtilityFunctions(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("utilityfunctions.gen.go").
		Funcs(template.FuncMap{
			"camelCase":           strcase.ToCamel,
			"goArgumentName":      goArgumentName,
			"goArgumentType":      goArgumentType,
			"goEncoder":           goEncoder,
			"goReturnType":        goReturnType,
			"coalesce":            coalesce,
			"goEncodeIsReference": goEncodeIsReference,
		}).
		Parse(utilityFunctionsText)
	if err != nil {
		return fmt.Errorf("parse template utilityfunctions.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template utilityfunctions.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "gdutilfunc", fmt.Sprintf("utilityfunctions.gen.go"))

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
