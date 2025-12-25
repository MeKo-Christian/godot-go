package pkg

import (
	"fmt"
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// ActionInputDemo implements GDClass evidence.
var _ GDClass = (*ActionInputDemo)(nil)

type ActionInputDemo struct {
	NodeImpl
	score      int64
	actionJump StringName
	actionClick StringName
	initialized bool
}

func (a *ActionInputDemo) GetClassName() string {
	return "ActionInputDemo"
}

func (a *ActionInputDemo) GetParentClassName() string {
	return "Node"
}

func (a *ActionInputDemo) V_Ready() {
	a.SetProcess(true)
	a.setupActions()
	a.updateLabel()
	printLine("Actions demo ready: press Space or click.")
}

func (a *ActionInputDemo) V_Process(delta float64) {
	input := GetInputSingleton()
	if input.IsActionJustPressed(a.actionJump, true) {
		a.score += 10
		a.updateLabel()
	}
	if input.IsActionJustPressed(a.actionClick, true) {
		a.score += 1
		a.updateLabel()
	}
}

func (a *ActionInputDemo) V_ExitTree() {
	if a.initialized {
		a.actionJump.Destroy()
		a.actionClick.Destroy()
	}
}

func (a *ActionInputDemo) setupActions() {
	if a.initialized {
		return
	}
	a.actionJump = NewStringNameWithLatin1Chars("demo_jump")
	a.actionClick = NewStringNameWithLatin1Chars("demo_click")
	a.initialized = true

	inputMap := getInputMapSingleton()
	if !inputMap.HasAction(a.actionJump) {
		inputMap.AddAction(a.actionJump, 0.5)
		keyEvent := newInputEventKey(KEY_SPACE)
		if keyEvent != nil {
			inputMap.ActionAddEvent(a.actionJump, NewRefInputEvent(keyEvent))
		}
	}
	if !inputMap.HasAction(a.actionClick) {
		inputMap.AddAction(a.actionClick, 0.2)
		mouseEvent := newInputEventMouseButton(MOUSE_BUTTON_LEFT)
		if mouseEvent != nil {
			inputMap.ActionAddEvent(a.actionClick, NewRefInputEvent(mouseEvent))
		}
	}
}

func (a *ActionInputDemo) updateLabel() {
	label := a.GetNodeOrNull(nodePath("ScoreLabel"))
	if label == nil {
		printLine("ActionInputDemo: Label not found")
		return
	}
	labelNode, ok := ObjectCastTo(label, "Label").(Label)
	if !ok || labelNode == nil {
		printLine("ActionInputDemo: Label cast failed")
		return
	}
	text := NewStringWithUtf8Chars(fmt.Sprintf("Score: %d", a.score))
	defer text.Destroy()
	labelNode.SetText(text)
}

func NewActionInputDemoFromOwnerObject(owner *GodotObject) GDClass {
	obj := &ActionInputDemo{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassActionInputDemo() {
	ClassDBRegisterClass(NewActionInputDemoFromOwnerObject, nil, nil, func(t *ActionInputDemo) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Process", "_process", []string{"delta"}, nil)
		ClassDBBindMethodVirtual(t, "V_ExitTree", "_exit_tree", nil, nil)
	})
}

func UnregisterClassActionInputDemo() {
	ClassDBUnregisterClass[*ActionInputDemo]()
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
	className := NewStringNameWithLatin1Chars("InputEventKey")
	defer className.Destroy()
	v := classDB.Instantiate(className)
	defer v.Destroy()
	obj := v.ToObject()
	keyEvent, ok := ObjectCastTo(obj, "InputEventKey").(InputEventKey)
	if !ok || keyEvent == nil {
		printLine("failed to create InputEventKey")
		return nil
	}
	keyEvent.SetKeycode(key)
	keyEvent.SetPressed(true)
	return keyEvent
}

func newInputEventMouseButton(button MouseButton) InputEventMouseButton {
	classDB := getClassDBSingleton()
	className := NewStringNameWithLatin1Chars("InputEventMouseButton")
	defer className.Destroy()
	v := classDB.Instantiate(className)
	defer v.Destroy()
	obj := v.ToObject()
	mouseEvent, ok := ObjectCastTo(obj, "InputEventMouseButton").(InputEventMouseButton)
	if !ok || mouseEvent == nil {
		printLine("failed to create InputEventMouseButton")
		return nil
	}
	mouseEvent.SetButtonIndex(button)
	mouseEvent.SetPressed(true)
	return mouseEvent
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
