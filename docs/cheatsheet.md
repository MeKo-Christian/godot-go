# GDScript to Go Cheat Sheet

This is a quick translation guide for common GDScript patterns.

## Class definition

GDScript:

```gdscript
extends CharacterBody2D
```

Go (composition over inheritance):

```go
type Player struct {
	CharacterBody2DImpl
}

func (p *Player) GetClassName() string   { return "Player" }
func (p *Player) GetParentClassName() string { return "CharacterBody2D" }
```

## Virtual methods

GDScript:

```gdscript
func _ready():
	pass
```

Go:

```go
func (p *Player) V_Ready() {
	// Ready logic
}

ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
```

## Properties

GDScript:

```gdscript
@export var speed := 400.0
```

Go:

```go
func (p *Player) GetSpeed() float32 { return p.speed }
func (p *Player) SetSpeed(v float32) { p.speed = v }

ClassDBBindMethod(t, "GetSpeed", "get_speed", nil, nil)
ClassDBBindMethod(t, "SetSpeed", "set_speed", []string{"value"}, nil)
ClassDBAddProperty(t, GDEXTENSION_VARIANT_TYPE_FLOAT, "speed", "set_speed", "get_speed")
```

## Signals

GDScript:

```gdscript
signal hit(damage: int)
```

Go:

```go
ClassDBAddSignal(t, "hit",
	SignalParam{Type: GDEXTENSION_VARIANT_TYPE_INT, Name: "damage"},
)
```

## Emit a signal

```go
signalName := NewStringNameWithLatin1Chars("hit")
defer signalName.Destroy()
arg := NewVariantInt64(10)
defer arg.Destroy()

p.EmitSignal(signalName, arg)
```

## Connect to a signal

```go
signalName := NewStringNameWithLatin1Chars("hit")
defer signalName.Destroy()
method := NewStringNameWithLatin1Chars("_on_hit")
defer method.Destroy()
callable := NewCallableWithObjectStringName(p, method)
defer callable.Destroy()

p.Connect(signalName, callable, 0)
```

## Arrays and dictionaries

```go
arr := NewArray()
defer arr.Destroy()
arr.Append(NewVariantInt64(1))
arr.Append(NewVariantInt64(2))

dict := NewDictionary()
defer dict.Destroy()
dict.SetKeyed("score", NewVariantInt64(100))
```

## Variant conversion

```go
v := NewVariantFloat32(1.5)
defer v.Destroy()

f := v.ToFloat32()
```

## Memory management

Builtin types like `String`, `StringName`, `Variant`, `Array`, `Dictionary`, `Callable`, and all `Packed*Array` types require manual `Destroy()`. See `docs/memory.md`.
