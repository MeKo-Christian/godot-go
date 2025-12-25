# Hello World (Go)

Minimal godot-go extension that registers a `HelloWorld` node and prints a message on `_ready`.

## Build the extension

```bash
# From the repo root
CGO_ENABLED=1 \
GOOS=$(go env GOOS) \
GOARCH=$(go env GOARCH) \
go build -buildmode=c-shared -tags tools -trimpath \
  -o examples/hello_world/lib/libgodotgo-hello-world-$(go env GOOS)-$(go env GOARCH).so \
  examples/hello_world/main.go
```

On macOS/Windows, make the output filename match the entries in `hello_world.gdextension` (or edit the `.gdextension` file to match your output).

## Run in Godot

```bash
GODOT=/path/to/godot $GODOT --path examples/hello_world/
```

You should see "Hello from godot-go" printed in the output.
