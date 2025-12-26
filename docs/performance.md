# Performance

This page collects practical guidance for profiling godot-go projects and
reducing overhead in hot paths.

## Profiling cgo overhead

godot-go uses cgo to cross the Go/Godot boundary. For tight loops, the cost of
cgo calls can dominate. Use the micro-benchmarks below to establish a baseline
on your machine.

```bash
# Cgo call overhead baseline

go test -bench=BenchmarkCgoTestCall -run ^$ ./pkg/util -benchmem

# Optional profiles for deeper analysis

go test -bench=BenchmarkCgoTestCall -run ^$ ./pkg/util -benchmem -cpuprofile cpu.out -memprofile mem.out
```

## Performance-critical patterns

- Minimize cross-boundary calls in per-frame loops. Prefer batch/bulk APIs.
- Cache `StringName`/`String` instances that are reused every frame.
- Avoid allocating Go slices in tight loops; reuse buffers where possible.
- Reuse `Array`/`Dictionary` values instead of recreating them each call.
- Prefer typed APIs over `Variant` conversions in performance-sensitive code.
- Keep `Callable` and method binds cached instead of recreating them on demand.

## Object pooling helpers

`pkg/pool` exposes small, reusable pools for common built-in types. Use them for
short-lived values inside hot loops.

```go
import "github.com/godot-go/godot-go/pkg/pool"

func buildPayload() {
	arr := pool.AcquireArray()
	defer pool.ReleaseArray(arr)

	v := pool.AcquireVariant()
	defer pool.ReleaseVariant(v)

	// Populate arr and v as needed...
}
```

## Physics-heavy benchmarks

Use Godot's profiler alongside a stress scene to isolate physics costs:

- Use the built-in benchmark harness in `test/demo` by setting
  `GODOT_GO_PHYSICS_BENCH=1` and running the demo in headless mode.
- Adjust parameters via environment variables:
  - `GODOT_GO_PHYSICS_BENCH_COUNT` (default: 1000)
  - `GODOT_GO_PHYSICS_BENCH_STEPS` (default: 300)
  - `GODOT_GO_PHYSICS_BENCH_RADIUS` (default: 6.0)
  - `GODOT_GO_PHYSICS_BENCH_SPACING` (default: 14.0)
  - `GODOT_GO_PHYSICS_BENCH_IMPULSE` (default: 12.0)
- Record `Performance.PHYSICS_TIME` and frame time over 300+ frames.
- Compare results between native GDScript and godot-go implementations.

This helps differentiate engine-side physics cost from Go/cgo overhead.
