package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// HelloWorld implements GDClass evidence.
var _ GDClass = (*HelloWorld)(nil)

type HelloWorld struct {
	Node2DImpl
}

func (h *HelloWorld) GetClassName() string {
	return "HelloWorld"
}

func (h *HelloWorld) GetParentClassName() string {
	return "Node2D"
}

func (h *HelloWorld) V_Ready() {
	msg := NewVariantGoString("Hello from godot-go")
	defer msg.Destroy()
	Print(msg)
}

func NewHelloWorldFromOwnerObject(owner *GodotObject) GDClass {
	obj := &HelloWorld{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassHelloWorld() {
	ClassDBRegisterClass(NewHelloWorldFromOwnerObject, nil, nil, func(t *HelloWorld) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
	})
}

func UnregisterClassHelloWorld() {
	ClassDBUnregisterClass[*HelloWorld]()
}
