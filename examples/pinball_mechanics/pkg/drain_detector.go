package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/ffi"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// DrainDetector implements GDClass evidence.
var _ GDClass = (*DrainDetector)(nil)

type DrainDetector struct {
	Area2DImpl
}

func (d *DrainDetector) GetClassName() string {
	return "DrainDetector"
}

func (d *DrainDetector) GetParentClassName() string {
	return "Area2D"
}

func (d *DrainDetector) V_Ready() {
	signal := NewStringNameWithLatin1Chars("body_entered")
	defer signal.Destroy()
	method := NewStringNameWithLatin1Chars("_on_body_entered")
	defer method.Destroy()
	callable := NewCallableWithObjectStringName(d, method)
	defer callable.Destroy()
	d.Connect(signal, callable, 0)
}

func (d *DrainDetector) OnBodyEntered(body Node2D) {
	if body == nil {
		return
	}
	name := NewStringNameWithLatin1Chars("drained")
	defer name.Destroy()
	d.EmitSignal(name, NewVariantGodotObject(body.GetGodotObjectOwner()))
}

func NewDrainDetectorFromOwnerObject(owner *GodotObject) GDClass {
	obj := &DrainDetector{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassDrainDetector() {
	ClassDBRegisterClass(NewDrainDetectorFromOwnerObject, nil, nil, func(t *DrainDetector) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethod(t, "OnBodyEntered", "_on_body_entered", []string{"body"}, nil)
		ClassDBAddSignal(t, "drained",
			SignalParam{Type: GDEXTENSION_VARIANT_TYPE_OBJECT, Name: "body"},
		)
	})
}

func UnregisterClassDrainDetector() {
	ClassDBUnregisterClass[*DrainDetector]()
}
