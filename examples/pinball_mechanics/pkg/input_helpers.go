package pkg

import (
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/ffi"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// getGDClassInstance looks up a custom GDExtension class instance by its Object.
// This is needed because ObjectCastTo doesn't work for custom GDExtension classes.
func getGDClassInstance[T GDClass](obj Object) T {
	var zero T
	if obj == nil {
		return zero
	}
	owner := obj.GetGodotObjectOwner()
	if owner == nil {
		return zero
	}
	id := CallFunc_GDExtensionInterfaceObjectGetInstanceId(
		(GDExtensionConstObjectPtr)(unsafe.Pointer(owner)),
	)
	inst, ok := Internal.GDClassInstances.Get(id)
	if !ok {
		return zero
	}
	result, ok := inst.(T)
	if !ok {
		return zero
	}
	return result
}

const (
	actionLeftFlipper  = "pinball_left_flipper"
	actionRightFlipper = "pinball_right_flipper"
	actionLaunch       = "pinball_launch"
)

func setupPinballActions() {
	inputMap := getInputMapSingleton()
	if inputMap == nil {
		printLine("PinballGame: InputMap singleton not available")
		return
	}

	ensureAction(inputMap, actionLeftFlipper, func() RefInputEvent {
		return newRefInputEventKey(KEY_LEFT)
	})
	ensureAction(inputMap, actionRightFlipper, func() RefInputEvent {
		return newRefInputEventKey(KEY_RIGHT)
	})
	ensureAction(inputMap, actionLaunch, func() RefInputEvent {
		return newRefInputEventKey(KEY_SPACE)
	})
}

func ensureAction(inputMap InputMap, name string, eventFactory func() RefInputEvent) {
	action := NewStringNameWithLatin1Chars(name)
	defer action.Destroy()
	if inputMap.HasAction(action) {
		return
	}
	inputMap.AddAction(action, 0.5)
	event := eventFactory()
	if event != nil && event.IsValid() {
		inputMap.ActionAddEvent(action, event)
	}
}

func getInputMapSingleton() InputMap {
	owner := (*GodotObject)(unsafe.Pointer(GetSingleton("InputMap")))
	return NewInputMapWithGodotOwnerObject(owner)
}

func getClassDBSingleton() ClassDB {
	owner := (*GodotObject)(unsafe.Pointer(GetSingleton("ClassDB")))
	return NewClassDBWithGodotOwnerObject(owner)
}

func newInputEventKey(key Key) InputEventKey {
	classDB := getClassDBSingleton()
	if classDB == nil {
		return nil
	}
	className := NewStringNameWithLatin1Chars("InputEventKey")
	defer className.Destroy()
	v := classDB.Instantiate(className)
	// NOTE: Do NOT call v.Destroy() here!
	// InputEventKey is a RefCounted object. When the Variant is destroyed,
	// it decrements the refcount. Since we're extracting the object to keep it,
	// destroying the Variant would free the object (dropping refcount from 1 to 0).
	obj := v.ToObject()
	keyEvent, ok := ObjectCastTo(obj, "InputEventKey").(InputEventKey)
	if !ok || keyEvent == nil {
		printLine("PinballGame: failed to create InputEventKey")
		v.Destroy() // Safe to destroy here since we're not keeping the object
		return nil
	}
	keyEvent.SetKeycode(key)
	keyEvent.SetPressed(true)
	return keyEvent
}

func newRefInputEventKey(key Key) RefInputEvent {
	keyEvent := newInputEventKey(key)
	if keyEvent == nil {
		return nil
	}
	return NewRefInputEvent(keyEvent)
}

func nodePath(path string) NodePath {
	str := NewStringWithUtf8Chars(path)
	defer str.Destroy()
	return NewNodePathWithString(str)
}

func printLine(text string) {
	v := NewVariantGoString(text)
	defer v.Destroy()
	Print(v)
}

func inputSingleton() Input {
	return GetInputSingleton()
}
