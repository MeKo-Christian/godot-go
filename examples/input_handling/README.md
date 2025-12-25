# Input Handling Demo

Minimal input demo that prints keyboard and mouse events using a Go-backed `InputDemo` node.

## Build the extension

```bash
# From the repo root
CGO_ENABLED=1 \
GOOS=$(go env GOOS) \
GOARCH=$(go env GOARCH) \
go build -buildmode=c-shared -tags tools -trimpath \
  -o examples/input_handling/lib/libgodotgo-input-demo-$(go env GOOS)-$(go env GOARCH).so \
  examples/input_handling/main.go
```

On macOS/Windows, make the output filename match the entries in `input_demo.gdextension`.

## Run in Godot

```bash
GODOT=/path/to/godot $GODOT --path examples/input_handling/
```

Use the keyboard or mouse to see events printed in the output.
