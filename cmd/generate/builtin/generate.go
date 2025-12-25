package builtin

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/godot-go/godot-go/cmd/extensionapiparser"
	"github.com/godot-go/godot-go/cmd/gdextensionparser/clang"
	"github.com/iancoleman/strcase"
)

var (
	//go:embed builtinclasses.bindings.go.tmpl
	builtinClassesBindingsText string

	//go:embed builtinclasses.go.tmpl
	builtinClassesText string

	//go:embed variant.go.tmpl
	variantGoText string

	//go:embed classes.interfaces.go.tmpl
	classesInterfacesText string

	//go:embed classes.ref.interfaces.go.tmpl
	classesRefInterfacesText string
)

func Generate(projectPath string, ast clang.CHeaderFileAST, eapi extensionapiparser.ExtensionApi) error {
	if err := GenerateBuiltinClasses(projectPath, eapi); err != nil {
		return fmt.Errorf("builtin classes: %w", err)
	}
	if err := GenerateBuiltinClassBindings(projectPath, eapi); err != nil {
		return fmt.Errorf("builtin class bindings: %w", err)
	}
	if err := GenerateClassInterfaces(projectPath, eapi); err != nil {
		return fmt.Errorf("class interfaces: %w", err)
	}
	if err := GenerateClassRefInterfaces(projectPath, eapi); err != nil {
		return fmt.Errorf("class ref interfaces: %w", err)
	}
	if err := GenerateVariantGoFile(projectPath, ast); err != nil {
		return fmt.Errorf("variant: %w", err)
	}
	return nil
}

func GenerateBuiltinClasses(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("builtinclasses.gen.go").
		Funcs(template.FuncMap{
			"upper":                    strings.ToUpper,
			"upperFirstChar":           upperFirstChar,
			"snakeCase":                snakeCase,
			"goMethodName":             goMethodName,
			"goArgumentName":           goArgumentName,
			"goArgumentType":           goArgumentType,
			"goHasArgumentTypeEncoder": goHasArgumentTypeEncoder,
			"goReturnType":             goReturnType,
			"goDecodeNumberType":       goDecodeNumberType,
			"getOperatorIdName":        getOperatorIdName,
			"typeHasPtr":               typeHasPtr,
			"goEncoder":                goEncoder,
			"goEncodeIsReference":      goEncodeIsReference,
		}).
		Parse(builtinClassesText)
	if err != nil {
		return fmt.Errorf("parse template builtinclasses.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template builtinclasses.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "builtin", fmt.Sprintf("builtinclasses.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateBuiltinClassBindings(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("builtinclasses.bindings.gen.go").
		Funcs(template.FuncMap{
			"upper":             strings.ToUpper,
			"lowerFirstChar":    lowerFirstChar,
			"screamingSnake":    screamingSnake,
			"getOperatorIdName": getOperatorIdName,
			"goEncoder":         goEncoder,
		}).
		Parse(builtinClassesBindingsText)
	if err != nil {
		return fmt.Errorf("parse template builtinclasses.bindings.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template builtinclasses.bindings.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "builtin", fmt.Sprintf("builtinclasses.bindings.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateClassInterfaces(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.interfaces.gen.go").
		Funcs(template.FuncMap{
			"goMethodName":         goMethodName,
			"goArgumentName":       goArgumentName,
			"goArgumentType":       goArgumentType,
			"goReturnType":         goReturnType,
			"goClassInterfaceName": goClassInterfaceName,
			"coalesce":             coalesce,
		}).
		Parse(classesInterfacesText)
	if err != nil {
		return fmt.Errorf("parse template classes.interfaces.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.interfaces.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "builtin", fmt.Sprintf("classes.interfaces.gen.go"))

	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateClassRefInterfaces(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("classes.ref.instances.gen.go").
		Funcs(template.FuncMap{
			"goClassInterfaceName": goClassInterfaceName,
			"goEncoder":            goEncoder,
		}).
		Parse(classesRefInterfacesText)
	if err != nil {
		return fmt.Errorf("parse template classes.ref.instances.gen.go: %w", err)
	}

	var b bytes.Buffer

	if err := tmpl.Execute(&b, extensionApi); err != nil {
		return fmt.Errorf("execute template classes.ref.instances.gen.go: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "builtin", fmt.Sprintf("classes.ref.interfaces.gen.go"))

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

func GenerateVariantGoFile(projectPath string, ast clang.CHeaderFileAST) error {
	funcs := template.FuncMap{
		"snakeCase":          strcase.ToSnake,
		"camelCase":          strcase.ToCamel,
		"goEncoder":          goEncoder,
		"astVariantMetadata": astVariantMetadata,
	}

	tmpl, err := template.New("variant.gen.go").
		Funcs(funcs).
		Parse(variantGoText)
	if err != nil {
		return fmt.Errorf("parse template variant.gen.go: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, ast); err != nil {
		return fmt.Errorf("execute template variant.gen.go: %w", err)
	}

	goFileName := filepath.Join(projectPath, "pkg", "builtin", "variant.gen.go")
	return writeGeneratedFile(goFileName, b.Bytes())
}
