package core

// #include <godot/gdextension_interface.h>
// #include "method_bind.h"
// #include <stdio.h>
// #include <stdlib.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/ffi"
	"github.com/godot-go/godot-go/pkg/log"
	"go.uber.org/zap"
)

// GoCallback_MethodBindMethodCall is called when GDScript vararg methods calls into Go.
//
//export GoCallback_MethodBindMethodCall
func GoCallback_MethodBindMethodCall(
	methodUserData unsafe.Pointer,
	instPtr C.GDExtensionClassInstancePtr,
	argPtrs *C.GDExtensionVariantPtr,
	argumentCount C.GDExtensionInt,
	rReturn C.GDExtensionVariantPtr,
	rError *C.GDExtensionCallError,
) {
	// Add panic recovery to catch any panics and report them via rError
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic in Go method callback",
				zap.Any("panic", r),
				zap.Stack("stack"),
			)
			// Set rError to indicate method call failed
			callErr := (*GDExtensionCallError)(unsafe.Pointer(rError))
			callErr.SetErrorFields(GDEXTENSION_CALL_ERROR_INVALID_METHOD, 0, 0)
		}
	}()

	ud := (cgo.Handle)(methodUserData)
	bind, ok := ud.Value().(*GoMethodMetadata)
	if !ok || bind == nil {
		log.Error("unable to retrieve methodUserData")
		callErr := (*GDExtensionCallError)(unsafe.Pointer(rError))
		callErr.SetErrorFields(GDEXTENSION_CALL_ERROR_INVALID_METHOD, 0, 0)
		return
	}
	pnr.Pin(instPtr)
	inst := ObjectClassFromGDExtensionClassInstancePtr((GDExtensionClassInstancePtr)(instPtr))
	if inst == nil {
		log.Error("GDExtensionClassInstancePtr cannot be null")
		callErr := (*GDExtensionCallError)(unsafe.Pointer(rError))
		callErr.SetErrorFields(GDEXTENSION_CALL_ERROR_INSTANCE_IS_NULL, 0, 0)
		return
	}
	pnr.Pin(inst)
	cn := inst.GetClass()
	defer cn.Destroy()
	log.Debug("GoCallback_MethodBindMethodCall called",
		zap.String("class", cn.ToUtf8()),
		zap.String("method", bind.GdMethodName),
		zap.String("bind", bind.String()),
	)
	argPtrSlice := unsafe.Slice((*GDExtensionConstVariantPtr)(argPtrs), int(argumentCount))
	args := make([]Variant, argumentCount)
	for i := range argPtrSlice {
		pnr.Pin(argPtrSlice[i])
		args[i] = NewVariantCopyWithGDExtensionConstVariantPtr(argPtrSlice[i])
	}

	// Call the Go method and check for errors
	retCall, callErr := bind.Call(inst, args...)
	if callErr != nil {
		log.Error("method call failed",
			zap.String("class", cn.ToUtf8()),
			zap.String("method", bind.GdMethodName),
			zap.Error(callErr),
		)
		// Copy error fields to rError
		*(*GDExtensionCallError)(unsafe.Pointer(rError)) = *callErr
		return
	}

	*(*Variant)(unsafe.Pointer(rReturn)) = retCall
	pnr.Pin(rReturn)
}

// called when godot calls into golang code
//
//export GoCallback_MethodBindMethodPtrcall
func GoCallback_MethodBindMethodPtrcall(
	methodUserData unsafe.Pointer,
	instPtr C.GDExtensionClassInstancePtr,
	argPtrs *C.GDExtensionConstTypePtr,
	rReturn C.GDExtensionTypePtr,
) {
	ud := (cgo.Handle)(methodUserData)
	bind, ok := ud.Value().(*GoMethodMetadata)
	if !ok || bind == nil {
		log.Panic("unable to retrieve methodUserData")
	}
	inst := ObjectClassFromGDExtensionClassInstancePtr((GDExtensionClassInstancePtr)(instPtr))
	if inst == nil {
		log.Panic("GDExtensionClassInstancePtr canoot be null")
	}
	cn := inst.GetClass()
	defer cn.Destroy()
	log.Debug("GoCallback_MethodBindMethodPtrcall called",
		zap.String("class", cn.ToUtf8()),
		zap.String("method", bind.String()),
	)
	sliceLen := len(bind.GoArgumentTypes)
	argsSlice := unsafe.Slice((*GDExtensionConstTypePtr)(unsafe.Pointer(argPtrs)), sliceLen)
	bind.Ptrcall(
		inst,
		argsSlice,
		(GDExtensionUninitializedTypePtr)(rReturn),
	)
	pnr.Pin(rReturn)
}
