package gdclassinit

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
	//go:embed classes.init.go.tmpl
	classesInitText string

	//go:embed classes.callbacks.h.tmpl
	cHeaderClassesText string

	//go:embed classes.callbacks.c.tmpl
	cClassesText string
)

// Generate will generate Go wrappers for all Godot base types
func Generate(projectPath string, eapi extensionapiparser.ExtensionApi) error {
	if err := GenerateCHeaderClassCallbacks(projectPath, eapi); err != nil {
		return fmt.Errorf("class callbacks header: %w", err)
	}
	if err := GenerateCClassCallbacks(projectPath, eapi); err != nil {
		return fmt.Errorf("class callbacks c: %w", err)
	}
	if err := GenerateClassInit(projectPath, eapi); err != nil {
		return fmt.Errorf("class init: %w", err)
	}
	return nil
}

func GenerateClassInit(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.init.gen.go").
		Funcs(template.FuncMap{
			"goVariantConstructor": goVariantConstructor,
			"goMethodName":         goMethodName,
			"goArgumentName":       goArgumentName,
			"goArgumentType":       goArgumentType,
			"goReturnType":         goReturnType,
			"goClassEnumName":      goClassEnumName,
			"goClassStructName":    goClassStructName,
			"goClassInterfaceName": goClassInterfaceName,
			"coalesce":             coalesce,
		}).
		Parse(classesInitText)
	if err != nil {
		return fmt.Errorf("parse template classes.init.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.init.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "gdclassinit", fmt.Sprintf("classes.init.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateCHeaderClassCallbacks(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.callbacks.gen.h").
		Funcs(template.FuncMap{
			"goMethodName":    goMethodName,
			"goArgumentName":  goArgumentName,
			"goArgumentType":  goArgumentType,
			"goReturnType":    goReturnType,
			"goClassEnumName": goClassEnumName,
			"coalesce":        coalesce,
		}).
		Parse(cHeaderClassesText)
	if err != nil {
		return fmt.Errorf("parse template classes.callbacks.gen.h: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.callbacks.gen.h: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "gdclassinit", fmt.Sprintf("classes.callbacks.gen.h"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateCClassCallbacks(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.callbacks.gen.c").
		Funcs(template.FuncMap{
			"goMethodName":    goMethodName,
			"goArgumentName":  goArgumentName,
			"goArgumentType":  goArgumentType,
			"goReturnType":    goReturnType,
			"goClassEnumName": goClassEnumName,
			"coalesce":        coalesce,
		}).
		Parse(cClassesText)
	if err != nil {
		return fmt.Errorf("parse template classes.callbacks.gen.c: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.callbacks.gen.c: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "gdclassinit", fmt.Sprintf("classes.callbacks.gen.c"))

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
