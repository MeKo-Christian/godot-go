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

func registerScoreOverlayTypes() {
	RegisterClassScoreOverlay()
}

func unregisterScoreOverlayTypes() {
	UnregisterClassScoreOverlay()
}

//export ScoreOverlayInit
func ScoreOverlayInit(pGetProcAddress unsafe.Pointer, pLibrary unsafe.Pointer, rInitialization unsafe.Pointer) bool {
	util.SetThreadName("score-overlay")
	log.Debug("ScoreOverlayInit called")
	initObj := NewInitObject(
		(ffi.GDExtensionInterfaceGetProcAddress)(pGetProcAddress),
		(ffi.GDExtensionClassLibraryPtr)(pLibrary),
		(*ffi.GDExtensionInitialization)(rInitialization),
	)
	initObj.RegisterSceneInitializer(registerScoreOverlayTypes)
	initObj.RegisterSceneTerminator(unregisterScoreOverlayTypes)
	return initObj.Init()
}
