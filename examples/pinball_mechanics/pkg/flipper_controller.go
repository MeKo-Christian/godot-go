package pkg

import (
	"strings"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
)

// FlipperController implements GDClass evidence.
var _ GDClass = (*FlipperController)(nil)

type FlipperController struct {
	Node2DImpl
	actionName  StringName
	initialized bool
	motorSpeed  float32
	joint       PinJoint2D
	flipperBody RigidBody2D
}

func (f *FlipperController) GetClassName() string {
	return "FlipperController"
}

func (f *FlipperController) GetParentClassName() string {
	return "Node2D"
}

func (f *FlipperController) V_Ready() {
	f.SetPhysicsProcess(true)
	f.configureSide()
	f.cacheNodes()
	f.configureJoint()
}

func (f *FlipperController) V_PhysicsProcess(_delta float64) {
	if !f.initialized || f.joint == nil {
		return
	}
	input := inputSingleton()
	if input == nil {
		return
	}
	if input.IsActionPressed(f.actionName, true) {
		f.joint.SetMotorTargetVelocity(f.motorSpeed)
	} else {
		f.joint.SetMotorTargetVelocity(0)
	}
}

func (f *FlipperController) V_ExitTree() {
	if f.initialized {
		f.actionName.Destroy()
	}
}

func (f *FlipperController) configureSide() {
	name := f.GetName()
	isLeft := strings.Contains(strings.ToLower(name.ToUtf8()), "left")
	name.Destroy()
	if isLeft {
		f.actionName = NewStringNameWithLatin1Chars(actionLeftFlipper)
		f.motorSpeed = 12.0
	} else {
		f.actionName = NewStringNameWithLatin1Chars(actionRightFlipper)
		f.motorSpeed = -12.0
	}
	f.initialized = true
}

func (f *FlipperController) cacheNodes() {
	bodyNode := f.GetNodeOrNull(nodePath("FlipperBody"))
	jointNode := f.GetNodeOrNull(nodePath("FlipperJoint"))
	anchorNode := f.GetNodeOrNull(nodePath("FlipperAnchor"))
	if bodyNode == nil || jointNode == nil || anchorNode == nil {
		printLine("FlipperController: missing child nodes")
		return
	}

	f.flipperBody, _ = ObjectCastTo(bodyNode, "RigidBody2D").(RigidBody2D)
	f.joint, _ = ObjectCastTo(jointNode, "PinJoint2D").(PinJoint2D)
	anchor, _ := ObjectCastTo(anchorNode, "StaticBody2D").(StaticBody2D)
	if f.flipperBody == nil || f.joint == nil || anchor == nil {
		printLine("FlipperController: child cast failed")
		return
	}

	anchorPath := f.joint.GetPathTo(anchor, false)
	defer anchorPath.Destroy()
	bodyPath := f.joint.GetPathTo(f.flipperBody, false)
	defer bodyPath.Destroy()
	f.joint.SetNodeA(anchorPath)
	f.joint.SetNodeB(bodyPath)
}

func (f *FlipperController) configureJoint() {
	if f.joint == nil {
		return
	}
	f.joint.SetMotorEnabled(true)
	f.joint.SetAngularLimitEnabled(true)
	f.joint.SetSoftness(0.2)
	if f.motorSpeed > 0 {
		f.joint.SetAngularLimitLower(-1.0)
		f.joint.SetAngularLimitUpper(0.2)
	} else {
		f.joint.SetAngularLimitLower(-0.2)
		f.joint.SetAngularLimitUpper(1.0)
	}
}

func NewFlipperControllerFromOwnerObject(owner *GodotObject) GDClass {
	obj := &FlipperController{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassFlipperController() {
	ClassDBRegisterClass(NewFlipperControllerFromOwnerObject, nil, nil, func(t *FlipperController) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_PhysicsProcess", "_physics_process", []string{"delta"}, nil)
		ClassDBBindMethodVirtual(t, "V_ExitTree", "_exit_tree", nil, nil)
	})
}

func UnregisterClassFlipperController() {
	ClassDBUnregisterClass[*FlipperController]()
}
