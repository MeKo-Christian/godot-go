# Memory Ownership and Manual Cleanup

godot-go exposes several Godot types that allocate native resources. These
must be released by calling `Destroy()` when you are done with them.

## Types Requiring Manual Cleanup

| Type                       | Requires `Destroy()` | Location                                                                    |
| -------------------------- | -------------------- | --------------------------------------------------------------------------- |
| String                     | ✅                   | [builtinclasses.gen.go:210](../pkg/builtin/builtinclasses.gen.go#L210)      |
| StringName                 | ✅                   | [builtinclasses.gen.go:16358](../pkg/builtin/builtinclasses.gen.go#L16358)  |
| Variant                    | ✅                   | [variant.go:394](../pkg/builtin/variant.go#L394)                            |
| Array                      | ✅                   | [builtinclasses.gen.go:21751](../pkg/builtin/builtinclasses.gen.go#L21751)  |
| Dictionary                 | ✅                   | [builtinclasses.gen.go:20653](../pkg/builtin/builtinclasses.gen.go#L20653)  |
| NodePath                   | ✅                   | [builtinclasses.gen.go:19347](../pkg/builtin/builtinclasses.gen.go#L19347)  |
| Callable                   | ✅                   | [builtinclasses.gen.go:19864](../pkg/builtin/builtinclasses.gen.go#L19864)  |
| Signal                     | ✅                   | [builtinclasses.gen.go:20319](../pkg/builtin/builtinclasses.gen.go#L20319)  |
| Packed\*Array (all 10)     | ✅                   | [builtinclasses.gen.go:23120+](../pkg/builtin/builtinclasses.gen.go#L23120) |
| GDExtensionPropertyInfo    | ✅                   | [property_info.go:47](../pkg/ffi/property_info.go#L47)                      |
| GDExtensionClassMethodInfo | ✅                   | [class_method_info.go:63](../pkg/ffi/class_method_info.go#L63)              |

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

## Finalizers (Not Used)

godot-go does not use `runtime.SetFinalizer` for `String`, `StringName`, or
other builtins. Finalizers are a poor fit here for a few reasons:

- Builtin types are small value types that get copied freely; a finalizer would
  run on a pointer to one copy and can double-free or free the wrong owner.
- Finalizers run on the GC thread at unpredictable times; calling into Godot
  from a finalizer is unsafe because the API is not thread-safe and the engine
  may already be shutting down.
- Cgo finalizers can race with in-flight pinned data and lifetime assumptions,
  which makes failures hard to debug and nondeterministic.

Instead, explicit `Destroy()` calls and scoped helpers are the supported
approach. Leak detection is tracked via long-running tests (see Task 1.7).

## Leak Test

The demo project includes a leak test that runs the game loop for a fixed
duration and checks Go heap growth via `runtime.MemStats`.

Run it with:

```bash
just leak_test
```

Environment overrides:

- `GODOT_GO_LEAK_TEST_SECONDS` (default: 600)
- `GODOT_GO_LEAK_TEST_INTERVAL_MS` (default: 100)
- `GODOT_GO_LEAK_TEST_ITERATIONS` (default: 1000)
- `GODOT_GO_LEAK_TEST_MAX_HEAP_BYTES` (default: 10485760)
- `GODOT_GO_LEAK_TEST_MAX_HEAP_OBJECTS` (default: 5000)

## Notes

- If you create a type listed above, you own it and must call `Destroy()`.
- If a function returns one of these types, assume you own it unless the API
  explicitly says Godot retains ownership.
