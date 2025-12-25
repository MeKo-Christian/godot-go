package pkg

import (
	"fmt"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// SignalListener implements GDClass evidence.
var _ GDClass = (*SignalListener)(nil)

type SignalListener struct {
	NodeImpl
}

func (s *SignalListener) GetClassName() string {
	return "SignalListener"
}

func (s *SignalListener) GetParentClassName() string {
	return "Node"
}

func (s *SignalListener) OnDemoSignal(message string, count int64) {
	printLine(fmt.Sprintf("SignalListener received: %s (%d)", message, count))
}

func NewSignalListenerFromOwnerObject(owner *GodotObject) GDClass {
	obj := &SignalListener{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassSignalListener() {
	ClassDBRegisterClass(NewSignalListenerFromOwnerObject, nil, nil, func(t *SignalListener) {
		ClassDBBindMethod(t, "OnDemoSignal", "_on_demo_signal", []string{"message", "count"}, nil)
	})
}

func UnregisterClassSignalListener() {
	ClassDBUnregisterClass[*SignalListener]()
}

func printLine(text string) {
	v := NewVariantGoString(text)
	defer v.Destroy()
	Print(v)
}
