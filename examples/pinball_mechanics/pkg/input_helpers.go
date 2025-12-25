package pkg

import (
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

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
	defer v.Destroy()
	obj := v.ToObject()
	keyEvent, ok := ObjectCastTo(obj, "InputEventKey").(InputEventKey)
	if !ok || keyEvent == nil {
		printLine("PinballGame: failed to create InputEventKey")
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
