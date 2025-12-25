package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// BouncingBall implements GDClass evidence.
var _ GDClass = (*BouncingBall)(nil)

type BouncingBall struct {
	RigidBody2DImpl
}

func (b *BouncingBall) GetClassName() string {
	return "BouncingBall"
}

func (b *BouncingBall) GetParentClassName() string {
	return "RigidBody2D"
}

func (b *BouncingBall) V_Ready() {
	b.SetContinuousCollisionDetectionMode(RigidBody2DCCDModeContinuous)
	b.ApplyCentralImpulse(NewVector2WithFloat32Float32(180, -240))
}

func NewBouncingBallFromOwnerObject(owner *GodotObject) GDClass {
	obj := &BouncingBall{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassBouncingBall() {
	ClassDBRegisterClass(NewBouncingBallFromOwnerObject, nil, nil, func(t *BouncingBall) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
	})
}

func UnregisterClassBouncingBall() {
	ClassDBUnregisterClass[*BouncingBall]()
}
