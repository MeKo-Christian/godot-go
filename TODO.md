# godot-go Roadmap

> **Goal:** Make godot-go production-ready for developing a 2D pinball game with Godot 4.5+

## Current State Summary

### What's Working

- **971 Godot classes** fully generated with Go bindings
- Class registration and method binding via ClassDB
- Virtual methods (`V_Ready()`, `V_Input()`, etc.)
- Physics: CharacterBody2D, RigidBody2D, Area2D, CollisionShape2D
- Input: Full InputEvent hierarchy, keyboard/mouse handling
- Math: Vector2/3/4, Transform2D/3D, Quaternion, etc.
- Containers: Array, Dictionary, Packed\*Array types
- Signals: Custom signal definition and emission
- Animation: AnimationPlayer, AnimatedSprite2D
- UI: Full Control hierarchy

### Known Issues

- Memory leaks (documented in README)
- String handling quirks between Go/Godot types
- No coroutine/await support (use timers instead)
- Error handling (`rError`) incomplete

---

## Phase 1: Stability & Memory Safety

> **Priority:** Critical — fixes must come before new features

### Task 1.1: Fix Critical Memory Leaks (Commented-Out Destroys)

These are explicit `Destroy()` calls that were commented out — uncomment or fix:

- [ ] [variant.go:453-470](pkg/builtin/variant.go#L453-L470) — `ZapVector2/3/4` logging functions leak Variants
- [ ] [variant.go:134](pkg/builtin/variant.go#L134) — `getObjectInstanceBinding()` StringName leak
- [ ] [classdb.go:197](pkg/core/classdb.go#L197) — Signal `PropertyInfo` array never destroyed
- [ ] [method_bind.go:353](pkg/core/method_bind.go#L353) — Method name StringName not destroyed
- [ ] [method_bind_reflect.go:206,515,524](pkg/core/method_bind_reflect.go#L206) — Multiple StringName leaks
- [ ] [variant_string_encoder.go:69,82](pkg/builtin/variant_string_encoder.go#L69) — Encoder Destroy commented
- [ ] [char_string.go:65](pkg/builtin/char_string.go#L65) — Buffer Destroy commented

### Task 1.2: Implement rError Checking

Currently ALL GDScript→Go method calls silently ignore errors:

- [ ] [method_bind_callback.go:30](pkg/core/method_bind_callback.go#L30) — `rError` param never checked (critical)
- [ ] Add error logging when `rError` indicates failure
- [ ] Propagate errors to Go callers where possible
- [ ] Update code generator templates to validate `rError` in generated bindings

### Task 1.3: Document Types Requiring Manual Cleanup

Create `docs/memory.md` documenting ownership for these types:

| Type | Requires `Destroy()` | Location |
|------|---------------------|----------|
| String | ✅ | [builtinclasses.gen.go:210](pkg/builtin/builtinclasses.gen.go#L210) |
| StringName | ✅ | [builtinclasses.gen.go:16358](pkg/builtin/builtinclasses.gen.go#L16358) |
| Variant | ✅ | [variant.go:394](pkg/builtin/variant.go#L394) |
| Array | ✅ | [builtinclasses.gen.go:21751](pkg/builtin/builtinclasses.gen.go#L21751) |
| Dictionary | ✅ | [builtinclasses.gen.go:20653](pkg/builtin/builtinclasses.gen.go#L20653) |
| NodePath | ✅ | [builtinclasses.gen.go:19347](pkg/builtin/builtinclasses.gen.go#L19347) |
| Callable | ✅ | [builtinclasses.gen.go:19864](pkg/builtin/builtinclasses.gen.go#L19864) |
| Signal | ✅ | [builtinclasses.gen.go:20319](pkg/builtin/builtinclasses.gen.go#L20319) |
| Packed*Array (all 10) | ✅ | [builtinclasses.gen.go:23120+](pkg/builtin/builtinclasses.gen.go#L23120) |
| GDExtensionPropertyInfo | ✅ | [property_info.go:47](pkg/ffi/property_info.go#L47) |
| GDExtensionClassMethodInfo | ✅ | [class_method_info.go:63](pkg/ffi/class_method_info.go#L63) |

### Task 1.4: Add Cleanup Helpers

- [ ] Create `DestroySlice[T]()` helper for batch cleanup
- [ ] Add `defer`-friendly patterns for common workflows
- [ ] Consider RAII wrapper: `WithString(s, func(str String) { ... })`

### Task 1.5: Evaluate Finalizers

Currently **zero** `runtime.SetFinalizer` usage in codebase:

- [ ] Evaluate adding finalizers for String/StringName as safety net
- [ ] Document why finalizers may not work (cgo pinning, GC timing)
- [ ] If not using finalizers, add leak detection tooling instead

### Task 1.6: Fix Pinner Lifecycle

Global `runtime.Pinner` objects are never unpinned:

- [ ] [builtin/lib.go:35](pkg/builtin/lib.go#L35) — `pnr` never unpinned
- [ ] [ffi/lib.go:18](pkg/ffi/lib.go#L18) — `pnr` never unpinned
- [ ] [core/lib.go:35](pkg/core/lib.go#L35) — `pnr` never unpinned
- [ ] Evaluate if unpinning is safe/necessary for long-running games

### Task 1.7: Memory Leak Detection Tests

- [ ] Create test that runs game loop for 10+ minutes
- [ ] Monitor memory growth with `runtime.MemStats`
- [ ] Add CI job that fails on memory growth above threshold

---

## Phase 2: Developer Experience

> **Priority:** High — makes the library usable for real projects

- [ ] **Documentation**
  - [ ] Complete API reference for common game dev classes
  - [ ] Add "Getting Started" guide with minimal project setup
  - [ ] Document GDScript-to-Go patterns (cheat sheet)
  - [ ] Add physics examples (RigidBody2D, impulses, collisions)
  - [ ] Document signal connection patterns

- [ ] **Examples**
  - [ ] Create `examples/` directory with standalone demos
  - [ ] Minimal "Hello World" scene
  - [ ] 2D physics demo (bouncing ball)
  - [ ] Input handling demo (keyboard + mouse)
  - [ ] Signal communication demo
  - [ ] UI overlay demo (score display)

- [ ] **Tooling**
  - [ ] Add `just new-class <ClassName>` generator for boilerplate
  - [ ] Improve error messages from code generator
  - [ ] Add VS Code snippets for common patterns

---

## Phase 3: Physics & Game Mechanics (Pinball Focus)

> **Priority:** High — enables the target use case

- [ ] **Physics validation**
  - [ ] Test RigidBody2D with continuous collision (ball physics)
  - [ ] Test impulse application (flipper hits)
  - [ ] Test Area2D triggers (bumpers, drain detection)
  - [ ] Verify physics material properties (friction, bounce)
  - [ ] Test joint constraints (PinJoint2D for flippers)

- [ ] **Create pinball-specific examples**
  - [ ] Flipper controller (RigidBody2D + input + joint)
  - [ ] Ball launcher (impulse application)
  - [ ] Bumper (Area2D + signal + impulse response)
  - [ ] Drain detection (Area2D trigger)
  - [ ] Score system (signals + UI binding)

- [ ] **Input responsiveness**
  - [ ] Test rapid input handling (critical for flipper response)
  - [ ] Verify `V_Input` vs `V_UnhandledInput` behavior
  - [ ] Document input buffering if needed

---

## Phase 4: Audio & Visual Polish

> **Priority:** Medium — enhances game feel

- [ ] **Audio**
  - [ ] Test AudioStreamPlayer2D with godot-go
  - [ ] Verify audio loading from resources
  - [ ] Create audio playback example
  - [ ] Document spatial audio setup

- [ ] **Visual effects**
  - [ ] Test GPUParticles2D for ball impacts
  - [ ] Verify AnimatedSprite2D frame control
  - [ ] Test shader material application
  - [ ] Document canvas drawing (if needed)

- [ ] **Camera & Viewport**
  - [ ] Test Camera2D following
  - [ ] Verify viewport scaling/resolution handling

---

## Phase 5: Production Readiness

> **Priority:** Medium — for shipping games

- [ ] **Performance**
  - [ ] Profile cgo overhead in hot loops
  - [ ] Identify and document performance-critical patterns
  - [ ] Add object pooling utilities for common types
  - [ ] Benchmark physics-heavy scenarios

- [ ] **Testing infrastructure**
  - [ ] Add unit tests for all builtin types
  - [ ] Add integration tests for class registration
  - [ ] Add physics simulation tests
  - [ ] CI: Run tests on multiple platforms

- [ ] **Build & Distribution**
  - [ ] Document cross-compilation for Windows/Mac/Linux
  - [ ] Test export workflow
  - [ ] Document debugging compiled extensions

---

## Phase 6: Advanced Features (Future)

> **Priority:** Low — nice to have

- [ ] **Goroutine integration**
  - [ ] Safe goroutine usage with Godot thread model
  - [ ] Signal-to-channel bridges
  - [ ] Async resource loading patterns

- [ ] **Editor integration**
  - [ ] Custom inspector properties
  - [ ] Tool scripts (running in editor)
  - [ ] Custom resource types

- [ ] **3D support validation**
  - [ ] Test RigidBody3D, CharacterBody3D
  - [ ] Verify 3D collision detection
  - [ ] Camera3D and viewport handling

---

## Technical Debt

### Code Quality

- [ ] Remove/resolve commented-out code (e.g., ExampleRef)
- [ ] Consistent error handling patterns across codebase
- [ ] Add linting for generated code quality

### Architecture

- [ ] Evaluate class-scoped virtual naming (`V_ClassName_Method`)
- [ ] Consider builder pattern for complex type construction
- [ ] Evaluate interface-based type system improvements

### Build System

- [ ] Parallel generation for faster builds
- [ ] Incremental generation (only regenerate changed classes)
- [ ] Better dependency tracking

---

## Reference Notes

### Property Hint Logic (needs placement)

```go
// behavior ported from godot-cpp
switch hint {
case PROPERTY_HINT_RESOURCE_TYPE:
    className = hintString
default:
    className = pClassName
}
```

### cgo Performance Resources

- [GopherCon 2018 - Adventures in Cgo Performance](https://about.sourcegraph.com/blog/go/gophercon-2018-adventures-in-cgo-performance)
- [FFI Overhead Benchmarks](https://github.com/dyu/ffi-overhead)

---

## Pinball Game Checklist

When Phase 3 is complete, you should be able to build:

- [ ] Ball physics with realistic bouncing
- [ ] Responsive flipper controls (<16ms input latency)
- [ ] Bumpers that apply impulses on contact
- [ ] Ramps and lanes using collision shapes
- [ ] Ball drain detection with signals
- [ ] Score tracking and display
- [ ] Multi-ball support
- [ ] Sound effects on collisions
- [ ] Basic particle effects

---

_Last updated: 2025-12-25_
