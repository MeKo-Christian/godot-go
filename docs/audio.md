# Audio

This guide covers AudioStreamPlayer2D usage, resource loading, and spatial audio setup in godot-go.

## AudioStreamPlayer2D + resource loading

The simplest path is to load a stream with ResourceLoader and assign it to an AudioStreamPlayer2D:

```go
loader := getResourceLoaderSingleton()
path := NewStringWithUtf8Chars("res://sfx/beep.wav")
typeHint := NewStringWithUtf8Chars("AudioStream")
stream := loader.Load(path, typeHint, RESOURCE_LOADER_CACHE_MODE_CACHE_MODE_REUSE)
```

Then cast the resource to `AudioStream` and call `SetStream` + `Play` on the player.

See `examples/audio_visual_fx/` for a full working setup.

## Spatial audio setup

For 2D spatial audio, configure the player and move it relative to the listener:

- `SetMaxDistance` controls how far the sound can be heard.
- `SetAttenuation` changes how quickly volume drops with distance.
- `SetPanningStrength` affects left/right stereo separation.
- `SetVolumeDb` gives a predictable default loudness.

Moving the AudioStreamPlayer2D node across the scene will create panning and attenuation automatically.

## Cleanup reminders

`String` and `StringName` values are manual-destroy types. Always `Destroy()` them when loading resources or building NodePath values.
