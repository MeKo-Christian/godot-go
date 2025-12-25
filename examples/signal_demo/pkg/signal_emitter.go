package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// SignalEmitter implements GDClass evidence.
var _ GDClass = (*SignalEmitter)(nil)

type SignalEmitter struct {
	NodeImpl
}

func (s *SignalEmitter) GetClassName() string {
	return "SignalEmitter"
}

func (s *SignalEmitter) GetParentClassName() string {
	return "Node"
}

func (s *SignalEmitter) V_Ready() {
	listener := s.getListener()
	if listener == nil {
		printLine("SignalEmitter: listener not found")
		return
	}

	signalName := NewStringNameWithLatin1Chars("demo_signal")
	defer signalName.Destroy()
	methodName := NewStringNameWithLatin1Chars("_on_demo_signal")
	defer methodName.Destroy()
	callable := NewCallableWithObjectStringName(listener, methodName)
	defer callable.Destroy()

	s.Connect(signalName, callable, 0)

	arg0 := NewVariantGoString("hello")
	defer arg0.Destroy()
	arg1 := NewVariantInt64(1)
	defer arg1.Destroy()
	s.EmitSignal(signalName, arg0, arg1)
}

func (s *SignalEmitter) getListener() Node {
	np := newNodePath("Listener")
	defer np.Destroy()
	return s.GetNodeOrNull(np)
}

func NewSignalEmitterFromOwnerObject(owner *GodotObject) GDClass {
	obj := &SignalEmitter{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassSignalEmitter() {
	ClassDBRegisterClass(NewSignalEmitterFromOwnerObject, nil, nil, func(t *SignalEmitter) {
		ClassDBAddSignal(t, "demo_signal",
			SignalParam{Type: GDEXTENSION_VARIANT_TYPE_STRING, Name: "message"},
			SignalParam{Type: GDEXTENSION_VARIANT_TYPE_INT, Name: "count"},
		)
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
	})
}

func UnregisterClassSignalEmitter() {
	ClassDBUnregisterClass[*SignalEmitter]()
}

func newNodePath(path string) NodePath {
	str := NewStringWithUtf8Chars(path)
	defer str.Destroy()
	return NewNodePathWithString(str)
}

func printLine(text string) {
	v := NewVariantGoString(text)
	defer v.Destroy()
	Print(v)
}
