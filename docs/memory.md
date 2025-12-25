# Memory Ownership and Manual Cleanup

godot-go exposes several Godot types that allocate native resources. These
must be released by calling `Destroy()` when you are done with them.

## Types Requiring Manual Cleanup

| Type | Requires `Destroy()` | Location |
|------|---------------------|----------|
| String | ✅ | [builtinclasses.gen.go:210](../pkg/builtin/builtinclasses.gen.go#L210) |
| StringName | ✅ | [builtinclasses.gen.go:16358](../pkg/builtin/builtinclasses.gen.go#L16358) |
| Variant | ✅ | [variant.go:394](../pkg/builtin/variant.go#L394) |
| Array | ✅ | [builtinclasses.gen.go:21751](../pkg/builtin/builtinclasses.gen.go#L21751) |
| Dictionary | ✅ | [builtinclasses.gen.go:20653](../pkg/builtin/builtinclasses.gen.go#L20653) |
| NodePath | ✅ | [builtinclasses.gen.go:19347](../pkg/builtin/builtinclasses.gen.go#L19347) |
| Callable | ✅ | [builtinclasses.gen.go:19864](../pkg/builtin/builtinclasses.gen.go#L19864) |
| Signal | ✅ | [builtinclasses.gen.go:20319](../pkg/builtin/builtinclasses.gen.go#L20319) |
| Packed*Array (all 10) | ✅ | [builtinclasses.gen.go:23120+](../pkg/builtin/builtinclasses.gen.go#L23120) |
| GDExtensionPropertyInfo | ✅ | [property_info.go:47](../pkg/ffi/property_info.go#L47) |
| GDExtensionClassMethodInfo | ✅ | [class_method_info.go:63](../pkg/ffi/class_method_info.go#L63) |

## Cleanup Helpers

For batch cleanup, use `util.DestroySlice`:

```go
import "github.com/godot-go/godot-go/pkg/util"

func cleanupVariants(values []builtin.Variant) {
	util.DestroySlice(values)
}
```

For scoped temporary values, prefer helper wrappers:

```go
import "github.com/godot-go/godot-go/pkg/builtin"

builtin.WithString("Player", func(name builtin.String) {
	// Use name here.
})

builtin.WithStringName("Node", func(nodeName builtin.StringName) {
	// Use nodeName here.
})
```

## Notes

- If you create a type listed above, you own it and must call `Destroy()`.
- If a function returns one of these types, assume you own it unless the API
  explicitly says Godot retains ownership.
