package pkg

import (
	"fmt"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// ScoreOverlay implements GDClass evidence.
var _ GDClass = (*ScoreOverlay)(nil)

type ScoreOverlay struct {
	ControlImpl
	score int64
}

func (s *ScoreOverlay) GetClassName() string {
	return "ScoreOverlay"
}

func (s *ScoreOverlay) GetParentClassName() string {
	return "Control"
}

func (s *ScoreOverlay) V_Ready() {
	s.updateLabel()
	timer := s.GetNodeOrNull(nodePath("ScoreTimer"))
	if timer == nil {
		printLine("ScoreOverlay: Timer not found")
		return
	}

	timerNode := ObjectCastTo(timer, "Timer").(Timer)
	signal := NewStringNameWithLatin1Chars("timeout")
	defer signal.Destroy()
	method := NewStringNameWithLatin1Chars("_on_timer_timeout")
	defer method.Destroy()
	callable := NewCallableWithObjectStringName(s, method)
	defer callable.Destroy()
	timerNode.Connect(signal, callable, 0)
	timerNode.Start(1)
}

func (s *ScoreOverlay) OnTimerTimeout() {
	s.score += 10
	s.updateLabel()
}

func (s *ScoreOverlay) updateLabel() {
	label := s.GetNodeOrNull(nodePath("ScoreLabel"))
	if label == nil {
		printLine("ScoreOverlay: Label not found")
		return
	}
	labelNode := ObjectCastTo(label, "Label").(Label)
	text := NewStringWithUtf8Chars(fmt.Sprintf("Score: %d", s.score))
	defer text.Destroy()
	labelNode.SetText(text)
}

func NewScoreOverlayFromOwnerObject(owner *GodotObject) GDClass {
	obj := &ScoreOverlay{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassScoreOverlay() {
	ClassDBRegisterClass(NewScoreOverlayFromOwnerObject, nil, nil, func(t *ScoreOverlay) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethod(t, "OnTimerTimeout", "_on_timer_timeout", nil, nil)
	})
}

func UnregisterClassScoreOverlay() {
	ClassDBUnregisterClass[*ScoreOverlay]()
}

func nodePath(path string) NodePath {
	str := NewStringWithUtf8Chars(path)
	defer str.Destroy()
	return NewNodePathWithString(str)
}

func printLine(text string) {
	v := NewVariantGoString(text)
	defer v.Destroy()
	Print(v)
}
