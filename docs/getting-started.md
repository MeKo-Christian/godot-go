# Getting Started

This guide focuses on building and running a minimal GDExtension with godot-go.

## Requirements

- Godot 4.5+ (matching headers)
- Go 1.25.x
- gcc/clang toolchain
- `just`
- `goimports` (via `just installdeps`)

## Build the demo extension

The repository ships with a test demo project in `test/demo/` that doubles as a working template.

```bash
just installdeps

# Optional if you updated your Godot version
GODOT=/path/to/godot just update_godot_headers_from_binary

just generate
just build
```

The build outputs a platform-specific shared library into `test/demo/lib/` and the demo is wired up via `test/demo/example.gdextension`.

## Run the demo in headless mode

```bash
GODOT=/path/to/godot just ci_gen_test_project_files
GODOT=/path/to/godot just test
```

## Creating your own extension

The shortest path is to copy the demo layout and adjust names.

1. Copy `test/demo/example.gdextension` into your Godot project, rename it, and update `entry_symbol` and library filenames.
2. Create a Go package that exports the entry symbol and registers your classes.
3. Build a shared library (`-buildmode=c-shared`) and place it where the `.gdextension` expects.

### Entry point skeleton

```go
package mygame

/*
#cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/path/to/godot_headers -I${SRCDIR}/path/to/pkg/log -I${SRCDIR}/path/to/pkg/gdextension
*/
import "C"

import (
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/godot-go/godot-go/pkg/util"
)

func RegisterMyTypes() {
	RegisterClassMyNode()
}

func UnregisterMyTypes() {
	UnregisterClassMyNode()
}

//export MyGameInit
func MyGameInit(pGetProcAddress unsafe.Pointer, pLibrary unsafe.Pointer, rInitialization unsafe.Pointer) bool {
	util.SetThreadName("my-game")
	log.Debug("MyGameInit called")
	initObj := NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(pGetProcAddress),
		(ffi.GDExtensionClassLibraryPtr)(pLibrary),
		(*ffi.GDExtensionInitialization)(rInitialization),
	)
	initObj.RegisterSceneInitializer(RegisterMyTypes)
	initObj.RegisterSceneTerminator(UnregisterMyTypes)
	return initObj.Init()
}
```

Adjust the `#cgo` include paths to point at your local godot-go checkout (see `test/pkg/lib.go` for a working example inside this repo).

### Minimal class skeleton

```go
type MyNode struct {
	Node2DImpl
}

func (n *MyNode) GetClassName() string {
	return "MyNode"
}

func (n *MyNode) GetParentClassName() string {
	return "Node2D"
}

func (n *MyNode) V_Ready() {
	// Called when the node enters the scene tree.
}

func NewMyNodeFromOwnerObject(owner *GodotObject) GDClass {
	obj := &MyNode{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassMyNode() {
	ClassDBRegisterClass(NewMyNodeFromOwnerObject, nil, nil, func(t *MyNode) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
	})
}

func UnregisterClassMyNode() {
	ClassDBUnregisterClass[*MyNode]()
}
```

## Memory cleanup reminders

Many builtin types (StringName, Variant, Array, Callable) require `Destroy()` to avoid leaks. See `docs/memory.md` for the full list and recommended patterns.
