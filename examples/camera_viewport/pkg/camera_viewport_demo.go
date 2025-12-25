package pkg

import (
	"fmt"
	"math"
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// CameraViewportDemo implements GDClass evidence.
var _ GDClass = (*CameraViewportDemo)(nil)

type CameraViewportDemo struct {
	Node2DImpl
	time       float64
	target     Node2D
	label      Label
	targetBase Vector2
}

func (d *CameraViewportDemo) GetClassName() string {
	return "CameraViewportDemo"
}

func (d *CameraViewportDemo) GetParentClassName() string {
	return "Node2D"
}

func (d *CameraViewportDemo) V_Ready() {
	d.SetProcess(true)
	d.setupTarget()
	d.setupCamera()
	d.setupLabel()
	d.connectViewportSignal()
	d.updateViewportLabel()
	printLine("Camera/viewport demo ready: the camera follows the moving target.")
}

func (d *CameraViewportDemo) V_Process(delta float64) {
	d.time += delta
	d.updateTarget()
}

func (d *CameraViewportDemo) OnViewportSizeChanged() {
	d.updateViewportLabel()
}

func (d *CameraViewportDemo) setupTarget() {
	node := d.GetNodeOrNull(nodePath("Target"))
	if node == nil {
		printLine("CameraViewportDemo: Target not found")
		return
	}
	target, ok := ObjectCastTo(node, "Node2D").(Node2D)
	if !ok || target == nil {
		printLine("CameraViewportDemo: Target cast failed")
		return
	}
	d.target = target
	d.targetBase = target.GetPosition()

	spriteNode := target.GetNodeOrNull(nodePath("TargetSprite"))
	if spriteNode == nil {
		return
	}
	sprite, ok := ObjectCastTo(spriteNode, "Sprite2D").(Sprite2D)
	if !ok || sprite == nil {
		return
	}
	texture := createSolidTexture(32, NewColorWithFloat32Float32Float32Float32(0.2, 0.9, 0.4, 1))
	if texture != nil {
		sprite.SetTexture(texture)
	}
}

func (d *CameraViewportDemo) setupCamera() {
	node := d.GetNodeOrNull(nodePath("Target/Camera2D"))
	if node == nil {
		return
	}
	camera, ok := ObjectCastTo(node, "Camera2D").(Camera2D)
	if !ok || camera == nil {
		return
	}
	camera.SetPositionSmoothingEnabled(true)
	camera.SetPositionSmoothingSpeed(6)
	camera.SetZoom(NewVector2WithFloat32Float32(1.1, 1.1))
	camera.MakeCurrent()
}

func (d *CameraViewportDemo) setupLabel() {
	node := d.GetNodeOrNull(nodePath("ViewportLabel"))
	if node == nil {
		return
	}
	label, ok := ObjectCastTo(node, "Label").(Label)
	if !ok || label == nil {
		return
	}
	d.label = label
}

func (d *CameraViewportDemo) updateTarget() {
	if d.target == nil {
		return
	}
	x := d.targetBase.MemberGetx() + float32(math.Sin(d.time*0.8))*260
	y := d.targetBase.MemberGety() + float32(math.Cos(d.time*1.1))*140
	pos := NewVector2WithFloat32Float32(x, y)
	d.target.SetPosition(pos)
}

func (d *CameraViewportDemo) updateViewportLabel() {
	if d.label == nil {
		return
	}
	viewport := d.GetViewport()
	if viewport == nil {
		return
	}
	rect := viewport.GetVisibleRect()
	size := rect.MemberGetsize()
	text := NewStringWithUtf8Chars(fmt.Sprintf("Viewport: %.0fx%.0f", size.MemberGetx(), size.MemberGety()))
	defer text.Destroy()
	d.label.SetText(text)
}

func (d *CameraViewportDemo) connectViewportSignal() {
	viewport := d.GetViewport()
	if viewport == nil {
		return
	}
	signal := NewStringNameWithLatin1Chars("size_changed")
	defer signal.Destroy()
	method := NewStringNameWithLatin1Chars("_on_viewport_size_changed")
	defer method.Destroy()
	callable := NewCallableWithObjectStringName(d, method)
	defer callable.Destroy()
	viewport.Connect(signal, callable, 0)
}

func NewCameraViewportDemoFromOwnerObject(owner *GodotObject) GDClass {
	obj := &CameraViewportDemo{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassCameraViewportDemo() {
	ClassDBRegisterClass(NewCameraViewportDemoFromOwnerObject, nil, nil, func(t *CameraViewportDemo) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Process", "_process", []string{"delta"}, nil)
		ClassDBBindMethod(t, "OnViewportSizeChanged", "_on_viewport_size_changed", nil, nil)
	})
}

func UnregisterClassCameraViewportDemo() {
	ClassDBUnregisterClass[*CameraViewportDemo]()
}

func instantiateObject(className string) Object {
	classDB := getClassDBSingleton()
	if classDB == nil {
		return nil
	}
	name := NewStringNameWithLatin1Chars(className)
	defer name.Destroy()
	v := classDB.Instantiate(name)
	defer v.Destroy()
	return v.ToObject()
}

func getClassDBSingleton() ClassDB {
	owner := (*GodotObject)(unsafe.Pointer(GetSingleton("ClassDB")))
	if owner == nil {
		printLine("CameraViewportDemo: ClassDB singleton not found")
		return nil
	}
	return NewClassDBWithGodotOwnerObject(owner)
}

func createSolidTexture(size int32, color Color) RefTexture2D {
	imageFactoryObj := instantiateObject("Image")
	imageFactory, ok := ObjectCastTo(imageFactoryObj, "Image").(Image)
	if !ok || imageFactory == nil {
		return nil
	}
	image := imageFactory.Create(size, size, false, IMAGE_FORMAT_FORMAT_RGBA_8)
	if image == nil {
		return nil
	}
	image.TypedPtr().Fill(color)

	textureFactoryObj := instantiateObject("ImageTexture")
	textureFactory, ok := ObjectCastTo(textureFactoryObj, "ImageTexture").(ImageTexture)
	if !ok || textureFactory == nil {
		return nil
	}
	texture := textureFactory.CreateFromImage(image)
	if texture == nil {
		return nil
	}
	return NewRefTexture2D(texture.TypedPtr())
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
