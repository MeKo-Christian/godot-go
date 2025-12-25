# Signals

This guide covers defining, emitting, and connecting signals in godot-go.

## Define a signal

Signal definitions happen during class registration.

```go
func RegisterClassExample() {
	ClassDBRegisterClass(NewExampleFromOwnerObject, nil, nil, func(t *Example) {
		ClassDBAddSignal(t, "custom_signal",
			SignalParam{Type: GDEXTENSION_VARIANT_TYPE_STRING, Name: "name"},
			SignalParam{Type: GDEXTENSION_VARIANT_TYPE_INT, Name: "value"},
		)

		ClassDBBindMethod(t, "EmitCustomSignal", "emit_custom_signal", []string{"name", "value"}, nil)
	})
}
```

## Emit a signal

```go
func (e *Example) EmitCustomSignal(name string, value int64) {
	signal := NewStringNameWithLatin1Chars("custom_signal")
	defer signal.Destroy()
	snName := NewStringWithLatin1Chars(name)
	defer snName.Destroy()

	arg0 := NewVariantString(snName)
	defer arg0.Destroy()
	arg1 := NewVariantInt64(value)
	defer arg1.Destroy()

	e.EmitSignal(signal, arg0, arg1)
}
```

## Connect a signal

Signals connect via `Object.Connect` with a `Callable`.

```go
signal := NewStringNameWithLatin1Chars("custom_signal")
defer signal.Destroy()
method := NewStringNameWithLatin1Chars("_on_custom_signal")
defer method.Destroy()
callable := NewCallableWithObjectStringName(e, method)
defer callable.Destroy()

err := e.Connect(signal, callable, 0)
if err != OK {
	// handle connection error
}
```

## Disconnect a signal

```go
signal := NewStringNameWithLatin1Chars("custom_signal")
defer signal.Destroy()
method := NewStringNameWithLatin1Chars("_on_custom_signal")
defer method.Destroy()
callable := NewCallableWithObjectStringName(e, method)
defer callable.Destroy()

e.Disconnect(signal, callable)
```

## Tips

- Prefer storing `StringName` values if you connect or emit frequently.
- Always `Destroy()` `StringName`, `Callable`, and `Variant` values you create.
- Use `IsConnected` to avoid duplicate connections.
