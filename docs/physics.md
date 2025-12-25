# Physics Patterns

This doc shows common 2D physics workflows in Go. Method names mirror the Godot API but use Go naming conventions.

## CharacterBody2D movement

```go
func (p *Player) V_PhysicsProcess(delta float64) {
	velocity := p.GetVelocity()
	velocity.X = p.speed
	p.SetVelocity(velocity)
	p.MoveAndSlide()
}

func RegisterClassPlayer() {
	ClassDBRegisterClass(NewPlayerFromOwnerObject, nil, nil, func(t *Player) {
		ClassDBBindMethodVirtual(t, "V_PhysicsProcess", "_physics_process", []string{"delta"}, nil)
	})
}
```

Useful methods: `SetVelocity`, `GetVelocity`, `MoveAndSlide`, `IsOnFloor`, `GetSlideCollisionCount`.

## RigidBody2D impulses

```go
func (b *Ball) Kick(impulse Vector2) {
	b.ApplyCentralImpulse(impulse)
}

func (b *Ball) EnableCCD() {
	b.SetContinuousCollisionDetectionMode(RigidBody2DCCDModeContinuous)
}
```

Useful methods: `ApplyCentralImpulse`, `ApplyImpulse`, `SetLinearVelocity`, `SetContactMonitor`.

## Area2D triggers

```go
func (a *GoalArea) V_Ready() {
	signal := NewStringNameWithLatin1Chars("body_entered")
	defer signal.Destroy()
	method := NewStringNameWithLatin1Chars("_on_body_entered")
	defer method.Destroy()
	callable := NewCallableWithObjectStringName(a, method)
	defer callable.Destroy()

	a.Connect(signal, callable, 0)
}

func (a *GoalArea) OnBodyEntered(body Node2D) {
	// Handle overlap
}

func RegisterClassGoalArea() {
	ClassDBRegisterClass(NewGoalAreaFromOwnerObject, nil, nil, func(t *GoalArea) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethod(t, "OnBodyEntered", "_on_body_entered", []string{"body"}, nil)
	})
}
```

Tip: For `body_entered` to fire, the Area2D needs a collision shape and `monitoring` enabled in the scene.
