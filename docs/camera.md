# Camera & Viewport

This guide covers Camera2D following and viewport size handling.

## Camera2D follow

A common pattern is to parent Camera2D under the target node (player, ball, etc.). When the target moves, the camera follows automatically.

Recommended settings for smooth motion:

- `SetPositionSmoothingEnabled(true)`
- `SetPositionSmoothingSpeed(4-8)`
- `SetZoom(Vector2)` for framing

See `examples/camera_viewport/` for a working setup.

## Viewport size handling

Use the viewport's `size_changed` signal to keep UI in sync with resizes:

- Connect `Viewport.size_changed` to a Go method (e.g. `_on_viewport_size_changed`).
- Query the size with `GetVisibleRect()` and update labels or layout accordingly.

The camera demo updates a label with the current viewport dimensions.
