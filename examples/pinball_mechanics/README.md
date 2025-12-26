# Pinball Mechanics Demo

This demo showcases pinball-focused mechanics implemented in Go:

- Flipper controller (RigidBody2D + PinJoint2D + input actions)
- Ball launcher (impulse application)
- Bumper (Area2D trigger + impulse + score signal)
- Drain detection (Area2D trigger)
- Score system (signals + UI label)

## Build

```bash
# from repo root
CGO_ENABLED=1 go build \
  -o examples/pinball_mechanics/lib/libgodotgo-pinball-mechanics-$(go env GOOS)-$(go env GOARCH).so \
  examples/pinball_mechanics/main.go
```

On macOS/Windows, make the output filename match the entries in `pinball_mechanics.gdextension`.

## Run

```bash
GODOT=/path/to/godot $GODOT --path examples/pinball_mechanics/
```

Controls:

- Left arrow: left flipper
- Right arrow: right flipper
- Space: launch ball
