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
func Generate(projectPath string, eapi extensionapiparser.ExtensionApi) {
	var err error
	if err = GenerateNativeStrucutres(projectPath, eapi); err != nil {
		panic(err)
	}
}

func GenerateNativeStrucutres(projectPath string, extensionApi extensionapiparser.ExtensionApi) error {
	tmpl, err := template.New("nativestructures.gen.go").
		Funcs(template.FuncMap{
			"nativeStructureFormatToFields": nativeStructureFormatToFields,
			"hasPrefix":                     strings.HasPrefix,
		}).
		Parse(nativeStructuresText)
	if err != nil {
		return err
	}

	var b bytes.Buffer

	err = tmpl.Execute(&b, extensionApi)
	if err != nil {
		return err
	}

	filename := filepath.Join(projectPath, "pkg", "nativestructure", fmt.Sprintf("nativestructures.gen.go"))

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(b.Bytes())
	if err != nil {
		return err
	}

	return nil
}
