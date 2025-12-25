# Actions Input Demo

Minimal input-actions demo that registers actions in `InputMap` and polls them with `Input.IsActionJustPressed`.

Actions:

- `demo_jump` bound to Space
- `demo_click` bound to Left Mouse Button

## Build the extension

```bash
# From the repo root
CGO_ENABLED=1 \
GOOS=$(go env GOOS) \
GOARCH=$(go env GOARCH) \
go build -buildmode=c-shared -tags tools -trimpath \
  -o examples/actions_input_demo/lib/libgodotgo-actions-demo-$(go env GOOS)-$(go env GOARCH).so \
  examples/actions_input_demo/main.go
```

On macOS/Windows, make the output filename match the entries in `actions_demo.gdextension`.

## Run in Godot

```bash
GODOT=/path/to/godot $GODOT --path examples/actions_input_demo/
```

Press Space or click to increment the score label.
