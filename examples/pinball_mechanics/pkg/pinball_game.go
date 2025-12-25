package pkg

import (
	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// PinballGame implements GDClass evidence.
var _ GDClass = (*PinballGame)(nil)

type PinballGame struct {
	Node2DImpl
}

func (p *PinballGame) GetClassName() string {
	return "PinballGame"
}

func (p *PinballGame) GetParentClassName() string {
	return "Node2D"
}

func (p *PinballGame) V_Ready() {
	setupPinballActions()
	p.connectSignals()
}

func (p *PinballGame) connectSignals() {
	scoreNode := p.GetNodeOrNull(nodePath("UI/ScoreSystem"))
	if scoreNode == nil {
		printLine("PinballGame: ScoreSystem not found")
		return
	}
	scoreSystem, ok := ObjectCastTo(scoreNode, "ScoreSystem").(*ScoreSystem)
	if !ok || scoreSystem == nil {
		printLine("PinballGame: ScoreSystem cast failed")
		return
	}

	connectBumper(scoreSystem, p.GetNodeOrNull(nodePath("BumperLeft")))
	connectBumper(scoreSystem, p.GetNodeOrNull(nodePath("BumperRight")))

	drainNode := p.GetNodeOrNull(nodePath("Drain"))
	if drainNode != nil {
		connectSignalTo(scoreSystem, drainNode, "drained", "_on_ball_drained")
	}

	launcherNode := p.GetNodeOrNull(nodePath("BallLauncher"))
	if launcherNode != nil {
		connectSignalTo(ObjectCastTo(launcherNode, "BallLauncher"), drainNode, "drained", "_on_ball_drained")
	}
}

func connectBumper(scoreSystem *ScoreSystem, node Node) {
	if node == nil {
		return
	}
	connectSignalTo(scoreSystem, node, "bumped", "_on_bumper_scored")
}

func connectSignalTo(target Object, node Node, signalName string, methodName string) {
	if node == nil || target == nil {
		return
	}
	signal := NewStringNameWithLatin1Chars(signalName)
	defer signal.Destroy()
	method := NewStringNameWithLatin1Chars(methodName)
	defer method.Destroy()
	callable := NewCallableWithObjectStringName(target, method)
	defer callable.Destroy()
	node.Connect(signal, callable, 0)
}

func NewPinballGameFromOwnerObject(owner *GodotObject) GDClass {
	obj := &PinballGame{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassPinballGame() {
	ClassDBRegisterClass(NewPinballGameFromOwnerObject, nil, nil, func(t *PinballGame) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
	})
}

func UnregisterClassPinballGame() {
	ClassDBUnregisterClass[*PinballGame]()
}
