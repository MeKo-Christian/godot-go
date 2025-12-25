# 2D Physics: Bouncing Ball

Minimal 2D physics demo using a `RigidBody2D` subclass implemented in Go. The ball is given an impulse on `_ready` and bounces off static walls.

## Build the extension

```bash
# From the repo root
CGO_ENABLED=1 \
GOOS=$(go env GOOS) \
GOARCH=$(go env GOARCH) \
go build -buildmode=c-shared -tags tools -trimpath \
  -o examples/physics_2d_bouncing_ball/lib/libgodotgo-bouncing-ball-$(go env GOOS)-$(go env GOARCH).so \
  examples/physics_2d_bouncing_ball/main.go
```

On macOS/Windows, make the output filename match the entries in `bouncing_ball.gdextension`.

## Run in Godot

```bash
GODOT=/path/to/godot $GODOT --path examples/physics_2d_bouncing_ball/
```

Tip: enable **Debug > Visible Collision Shapes** to see the physics bodies if you disable the polygon visual.
