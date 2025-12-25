# Audio + Visual FX Demo

This demo validates AudioStreamPlayer2D playback, resource loading, GPUParticles2D, AnimatedSprite2D frame control, and shader materials.

## What it shows

- Loads `res://sfx/beep.wav` via `ResourceLoader` and plays it through AudioStreamPlayer2D.
- Moves the audio source to demonstrate spatial panning/attenuation.
- Runs GPUParticles2D with a ParticleProcessMaterial.
- Generates AnimatedSprite2D frames at runtime and steps frames manually.
- Applies a custom canvas-item shader to a ColorRect.

## Run

1. Build the extension for this example.
2. Open `examples/audio_visual_fx/project.godot` in Godot 4.5+.
3. Run the scene.

The beep should pan left-to-right while particles and visual effects animate.
