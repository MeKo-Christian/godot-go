# Input Responsiveness

This doc covers input handling patterns that minimize latency and clarify `_input` vs `_unhandled_input` behavior.

## V_Input vs V_UnhandledInput

- `V_Input` (`_input`) receives raw input events first.
- `V_UnhandledInput` (`_unhandled_input`) runs after `_input` if the event was not marked as handled.
- Call `GetViewport().SetInputAsHandled()` inside `V_Input` to stop propagation to `_unhandled_input`.

Use `V_Input` for latency-sensitive actions (flippers), and reserve `V_UnhandledInput` for UI/secondary handling.

## Rapid Input Handling

Godot queues input events and delivers them during the frame. To keep flipper response tight:

- Prefer polling in `_physics_process` using `Input.IsActionPressed` or `IsActionJustPressed` for frame-stable results.
- Use `InputMap` actions so multiple devices map to the same action.
- Avoid heavy work in `V_Input`; instead, set state and act in `_physics_process`.

## Input Buffering Notes

- `Input.SetUseAccumulatedInput(true)` affects mouse motion accumulation, not keyboard events.
- For keyboard input, each event is delivered as-is; no extra buffering is applied.
- If you need buffering (e.g., to accept input within a short window), implement it explicitly by timestamping events in `V_Input` and consuming them in `_physics_process`.
