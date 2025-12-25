package pkg

/*
#cgo CFLAGS: -I${SRCDIR} -I${SRCDIR}/../../../godot_headers -I${SRCDIR}/../../../pkg/log -I${SRCDIR}/../../../pkg/gdextension
*/
import "C"

import (
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/core"
	"github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/log"
	"github.com/godot-go/godot-go/pkg/util"
)

func registerInputDemoTypes() {
	RegisterClassInputDemo()
}

func unregisterInputDemoTypes() {
	UnregisterClassInputDemo()
}

//export InputDemoInit
func InputDemoInit(pGetProcAddress unsafe.Pointer, pLibrary unsafe.Pointer, rInitialization unsafe.Pointer) bool {
	util.SetThreadName("input-demo")
	log.Debug("InputDemoInit called")
	initObj := NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(pGetProcAddress),
		(ffi.GDExtensionClassLibraryPtr)(pLibrary),
		(*ffi.GDExtensionInitialization)(rInitialization),
	)
	initObj.RegisterSceneInitializer(registerInputDemoTypes)
	initObj.RegisterSceneTerminator(unregisterInputDemoTypes)
	return initObj.Init()
}
