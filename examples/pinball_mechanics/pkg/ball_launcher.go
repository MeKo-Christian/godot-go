package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// BallLauncher implements GDClass evidence.
var _ GDClass = (*BallLauncher)(nil)

type BallLauncher struct {
	Node2DImpl
	actionName  StringName
	initialized bool
	ball        RigidBody2D
	startPos    Vector2
}

func (b *BallLauncher) GetClassName() string {
	return "BallLauncher"
}

func (b *BallLauncher) GetParentClassName() string {
	return "Node2D"
}

func (b *BallLauncher) V_Ready() {
	b.SetPhysicsProcess(true)
	b.actionName = NewStringNameWithLatin1Chars(actionLaunch)
	b.initialized = true
	b.cacheBall()
}

func (b *BallLauncher) V_PhysicsProcess(_delta float64) {
	if !b.initialized || b.ball == nil {
		return
	}
	input := inputSingleton()
	if input == nil {
		return
	}
	if input.IsActionJustPressed(b.actionName, true) {
		b.launchBall()
	}
}

func (b *BallLauncher) V_ExitTree() {
	if b.initialized {
		b.actionName.Destroy()
	}
}

func (b *BallLauncher) OnBallDrained(body Node2D) {
	if b.ball == nil || body == nil {
		return
	}
	if body.GetInstanceId() != b.ball.GetInstanceId() {
		return
	}
	b.resetBall()
}

func (b *BallLauncher) cacheBall() {
	ballNode := b.GetNodeOrNull(nodePath("../Ball"))
	if ballNode == nil {
		printLine("BallLauncher: Ball not found")
		return
	}
	b.ball, _ = ObjectCastTo(ballNode, "RigidBody2D").(RigidBody2D)
	if b.ball == nil {
		printLine("BallLauncher: Ball cast failed")
		return
	}
	b.startPos = b.ball.GetGlobalPosition()
}

func (b *BallLauncher) launchBall() {
	impulse := NewVector2WithFloat32Float32(0, -220)
	b.ball.ApplyCentralImpulse(impulse)
}

func (b *BallLauncher) resetBall() {
	b.ball.SetLinearVelocity(NewVector2WithFloat32Float32(0, 0))
	b.ball.SetAngularVelocity(0)
	b.ball.SetGlobalPosition(b.startPos)
}

func NewBallLauncherFromOwnerObject(owner *GodotObject) GDClass {
	obj := &BallLauncher{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassBallLauncher() {
	ClassDBRegisterClass(NewBallLauncherFromOwnerObject, nil, nil, func(t *BallLauncher) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_PhysicsProcess", "_physics_process", []string{"delta"}, nil)
		ClassDBBindMethodVirtual(t, "V_ExitTree", "_exit_tree", nil, nil)
		ClassDBBindMethod(t, "OnBallDrained", "_on_ball_drained", []string{"body"}, nil)
	})
}

func UnregisterClassBallLauncher() {
	ClassDBUnregisterClass[*BallLauncher]()
}
