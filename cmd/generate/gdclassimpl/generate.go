package gdclassimpl

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/godot-go/godot-go/cmd/extensionapiparser"
)

var (
	//go:embed classes.go.tmpl
	classesText string

	//go:embed classes.refs.go.tmpl
	classesRefsText string
)

// Generate will generate Go wrappers for all Godot base types
func Generate(projectPath string, eapi extensionapiparser.ExtensionApi) error {
	if err := GenerateClasses(projectPath, eapi); err != nil {
		return fmt.Errorf("classes: %w", err)
	}
	if err := GenerateClassRefs(projectPath, eapi); err != nil {
		return fmt.Errorf("class refs: %w", err)
	}
	return nil
}

func GenerateClasses(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.gen.go").
		Funcs(template.FuncMap{
			"isSetterMethodName":   isSetterMethodName,
			"goVariantConstructor": goVariantConstructor,
			"goMethodName":         goMethodName,
			"goArgumentName":       goArgumentName,
			"goArgumentType":       goArgumentType,
			"goVariantFunc":        goVariantFunc,
			"goReturnType":         goReturnType,
			"goClassEnumName":      goClassEnumName,
			"goClassStructName":    goClassStructName,
			"goClassInterfaceName": goClassInterfaceName,
			"goEncoder":            goEncoder,
			"goEncodeIsReference":  goEncodeIsReference,
			"coalesce":             coalesce,
		}).
		Parse(classesText)
	if err != nil {
		return fmt.Errorf("parse template classes.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "gdclassimpl", fmt.Sprintf("classes.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateClassRefs(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.refs.gen.go").
		Funcs(template.FuncMap{
			"isSetterMethodName":   isSetterMethodName,
			"goVariantConstructor": goVariantConstructor,
			"goMethodName":         goMethodName,
			"goArgumentName":       goArgumentName,
			"goArgumentType":       goArgumentType,
			"goVariantFunc":        goVariantFunc,
			"goReturnType":         goReturnType,
			"goClassEnumName":      goClassEnumName,
			"goClassStructName":    goClassStructName,
			"goClassInterfaceName": goClassInterfaceName,
			"goEncoder":            goEncoder,
			"goEncodeIsReference":  goEncodeIsReference,
			"coalesce":             coalesce,
		}).
		Parse(classesRefsText)
	if err != nil {
		return fmt.Errorf("parse template classes.refs.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.refs.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "gdclassimpl", fmt.Sprintf("classes.refs.gen.go"))

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
