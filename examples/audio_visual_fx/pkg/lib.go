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

func registerAudioVisualFXTypes() {
	RegisterClassAudioVisualDemo()
}

func unregisterAudioVisualFXTypes() {
	UnregisterClassAudioVisualDemo()
}

//export AudioVisualFXInit
func AudioVisualFXInit(pGetProcAddress unsafe.Pointer, pLibrary unsafe.Pointer, rInitialization unsafe.Pointer) bool {
	util.SetThreadName("audio-visual-fx")
	log.Debug("AudioVisualFXInit called")
	initObj := NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(pGetProcAddress),
		(ffi.GDExtensionClassLibraryPtr)(pLibrary),
		(*ffi.GDExtensionInitialization)(rInitialization),
	)
	initObj.RegisterSceneInitializer(registerAudioVisualFXTypes)
	initObj.RegisterSceneTerminator(unregisterAudioVisualFXTypes)
	return initObj.Init()
}
