# UI Overlay: Score Display

Minimal UI overlay demo with a Go-backed `ScoreOverlay` node that updates a `Label` every second.

## Build the extension

```bash
# From the repo root
CGO_ENABLED=1 \
GOOS=$(go env GOOS) \
GOARCH=$(go env GOARCH) \
go build -buildmode=c-shared -tags tools -trimpath \
  -o examples/ui_overlay_score/lib/libgodotgo-score-overlay-$(go env GOOS)-$(go env GOARCH).so \
  examples/ui_overlay_score/main.go
```

On macOS/Windows, make the output filename match the entries in `score_overlay.gdextension`.

## Run in Godot

```bash
GODOT=/path/to/godot $GODOT --path examples/ui_overlay_score/
```
