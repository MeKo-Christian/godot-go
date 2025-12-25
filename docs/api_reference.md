# Common API Reference

This is a concise reference for common game-dev classes. The full generated interfaces live in `pkg/builtin/classes.interfaces.gen.go`.

## Node

Use for scene graph operations.

Common methods:

- `AddChild(node Node, force_readable_name bool, internalMode NodeInternalMode)`
- `RemoveChild(node Node)`
- `GetParent() Node`
- `GetNode(path NodePath) Node` / `GetNodeOrNull(path NodePath) Node`
- `IsInsideTree() bool`
- `SetProcess(enable bool)` / `SetPhysicsProcess(enable bool)`
- `AddToGroup(group StringName, persistent bool)`

## Node2D

Use for 2D transforms.

Common methods:

- `SetPosition(position Vector2)` / `GetPosition() Vector2`
- `SetRotation(radians float32)` / `GetRotation() float32`
- `SetScale(scale Vector2)` / `GetScale() Vector2`
- `Translate(offset Vector2)`
- `LookAt(point Vector2)`

## CharacterBody2D

Kinematic-style movement.

Common methods:

- `SetVelocity(velocity Vector2)` / `GetVelocity() Vector2`
- `MoveAndSlide() bool`
- `IsOnFloor() bool` / `IsOnWall() bool`
- `GetSlideCollisionCount() int32`
- `GetSlideCollision(slide_idx int32) RefKinematicCollision2D`

## RigidBody2D

Physics-driven movement.

Common methods:

- `ApplyCentralImpulse(impulse Vector2)`
- `ApplyImpulse(impulse Vector2, position Vector2)`
- `SetLinearVelocity(linear_velocity Vector2)`
- `SetAngularVelocity(angular_velocity float32)`
- `SetContinuousCollisionDetectionMode(mode RigidBody2DCCDMode)`
- `SetContactMonitor(enabled bool)`

## Area2D

Overlap detection and triggers.

Common patterns:

- Connect to signals: `body_entered`, `area_entered`
- Use `Connect` on the `Area2D` or `Node` base to wire callables

## Input

Polling and action checks.

Common methods:

- `IsActionPressed(action StringName, exact_match bool) bool`
- `IsActionJustPressed(action StringName, exact_match bool) bool`
- `IsActionJustReleased(action StringName, exact_match bool) bool`
- `GetActionStrength(action StringName, exact_match bool) float32`
- `GetVector(negative_x, positive_x, negative_y, positive_y StringName, deadzone float32) Vector2`

## Timer

Use for one-shot or repeating callbacks.

Common methods:

- `SetWaitTime(time_sec float64)`
- `SetOneShot(enable bool)`
- `Start(time_sec float64)`
- `Stop()`
- `IsStopped() bool`

## Control and Label

UI basics.

Control common methods:

- `SetAnchorsPreset(preset ControlLayoutPreset, keep_offsets bool)`
- `SetOffsetsPreset(preset ControlLayoutPreset, resize_mode ControlLayoutPresetMode, margin int32)`
- `SetSize(size Vector2, keep_offsets bool)` / `GetSize() Vector2`

Label common methods:

- `SetText(text String)`
- `SetHorizontalAlignment(alignment HorizontalAlignment)`
- `SetVerticalAlignment(alignment VerticalAlignment)`

## AnimationPlayer

Animation playback.

Common methods:

- `Play(name StringName, custom_blend float32, custom_speed float32, from_end bool)`
- `Stop()`
- `IsPlaying() bool`
- `Seek(seconds float64, update bool)`
