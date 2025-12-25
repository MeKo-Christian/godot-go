package constant

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
	//go:embed classes.constants.go.tmpl
	classesConstantsText string

	//go:embed classes.enums.go.tmpl
	classesEnumsText string

	//go:embed globalconstants.go.tmpl
	globalConstantsText string

	//go:embed globalenums.go.tmpl
	globalEnumsText string
)

// Generate will generate Go wrappers for all Godot base types
func Generate(projectPath string, eapi extensionapiparser.ExtensionApi) error {
	if err := GenerateClassConstants(projectPath, eapi); err != nil {
		return fmt.Errorf("class constants: %w", err)
	}
	if err := GenerateClassEnums(projectPath, eapi); err != nil {
		return fmt.Errorf("class enums: %w", err)
	}
	if err := GenerateGlobalConstants(projectPath, eapi); err != nil {
		return fmt.Errorf("global constants: %w", err)
	}
	if err := GenerateGlobalEnums(projectPath, eapi); err != nil {
		return fmt.Errorf("global enums: %w", err)
	}
	return nil
}

func GenerateClassConstants(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.constants.gen.go").
		Funcs(template.FuncMap{
			"goVariantConstructor": goVariantConstructor,
			"goMethodName":         goMethodName,
			"goArgumentName":       goArgumentName,
			"goArgumentType":       goArgumentType,
			"goReturnType":         goReturnType,
			"goClassEnumName":      goClassEnumName,
			"goClassConstantName":  goClassConstantName,
			"goClassStructName":    goClassStructName,
			"goClassInterfaceName": goClassInterfaceName,
			"coalesce":             coalesce,
		}).
		Parse(classesConstantsText)
	if err != nil {
		return fmt.Errorf("parse template classes.constants.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.constants.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "constant", fmt.Sprintf("classes.constants.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateClassEnums(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.enums.gen.go").
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
		Parse(classesEnumsText)
	if err != nil {
		return fmt.Errorf("parse template classes.enums.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.enums.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "constant", fmt.Sprintf("classes.enums.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateGlobalConstants(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	if len(extensionApi.GlobalConstants) == 0 {
		return nil
	}

	tmpl, err := template.New("globalconstants.gen.go").
		Parse(globalConstantsText)
	if err != nil {
		return fmt.Errorf("parse template globalconstants.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template globalconstants.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "constant", fmt.Sprintf("globalconstants.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateGlobalEnums(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	if len(extensionApi.GlobalEnums) == 0 {
		return nil
	}

	tmpl, err := template.New("globalenums.gen.go").
		Parse(globalEnumsText)
	if err != nil {
		return fmt.Errorf("parse template globalenums.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template globalenums.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "constant", fmt.Sprintf("globalenums.gen.go"))

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
