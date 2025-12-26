package pkg

import (
	"math"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/ffi"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	"github.com/godot-go/godot-go/pkg/log"
)

type PhysicsBenchmark struct {
	Node2DImpl
	bodies []RigidBody2D
	floor  StaticBody2D
}

func (p *PhysicsBenchmark) GetClassName() string {
	return "PhysicsBenchmark"
}

func (p *PhysicsBenchmark) GetParentClassName() string {
	return "Node2D"
}

func (p *PhysicsBenchmark) Setup(count int32, radius, spacing float32) bool {
	p.ClearBodies()
	if count <= 0 {
		return false
	}
	if radius <= 0 {
		radius = 4
	}
	if spacing <= 0 {
		spacing = radius * 3
	}
	cols := int32(math.Ceil(math.Sqrt(float64(count))))
	if cols < 1 {
		cols = 1
	}
	width := float32(cols) * spacing
	rows := int32(math.Ceil(float64(count) / float64(cols)))
	floorY := float32(rows+2) * spacing
	p.ensureFloor(width, spacing, floorY)
	for i := int32(0); i < count; i++ {
		body := p.newRigidBody2D()
		if body == nil {
			return false
		}
		body.SetCanSleep(false)
		body.SetGravityScale(1)
		x := float32(i%cols) * spacing
		y := float32(i/cols) * spacing
		body.SetPosition(NewVector2WithFloat32Float32(x, y))
		shapeNode := p.newCollisionShape2D()
		if shapeNode == nil {
			return false
		}
		circle := p.newCircleShape2D(radius)
		if circle == nil {
			return false
		}
		shapeRef := NewRefShape2DGDExtensionIternalConstructor(circle.TypedPtr())
		shapeNode.SetShape(shapeRef)
		body.AddChild(shapeNode, false, NODE_INTERNAL_MODE_INTERNAL_MODE_DISABLED)
		p.AddChild(body, false, NODE_INTERNAL_MODE_INTERNAL_MODE_DISABLED)
		p.bodies = append(p.bodies, body)
	}
	return true
}

func (p *PhysicsBenchmark) ApplyImpulseBatch(impulse Vector2) int32 {
	var applied int32
	for _, body := range p.bodies {
		if body == nil {
			continue
		}
		body.ApplyImpulse(impulse, Vector2{})
		applied++
	}
	return applied
}

func (p *PhysicsBenchmark) ClearBodies() {
	for _, body := range p.bodies {
		if body == nil {
			continue
		}
		p.RemoveChild(body)
		body.QueueFree()
	}
	p.bodies = nil
}

func (p *PhysicsBenchmark) GetBodyCount() int32 {
	return int32(len(p.bodies))
}

func (p *PhysicsBenchmark) ensureFloor(width, height, y float32) {
	if p.floor != nil {
		return
	}
	body := p.newStaticBody2D()
	if body == nil {
		return
	}
	shapeNode := p.newCollisionShape2D()
	if shapeNode == nil {
		return
	}
	rect := p.newRectangleShape2D(width, height)
	if rect == nil {
		return
	}
	shapeRef := NewRefShape2DGDExtensionIternalConstructor(rect.TypedPtr())
	shapeNode.SetShape(shapeRef)
	body.AddChild(shapeNode, false, NODE_INTERNAL_MODE_INTERNAL_MODE_DISABLED)
	body.SetPosition(NewVector2WithFloat32Float32(width/2, y))
	p.AddChild(body, false, NODE_INTERNAL_MODE_INTERNAL_MODE_DISABLED)
	p.floor = body
}

func (p *PhysicsBenchmark) newRigidBody2D() RigidBody2D {
	owner := constructObject("RigidBody2D")
	if owner == nil {
		return nil
	}
	return NewRigidBody2DWithGodotOwnerObject(owner)
}

func (p *PhysicsBenchmark) newStaticBody2D() StaticBody2D {
	owner := constructObject("StaticBody2D")
	if owner == nil {
		return nil
	}
	return NewStaticBody2DWithGodotOwnerObject(owner)
}

func (p *PhysicsBenchmark) newCollisionShape2D() CollisionShape2D {
	owner := constructObject("CollisionShape2D")
	if owner == nil {
		return nil
	}
	return NewCollisionShape2DWithGodotOwnerObject(owner)
}

func (p *PhysicsBenchmark) newCircleShape2D(radius float32) RefCircleShape2D {
	owner := constructObject("CircleShape2D")
	if owner == nil {
		return nil
	}
	shape := NewCircleShape2DWithGodotOwnerObject(owner)
	if shape != nil {
		shape.TypedPtr().SetRadius(radius)
	}
	return shape
}

func (p *PhysicsBenchmark) newRectangleShape2D(width, height float32) RefRectangleShape2D {
	owner := constructObject("RectangleShape2D")
	if owner == nil {
		return nil
	}
	shape := NewRectangleShape2DWithGodotOwnerObject(owner)
	if shape != nil {
		shape.TypedPtr().SetSize(NewVector2WithFloat32Float32(width, height))
	}
	return shape
}

func constructObject(className string) *GodotObject {
	sn := NewStringNameWithLatin1Chars(className)
	defer sn.Destroy()
	owner := CallFunc_GDExtensionInterfaceClassdbConstructObject(sn.AsGDExtensionConstStringNamePtr())
	return (*GodotObject)(owner)
}

func NewPhysicsBenchmarkFromOwnerObject(owner *GodotObject) GDClass {
	obj := &PhysicsBenchmark{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassPhysicsBenchmark() {
	ClassDBRegisterClass(NewPhysicsBenchmarkFromOwnerObject, nil, nil, func(t *PhysicsBenchmark) {
		ClassDBBindMethod(t, "Setup", "setup", []string{"count", "radius", "spacing"}, nil)
		ClassDBBindMethod(t, "ApplyImpulseBatch", "apply_impulse_batch", []string{"impulse"}, nil)
		ClassDBBindMethod(t, "ClearBodies", "clear_bodies", nil, nil)
		ClassDBBindMethod(t, "GetBodyCount", "get_body_count", nil, nil)
		log.Debug("PhysicsBenchmark registered")
	})
}

func UnregisterClassPhysicsBenchmark() {
	ClassDBUnregisterClass[*PhysicsBenchmark]()
	log.Debug("PhysicsBenchmark unregistered")
}
