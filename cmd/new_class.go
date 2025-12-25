package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

type newClassOptions struct {
	parentClass string
	outDir      string
	packageName string
	withReady   bool
	force       bool
}

const newClassTemplate = `package {{.Package}}

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// {{.ClassName}} implements GDClass evidence.
var _ GDClass = (*{{.ClassName}})(nil)

type {{.ClassName}} struct {
	{{.ParentClass}}Impl
}

func (c *{{.ClassName}}) GetClassName() string {
	return "{{.ClassName}}"
}

func (c *{{.ClassName}}) GetParentClassName() string {
	return "{{.ParentClass}}"
}

{{- if .WithReady }}
func (c *{{.ClassName}}) V_Ready() {
	// TODO: initialize
}

{{- end }}
func New{{.ClassName}}FromOwnerObject(owner *GodotObject) GDClass {
	obj := &{{.ClassName}}{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClass{{.ClassName}}() {
	ClassDBRegisterClass(New{{.ClassName}}FromOwnerObject, nil, nil, func(t *{{.ClassName}}) {
		{{- if .WithReady }}
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		{{- end }}
	})
}

func UnregisterClass{{.ClassName}}() {
	ClassDBUnregisterClass[*{{.ClassName}}]()
}
`

func init() {
	opts := &newClassOptions{}
	cmd := &cobra.Command{
		Use:   "new-class <ClassName>",
		Short: "Generate a Go class skeleton",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			className := args[0]
			return runNewClass(className, opts)
		},
	}
	cmd.Flags().StringVar(&opts.parentClass, "parent", "Node", "Parent Godot class name (without Impl)")
	cmd.Flags().StringVar(&opts.outDir, "out", ".", "Output directory")
	cmd.Flags().StringVar(&opts.packageName, "package", "", "Package name (defaults to output directory name)")
	cmd.Flags().BoolVar(&opts.withReady, "ready", true, "Include V_Ready virtual method")
	cmd.Flags().BoolVar(&opts.force, "force", false, "Overwrite existing file")
	rootCmd.AddCommand(cmd)
}

type newClassTemplateData struct {
	ClassName   string
	ParentClass string
	Package     string
	WithReady   bool
}

func runNewClass(className string, opts *newClassOptions) error {
	if !isValidClassName(className) {
		return fmt.Errorf("invalid class name %q: must be a valid Go identifier starting with an uppercase letter", className)
	}
	if !isValidIdentifier(opts.parentClass) {
		return fmt.Errorf("invalid parent class name %q", opts.parentClass)
	}

	outDir := filepath.Clean(opts.outDir)
	if outDir == "." {
		outDir = "."
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", outDir, err)
	}

	pkgName := opts.packageName
	if pkgName == "" {
		pkgName = inferPackageName(outDir)
	}
	if pkgName == "" {
		return errors.New("package name could not be inferred; pass --package")
	}
	if !isValidIdentifier(pkgName) {
		return fmt.Errorf("invalid package name %q", pkgName)
	}

	fileName := fmt.Sprintf("%s.go", strcase.ToSnake(className))
	filePath := filepath.Join(outDir, fileName)
	if !opts.force {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("file already exists: %s (use --force to overwrite)", filePath)
		}
	}

	data := newClassTemplateData{
		ClassName:   className,
		ParentClass: opts.parentClass,
		Package:     pkgName,
		WithReady:   opts.withReady,
	}

	tmpl, err := template.New("new-class").Parse(newClassTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("format generated code: %w", err)
	}

	if err := os.WriteFile(filePath, formatted, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", filePath, err)
	}

	fmt.Printf("Generated %s\n", filePath)
	return nil
}

func isValidClassName(name string) bool {
	if !isValidIdentifier(name) {
		return false
	}
	runes := []rune(name)
	if len(runes) == 0 {
		return false
	}
	return unicode.IsUpper(runes[0])
}

func isValidIdentifier(name string) bool {
	if name == "" {
		return false
	}
	identifier := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	return identifier.MatchString(name)
}

func inferPackageName(outDir string) string {
	base := filepath.Base(outDir)
	if base == "." || base == string(filepath.Separator) {
		return "main"
	}
	return base
}
