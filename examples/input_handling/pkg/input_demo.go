package pkg

import (
	"fmt"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// InputDemo implements GDClass evidence.
var _ GDClass = (*InputDemo)(nil)

type InputDemo struct {
	NodeImpl
}

func (d *InputDemo) GetClassName() string {
	return "InputDemo"
}

func (d *InputDemo) GetParentClassName() string {
	return "Node"
}

func (d *InputDemo) V_Ready() {
	d.SetProcessInput(true)
	printLine("Input demo ready: press keys or click the mouse.")
}

func (d *InputDemo) V_Input(refEvent RefInputEvent) {
	event := refEvent.TypedPtr()
	if event == nil {
		printLine("InputDemo.V_Input: null event")
		return
	}

	if keyEvent, ok := ObjectCastTo(event, "InputEventKey").(InputEventKey); ok {
		text := keyEvent.AsTextKeyLabel()
		defer text.Destroy()
		printLine(fmt.Sprintf("Key: %s pressed=%v", text.ToUtf8(), keyEvent.IsPressed()))
		return
	}

	if mouseButton, ok := ObjectCastTo(event, "InputEventMouseButton").(InputEventMouseButton); ok {
		printLine(fmt.Sprintf("Mouse button=%d pressed=%v", mouseButton.GetButtonIndex(), mouseButton.IsPressed()))
		return
	}

	if mouseMotion, ok := ObjectCastTo(event, "InputEventMouseMotion").(InputEventMouseMotion); ok {
		relative := mouseMotion.GetRelative()
		printLine(fmt.Sprintf("Mouse motion dx=%.1f dy=%.1f", relative.MemberGetx(), relative.MemberGety()))
		return
	}
}

func NewInputDemoFromOwnerObject(owner *GodotObject) GDClass {
	obj := &InputDemo{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassInputDemo() {
	ClassDBRegisterClass(NewInputDemoFromOwnerObject, nil, nil, func(t *InputDemo) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Input", "_input", []string{"event"}, nil)
	})
}

func UnregisterClassInputDemo() {
	ClassDBUnregisterClass[*InputDemo]()
}

func printLine(text string) {
	v := NewVariantGoString(text)
	defer v.Destroy()
	Print(v)
}
