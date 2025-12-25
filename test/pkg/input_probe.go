package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// InputProbe implements GDClass evidence.
var _ GDClass = (*InputProbe)(nil)

type InputProbe struct {
	NodeImpl
	handleInput    bool
	inputCount     int32
	unhandledCount int32
}

func (p *InputProbe) GetClassName() string {
	return "InputProbe"
}

func (p *InputProbe) GetParentClassName() string {
	return "Node"
}

func (p *InputProbe) V_Ready() {
	p.SetProcessInput(true)
	p.SetProcessUnhandledInput(true)
}

func (p *InputProbe) V_Input(_event RefInputEvent) {
	p.inputCount++
	if p.handleInput {
		viewport := p.GetViewport()
		if viewport != nil {
			viewport.SetInputAsHandled()
		}
	}
}

func (p *InputProbe) V_UnhandledInput(_event RefInputEvent) {
	p.unhandledCount++
}

func (p *InputProbe) SetHandleInput(handled bool) {
	p.handleInput = handled
}

func (p *InputProbe) ResetCounts() {
	p.inputCount = 0
	p.unhandledCount = 0
}

func (p *InputProbe) GetInputCount() int32 {
	return p.inputCount
}

func (p *InputProbe) GetUnhandledCount() int32 {
	return p.unhandledCount
}

func NewInputProbeFromOwnerObject(owner *GodotObject) GDClass {
	obj := &InputProbe{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassInputProbe() {
	ClassDBRegisterClass(NewInputProbeFromOwnerObject, nil, nil, func(t *InputProbe) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Input", "_input", []string{"event"}, nil)
		ClassDBBindMethodVirtual(t, "V_UnhandledInput", "_unhandled_input", []string{"event"}, nil)
		ClassDBBindMethod(t, "SetHandleInput", "set_handle_input", []string{"handled"}, nil)
		ClassDBBindMethod(t, "ResetCounts", "reset_counts", nil, nil)
		ClassDBBindMethod(t, "GetInputCount", "get_input_count", nil, nil)
		ClassDBBindMethod(t, "GetUnhandledCount", "get_unhandled_count", nil, nil)
	})
}

func UnregisterClassInputProbe() {
	ClassDBUnregisterClass[*InputProbe]()
}
