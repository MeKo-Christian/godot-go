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

func registerSignalDemoTypes() {
	RegisterClassSignalEmitter()
	RegisterClassSignalListener()
}

func unregisterSignalDemoTypes() {
	UnregisterClassSignalEmitter()
	UnregisterClassSignalListener()
}

//export SignalDemoInit
func SignalDemoInit(pGetProcAddress unsafe.Pointer, pLibrary unsafe.Pointer, rInitialization unsafe.Pointer) bool {
	util.SetThreadName("signal-demo")
	log.Debug("SignalDemoInit called")
	initObj := NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(pGetProcAddress),
		(ffi.GDExtensionClassLibraryPtr)(pLibrary),
		(*ffi.GDExtensionInitialization)(rInitialization),
	)
	initObj.RegisterSceneInitializer(registerSignalDemoTypes)
	initObj.RegisterSceneTerminator(unregisterSignalDemoTypes)
	return initObj.Init()
}
