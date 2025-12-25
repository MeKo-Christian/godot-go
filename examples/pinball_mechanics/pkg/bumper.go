package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/ffi"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// Bumper implements GDClass evidence.
var _ GDClass = (*Bumper)(nil)

type Bumper struct {
	Area2DImpl
	scoreValue int64
	strength   float32
}

func (b *Bumper) GetClassName() string {
	return "Bumper"
}

func (b *Bumper) GetParentClassName() string {
	return "Area2D"
}

func (b *Bumper) V_Ready() {
	b.scoreValue = 100
	b.strength = 180

	signal := NewStringNameWithLatin1Chars("body_entered")
	defer signal.Destroy()
	method := NewStringNameWithLatin1Chars("_on_body_entered")
	defer method.Destroy()
	callable := NewCallableWithObjectStringName(b, method)
	defer callable.Destroy()
	b.Connect(signal, callable, 0)
}

func (b *Bumper) OnBodyEntered(body Node2D) {
	if body == nil {
		return
	}
	rigidBody, ok := ObjectCastTo(body, "RigidBody2D").(RigidBody2D)
	if !ok || rigidBody == nil {
		return
	}
	from := b.GetGlobalPosition()
	to := rigidBody.GetGlobalPosition()
	dir := from.DirectionTo(to)
	impulse := NewVector2WithFloat32Float32(dir.MemberGetx()*b.strength, dir.MemberGety()*b.strength)
	rigidBody.ApplyCentralImpulse(impulse)
	b.emitBumped()
}

func (b *Bumper) emitBumped() {
	name := NewStringNameWithLatin1Chars("bumped")
	defer name.Destroy()
	b.EmitSignal(name, NewVariantInt64(b.scoreValue))
}

func NewBumperFromOwnerObject(owner *GodotObject) GDClass {
	obj := &Bumper{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassBumper() {
	ClassDBRegisterClass(NewBumperFromOwnerObject, nil, nil, func(t *Bumper) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethod(t, "OnBodyEntered", "_on_body_entered", []string{"body"}, nil)
		ClassDBAddSignal(t, "bumped",
			SignalParam{Type: GDEXTENSION_VARIANT_TYPE_INT, Name: "points"},
		)
	})
}

func UnregisterClassBumper() {
	ClassDBUnregisterClass[*Bumper]()
}
