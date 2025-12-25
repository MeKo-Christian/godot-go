// Package gdextensionwrapper generates C code to wrap all of the gdextension
// methods to call functions on the gdextension_api_structs to work
// around the cgo C function pointer limitation.
package ffi

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/godot-go/godot-go/cmd/gdextensionparser/clang"
	"github.com/iancoleman/strcase"
)

var (
	//go:embed ffi_wrapper.h.tmpl
	ffiWrapperHeaderFileText string

	//go:embed ffi_wrapper.c.tmpl
	ffiWrapperSrcFileText string

	//go:embed ffi_wrapper.go.tmpl
	ffiWrapperGoFileText string

	//go:embed ffi.go.tmpl
	ffiFileText string
)

func Generate(projectPath string, ast clang.CHeaderFileAST) error {
	if err := GenerateGDExtensionWrapperHeaderFile(projectPath, ast); err != nil {
		return fmt.Errorf("ffi wrapper header: %w", err)
	}
	if err := GenerateGDExtensionWrapperSrcFile(projectPath, ast); err != nil {
		return fmt.Errorf("ffi wrapper source: %w", err)
	}
	if err := GenerateGDExtensionWrapperGoFile(projectPath, ast); err != nil {
		return fmt.Errorf("ffi wrapper go: %w", err)
	}
	if err := GenerateGDExtensionInterfaceGoFile(projectPath, ast); err != nil {
		return fmt.Errorf("ffi interface go: %w", err)
	}
	return nil
}

func GenerateGDExtensionWrapperHeaderFile(projectPath string, ast clang.CHeaderFileAST) error {
	tmpl, err := template.New("ffi_wrapper.gen.h").
		Funcs(template.FuncMap{
			"snakeCase": strcase.ToSnake,
		}).
		Parse(ffiWrapperHeaderFileText)
	if err != nil {
		return fmt.Errorf("parse template ffi_wrapper.gen.h: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, ast); err != nil {
		return fmt.Errorf("execute template ffi_wrapper.gen.h: %w", err)
	}

	filename := filepath.Join(projectPath, "pkg", "ffi", "ffi_wrapper.gen.h")
	return writeGeneratedFile(filename, b.Bytes())
}

func GenerateGDExtensionWrapperSrcFile(projectPath string, ast clang.CHeaderFileAST) error {
	tmpl, err := template.New("ffi_wrapper.gen.c").
		Funcs(template.FuncMap{
			"snakeCase": strcase.ToSnake,
		}).
		Parse(ffiWrapperSrcFileText)
	if err != nil {
		return fmt.Errorf("parse template ffi_wrapper.gen.c: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, ast); err != nil {
		return fmt.Errorf("execute template ffi_wrapper.gen.c: %w", err)
	}

	headerFileName := filepath.Join(projectPath, "pkg", "ffi", "ffi_wrapper.gen.c")
	return writeGeneratedFile(headerFileName, b.Bytes())
}

func GenerateGDExtensionWrapperGoFile(projectPath string, ast clang.CHeaderFileAST) error {
	funcs := template.FuncMap{
		"gdiVariableName":    gdiVariableName,
		"snakeCase":          strcase.ToSnake,
		"camelCase":          strcase.ToCamel,
		"goReturnType":       goReturnType,
		"goArgumentType":     goArgumentType,
		"goEnumValue":        goEnumValue,
		"add":                add,
		"cgoCastArgument":    cgoCastArgument,
		"cgoCastReturnType":  cgoCastReturnType,
		"cgoPinReturnType":   cgoPinReturnType,
		"cgoCleanUpArgument": cgoCleanUpArgument,
	}

	tmpl, err := template.New("ffi_wrapper.gen.go").
		Funcs(funcs).
		Parse(ffiWrapperGoFileText)
	if err != nil {
		return fmt.Errorf("parse template ffi_wrapper.gen.go: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, ast); err != nil {
		return fmt.Errorf("execute template ffi_wrapper.gen.go: %w", err)
	}

	headerFileName := filepath.Join(projectPath, "pkg", "ffi", "ffi_wrapper.gen.go")
	return writeGeneratedFile(headerFileName, b.Bytes())
}

func GenerateGDExtensionInterfaceGoFile(projectPath string, ast clang.CHeaderFileAST) error {
	funcs := template.FuncMap{
		"gdiVariableName":     gdiVariableName,
		"snakeCase":           strcase.ToSnake,
		"camelCase":           strcase.ToCamel,
		"goReturnType":        goReturnType,
		"goArgumentType":      goArgumentType,
		"goEnumValue":         goEnumValue,
		"add":                 add,
		"cgoCastArgument":     cgoCastArgument,
		"cgoCastReturnType":   cgoCastReturnType,
		"cgoCleanUpArgument":  cgoCleanUpArgument,
		"trimPrefix":          trimPrefix,
		"loadProcAddressName": loadProcAddressName,
	}

	tmpl, err := template.New("ffi.gen.go").
		Funcs(funcs).
		Parse(ffiFileText)
	if err != nil {
		return fmt.Errorf("parse template ffi.gen.go: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, ast); err != nil {
		return fmt.Errorf("execute template ffi.gen.go: %w", err)
	}

	headerFileName := filepath.Join(projectPath, "pkg", "ffi", "ffi.gen.go")
	return writeGeneratedFile(headerFileName, b.Bytes())
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
