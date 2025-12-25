# Visual Effects

This guide covers GPUParticles2D, AnimatedSprite2D frame control, and shader materials.

## GPUParticles2D

Use a `ParticleProcessMaterial` for procedural particle behavior:

- `SetDirection` and `SetSpread` define emission direction.
- `SetInitialVelocity` and `SetGravity` shape motion.
- Assign the material with `GPUParticles2D.SetProcessMaterial`.
- Set `SetEmitting(true)` to start the effect.

The example at `examples/audio_visual_fx/` initializes a ParticleProcessMaterial at runtime.

## AnimatedSprite2D frame control

Create `SpriteFrames` in Go and add frames via `AddAnimation` + `AddFrame`.

Manual stepping is useful for frame-precise control:

```go
sprite.SetPlaying(false)
sprite.SetFrame(0)
// then advance frames in _process
```

## Shader materials (CanvasItem)

Apply a custom shader to `ColorRect` (or any CanvasItem-derived node):

- Create a `Shader` resource and call `SetCode`.
- Create `ShaderMaterial` and assign the shader.
- Use `SetMaterial` on the target node.

The demo shader in `examples/audio_visual_fx/` uses `TIME` to animate a gradient.

## Optional: CanvasItem drawing

For custom drawing, implement `_draw` and call CanvasItem draw methods (e.g., `DrawLine`, `DrawCircle`). This is useful for debug overlays or procedural UI, but shaders and sprites cover most visual effects.
