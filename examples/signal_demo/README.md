# Signal Demo

Minimal signal demo with two Go classes: `SignalEmitter` emits a custom signal, and `SignalListener` receives it.

## Build the extension

```bash
# From the repo root
CGO_ENABLED=1 \
GOOS=$(go env GOOS) \
GOARCH=$(go env GOARCH) \
go build -buildmode=c-shared -tags tools -trimpath \
  -o examples/signal_demo/lib/libgodotgo-signal-demo-$(go env GOOS)-$(go env GOARCH).so \
  examples/signal_demo/main.go
```

On macOS/Windows, make the output filename match the entries in `signal_demo.gdextension`.

## Run in Godot

```bash
GODOT=/path/to/godot $GODOT --path examples/signal_demo/
```

You should see log output confirming the signal connection and emission.
