package pkg

import (
	"fmt"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// ScoreSystem implements GDClass evidence.
var _ GDClass = (*ScoreSystem)(nil)

type ScoreSystem struct {
	ControlImpl
	score int64
}

func (s *ScoreSystem) GetClassName() string {
	return "ScoreSystem"
}

func (s *ScoreSystem) GetParentClassName() string {
	return "Control"
}

func (s *ScoreSystem) V_Ready() {
	s.updateLabel()
}

func (s *ScoreSystem) OnBumperScored(points int64) {
	s.score += points
	s.updateLabel()
}

func (s *ScoreSystem) OnBallDrained(_body Node2D) {
	s.score = 0
	s.updateLabel()
}

func (s *ScoreSystem) updateLabel() {
	label := s.GetNodeOrNull(nodePath("ScoreLabel"))
	if label == nil {
		printLine("ScoreSystem: Label not found")
		return
	}
	labelNode, ok := ObjectCastTo(label, "Label").(Label)
	if !ok || labelNode == nil {
		printLine("ScoreSystem: Label cast failed")
		return
	}
	text := NewStringWithUtf8Chars(fmt.Sprintf("Score: %d", s.score))
	defer text.Destroy()
	labelNode.SetText(text)
}

func NewScoreSystemFromOwnerObject(owner *GodotObject) GDClass {
	obj := &ScoreSystem{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassScoreSystem() {
	ClassDBRegisterClass(NewScoreSystemFromOwnerObject, nil, nil, func(t *ScoreSystem) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethod(t, "OnBumperScored", "_on_bumper_scored", []string{"points"}, nil)
		ClassDBBindMethod(t, "OnBallDrained", "_on_ball_drained", []string{"body"}, nil)
	})
}

func UnregisterClassScoreSystem() {
	ClassDBUnregisterClass[*ScoreSystem]()
}
