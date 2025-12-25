package pkg

import (
	"math"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	"github.com/godot-go/godot-go/pkg/log"
	"go.uber.org/zap"
)

type PhysicsValidation struct {
	Node2DImpl
	areaEnterCount int32
	areaExitCount  int32
}

func (p *PhysicsValidation) GetClassName() string {
	return "PhysicsValidation"
}

func (p *PhysicsValidation) GetParentClassName() string {
	return "Node2D"
}

func (p *PhysicsValidation) EnableCCD(body RigidBody2D) bool {
	if body == nil {
		log.Warn("EnableCCD called with nil body")
		return false
	}
	body.SetContinuousCollisionDetectionMode(RigidBody2DCCDModeContinuous)
	return body.GetContinuousCollisionDetectionMode() == RigidBody2DCCDModeContinuous
}

func (p *PhysicsValidation) ApplyFlipperImpulse(body RigidBody2D, impulse Vector2, position Vector2) {
	if body == nil {
		log.Warn("ApplyFlipperImpulse called with nil body")
		return
	}
	body.ApplyImpulse(impulse, position)
}

func (p *PhysicsValidation) GetLinearSpeed(body RigidBody2D) float32 {
	if body == nil {
		log.Warn("GetLinearSpeed called with nil body")
		return 0
	}
	vel := body.GetLinearVelocity()
	return vel.Length()
}

func (p *PhysicsValidation) ConfigureMaterial(body RigidBody2D, material RefPhysicsMaterial, friction, bounce float32) bool {
	if body == nil {
		log.Warn("ConfigureMaterial called with nil body")
		return false
	}
	if material == nil || !material.IsValid() {
		log.Warn("ConfigureMaterial called with invalid material")
		return false
	}
	mat := material.TypedPtr()
	mat.SetFriction(friction)
	mat.SetBounce(bounce)
	body.SetPhysicsMaterialOverride(material)

	assigned := body.GetPhysicsMaterialOverride()
	if assigned == nil || !assigned.IsValid() {
		log.Warn("ConfigureMaterial failed to read override material")
		return false
	}
	assignedMat := assigned.TypedPtr()
	gotFriction := assignedMat.GetFriction()
	gotBounce := assignedMat.GetBounce()
	return approxEqual32(gotFriction, friction) && approxEqual32(gotBounce, bounce)
}

func (p *PhysicsValidation) BindArea(area Area2D) bool {
	if area == nil {
		log.Warn("BindArea called with nil area")
		return false
	}
	signalEntered := NewStringNameWithLatin1Chars("body_entered")
	defer signalEntered.Destroy()
	methodEntered := NewStringNameWithLatin1Chars("_on_area_body_entered")
	defer methodEntered.Destroy()
	callableEntered := NewCallableWithObjectStringName(p, methodEntered)
	defer callableEntered.Destroy()

	signalExited := NewStringNameWithLatin1Chars("body_exited")
	defer signalExited.Destroy()
	methodExited := NewStringNameWithLatin1Chars("_on_area_body_exited")
	defer methodExited.Destroy()
	callableExited := NewCallableWithObjectStringName(p, methodExited)
	defer callableExited.Destroy()

	enteredErr := area.Connect(signalEntered, callableEntered, 0)
	exitedErr := area.Connect(signalExited, callableExited, 0)
	if enteredErr != OK || exitedErr != OK {
		log.Warn("BindArea connect failed",
			zap.Any("entered_err", enteredErr),
			zap.Any("exited_err", exitedErr),
		)
		return false
	}
	return true
}

func (p *PhysicsValidation) ResetAreaCounts() {
	p.areaEnterCount = 0
	p.areaExitCount = 0
}

func (p *PhysicsValidation) GetAreaEnterCount() int32 {
	return p.areaEnterCount
}

func (p *PhysicsValidation) GetAreaExitCount() int32 {
	return p.areaExitCount
}

func (p *PhysicsValidation) OnAreaBodyEntered(body Node2D) {
	p.areaEnterCount++
	if body != nil {
		log.Debug("Area body entered", zap.Uint64("body_id", body.GetInstanceId()))
	}
}

func (p *PhysicsValidation) OnAreaBodyExited(body Node2D) {
	p.areaExitCount++
	if body != nil {
		log.Debug("Area body exited", zap.Uint64("body_id", body.GetInstanceId()))
	}
}

func (p *PhysicsValidation) ConfigurePinJoint(joint PinJoint2D, nodeAPath, nodeBPath string, lower, upper, softness, motorVelocity float32) bool {
	if joint == nil {
		log.Warn("ConfigurePinJoint called with nil joint")
		return false
	}
	nodeAString := NewStringWithUtf8Chars(nodeAPath)
	defer nodeAString.Destroy()
	nodeBString := NewStringWithUtf8Chars(nodeBPath)
	defer nodeBString.Destroy()

	nodeA := NewNodePathWithString(nodeAString)
	defer nodeA.Destroy()
	nodeB := NewNodePathWithString(nodeBString)
	defer nodeB.Destroy()

	joint.SetNodeA(nodeA)
	joint.SetNodeB(nodeB)
	joint.SetAngularLimitEnabled(true)
	joint.SetAngularLimitLower(lower)
	joint.SetAngularLimitUpper(upper)
	joint.SetMotorEnabled(true)
	joint.SetMotorTargetVelocity(motorVelocity)
	joint.SetSoftness(softness)

	gotNodeA := joint.GetNodeA()
	defer gotNodeA.Destroy()
	gotNodeB := joint.GetNodeB()
	defer gotNodeB.Destroy()

	okPaths := nodeA.Equal_NodePath(gotNodeA) && nodeB.Equal_NodePath(gotNodeB)
	okLimits := joint.IsAngularLimitEnabled() &&
		approxEqual32(joint.GetAngularLimitLower(), lower) &&
		approxEqual32(joint.GetAngularLimitUpper(), upper)
	okMotor := joint.IsMotorEnabled() &&
		approxEqual32(joint.GetMotorTargetVelocity(), motorVelocity)
	okSoftness := approxEqual32(joint.GetSoftness(), softness)

	return okPaths && okLimits && okMotor && okSoftness
}

func NewPhysicsValidationFromOwnerObject(owner *GodotObject) GDClass {
	obj := &PhysicsValidation{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassPhysicsValidation() {
	ClassDBRegisterClass(NewPhysicsValidationFromOwnerObject, nil, nil, func(t *PhysicsValidation) {
		ClassDBBindMethod(t, "EnableCCD", "enable_ccd", []string{"body"}, nil)
		ClassDBBindMethod(t, "ApplyFlipperImpulse", "apply_flipper_impulse", []string{"body", "impulse", "position"}, nil)
		ClassDBBindMethod(t, "GetLinearSpeed", "get_linear_speed", []string{"body"}, nil)
		ClassDBBindMethod(t, "ConfigureMaterial", "configure_material", []string{"body", "material", "friction", "bounce"}, nil)
		ClassDBBindMethod(t, "BindArea", "bind_area", []string{"area"}, nil)
		ClassDBBindMethod(t, "ResetAreaCounts", "reset_area_counts", nil, nil)
		ClassDBBindMethod(t, "GetAreaEnterCount", "get_area_enter_count", nil, nil)
		ClassDBBindMethod(t, "GetAreaExitCount", "get_area_exit_count", nil, nil)
		ClassDBBindMethod(t, "OnAreaBodyEntered", "_on_area_body_entered", []string{"body"}, nil)
		ClassDBBindMethod(t, "OnAreaBodyExited", "_on_area_body_exited", []string{"body"}, nil)
		ClassDBBindMethod(t, "ConfigurePinJoint", "configure_pin_joint", []string{"joint", "node_a_path", "node_b_path", "lower", "upper", "softness", "motor_velocity"}, nil)
		log.Debug("PhysicsValidation registered")
	})
}

func UnregisterClassPhysicsValidation() {
	ClassDBUnregisterClass[*PhysicsValidation]()
	log.Debug("PhysicsValidation unregistered")
}

func approxEqual32(a, b float32) bool {
	return float32(math.Abs(float64(a-b))) <= 0.0001
}
