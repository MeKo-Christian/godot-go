package pkg

import (
	"math"
	"unsafe"

	. "github.com/godot-go/godot-go/pkg/builtin"
	. "github.com/godot-go/godot-go/pkg/constant"
	. "github.com/godot-go/godot-go/pkg/core"
	. "github.com/godot-go/godot-go/pkg/gdclassimpl"
	. "github.com/godot-go/godot-go/pkg/gdutilfunc"
)

// AudioVisualDemo implements GDClass evidence.
var _ GDClass = (*AudioVisualDemo)(nil)

type AudioVisualDemo struct {
	Node2DImpl
	time         float64
	nextPlayAt   float64
	frameTimer   float64
	frameIndex   int32
	sfxBaseX     float32
	sfxBaseY     float32
	sfxPlayer    AudioStreamPlayer2D
	beepStream   RefAudioStream
	sprite       AnimatedSprite2D
	spriteFrames RefSpriteFrames
	particles    GPUParticles2D
	particleMat  RefMaterial
	particleTex  RefTexture2D
	shaderMat    RefShaderMaterial
}

func (d *AudioVisualDemo) GetClassName() string {
	return "AudioVisualDemo"
}

func (d *AudioVisualDemo) GetParentClassName() string {
	return "Node2D"
}

func (d *AudioVisualDemo) V_Ready() {
	d.SetProcess(true)
	d.setupAudio()
	d.setupParticles()
	d.setupAnimatedSprite()
	d.setupShaderPanel()
	printLine("Audio/visual demo ready: listen for the beep and watch the effects.")
}

func (d *AudioVisualDemo) V_Process(delta float64) {
	d.time += delta
	d.updateAudioMotion()
	d.updateAnimatedSprite(delta)
}

func (d *AudioVisualDemo) setupAudio() {
	node := d.GetNodeOrNull(nodePath("SFXPlayer"))
	if node == nil {
		printLine("AudioVisualDemo: AudioStreamPlayer2D not found")
		return
	}
	player, ok := ObjectCastTo(node, "AudioStreamPlayer2D").(AudioStreamPlayer2D)
	if !ok || player == nil {
		printLine("AudioVisualDemo: AudioStreamPlayer2D cast failed")
		return
	}
	d.sfxPlayer = player

	stream := d.loadAudioStream("res://sfx/beep.wav")
	if stream == nil {
		printLine("AudioVisualDemo: failed to load res://sfx/beep.wav")
		return
	}
	d.beepStream = stream
	player.SetStream(stream)
	player.SetVolumeDb(-6)
	player.SetMaxDistance(520)
	player.SetAttenuation(2.2)
	player.SetPanningStrength(1.0)
	player.Play(0)

	pos := player.GetPosition()
	d.sfxBaseX = pos.MemberGetx()
	d.sfxBaseY = pos.MemberGety()
}

func (d *AudioVisualDemo) updateAudioMotion() {
	if d.sfxPlayer == nil {
		return
	}
	offset := float32(math.Sin(d.time*1.2)) * 160
	pos := NewVector2WithFloat32Float32(d.sfxBaseX+offset, d.sfxBaseY)
	d.sfxPlayer.SetPosition(pos)

	if !d.sfxPlayer.IsPlaying() && d.time >= d.nextPlayAt {
		d.sfxPlayer.Play(0)
		d.nextPlayAt = d.time + 0.85
	}
}

func (d *AudioVisualDemo) setupParticles() {
	node := d.GetNodeOrNull(nodePath("ImpactParticles"))
	if node == nil {
		printLine("AudioVisualDemo: GPUParticles2D not found")
		return
	}
	particles, ok := ObjectCastTo(node, "GPUParticles2D").(GPUParticles2D)
	if !ok || particles == nil {
		printLine("AudioVisualDemo: GPUParticles2D cast failed")
		return
	}
	d.particles = particles

	mat := d.newParticleProcessMaterial()
	if mat == nil {
		printLine("AudioVisualDemo: failed to create particle material")
		return
	}
	texture := createSolidTexture(12, NewColorWithFloat32Float32Float32Float32(1, 0.75, 0.3, 1))
	if texture != nil {
		d.particleTex = texture
		particles.SetTexture(texture)
	}
	d.particleMat = NewRefMaterial(mat)
	particles.SetProcessMaterial(d.particleMat)
	particles.SetAmount(72)
	particles.SetLifetime(0.9)
	particles.SetEmitting(true)
}

func (d *AudioVisualDemo) setupAnimatedSprite() {
	node := d.GetNodeOrNull(nodePath("GlowSprite"))
	if node == nil {
		printLine("AudioVisualDemo: AnimatedSprite2D not found")
		return
	}
	sprite, ok := ObjectCastTo(node, "AnimatedSprite2D").(AnimatedSprite2D)
	if !ok || sprite == nil {
		printLine("AudioVisualDemo: AnimatedSprite2D cast failed")
		return
	}
	d.sprite = sprite

	frames := d.createSpriteFrames()
	if frames == nil {
		printLine("AudioVisualDemo: failed to create sprite frames")
		return
	}
	d.spriteFrames = frames
	sprite.SetSpriteFrames(frames)
	anim := NewStringNameWithLatin1Chars("pulse")
	defer anim.Destroy()
	sprite.SetAnimation(anim)
	sprite.SetPlaying(false)
	sprite.SetFrame(0)
}

func (d *AudioVisualDemo) updateAnimatedSprite(delta float64) {
	if d.sprite == nil {
		return
	}
	d.frameTimer += delta
	if d.frameTimer < 0.18 {
		return
	}
	d.frameTimer = 0
	d.frameIndex = (d.frameIndex + 1) % 2
	d.sprite.SetFrame(d.frameIndex)
}

func (d *AudioVisualDemo) setupShaderPanel() {
	node := d.GetNodeOrNull(nodePath("ShaderPanel"))
	if node == nil {
		printLine("AudioVisualDemo: Shader panel not found")
		return
	}
	panel, ok := ObjectCastTo(node, "ColorRect").(ColorRect)
	if !ok || panel == nil {
		printLine("AudioVisualDemo: ColorRect cast failed")
		return
	}

	shader := d.newShader(shaderCode())
	if shader == nil {
		printLine("AudioVisualDemo: failed to create shader")
		return
	}
	shaderMat := d.newShaderMaterial(shader)
	if shaderMat == nil {
		printLine("AudioVisualDemo: failed to create shader material")
		return
	}
	d.shaderMat = NewRefShaderMaterial(shaderMat)
	panel.SetMaterial(NewRefMaterial(shaderMat))
}

func (d *AudioVisualDemo) createSpriteFrames() RefSpriteFrames {
	framesObj := instantiateObject("SpriteFrames")
	frames, ok := ObjectCastTo(framesObj, "SpriteFrames").(SpriteFrames)
	if !ok || frames == nil {
		return nil
	}
	anim := NewStringNameWithLatin1Chars("pulse")
	defer anim.Destroy()
	frames.AddAnimation(anim)
	frames.SetAnimationLoop(anim, true)
	frames.SetAnimationSpeed(anim, 6)

	colorA := NewColorWithFloat32Float32Float32Float32(0.2, 0.9, 1.0, 1)
	colorB := NewColorWithFloat32Float32Float32Float32(1.0, 0.2, 0.8, 1)
	texA := createSolidTexture(24, colorA)
	texB := createSolidTexture(24, colorB)
	if texA == nil || texB == nil {
		return nil
	}
	frames.AddFrame(anim, texA, 0.0, -1)
	frames.AddFrame(anim, texB, 0.0, -1)
	return NewRefSpriteFrames(frames)
}

func (d *AudioVisualDemo) newParticleProcessMaterial() ParticleProcessMaterial {
	matObj := instantiateObject("ParticleProcessMaterial")
	mat, ok := ObjectCastTo(matObj, "ParticleProcessMaterial").(ParticleProcessMaterial)
	if !ok || mat == nil {
		return nil
	}
	mat.SetDirection(NewVector3WithFloat32Float32Float32(0, -1, 0))
	mat.SetSpread(30)
	mat.SetGravity(NewVector3WithFloat32Float32Float32(0, 220, 0))
	return mat
}

func (d *AudioVisualDemo) newShader(code string) Shader {
	shaderObj := instantiateObject("Shader")
	shader, ok := ObjectCastTo(shaderObj, "Shader").(Shader)
	if !ok || shader == nil {
		return nil
	}
	codeStr := NewStringWithUtf8Chars(code)
	defer codeStr.Destroy()
	shader.SetCode(codeStr)
	return shader
}

func (d *AudioVisualDemo) newShaderMaterial(shader Shader) ShaderMaterial {
	matObj := instantiateObject("ShaderMaterial")
	mat, ok := ObjectCastTo(matObj, "ShaderMaterial").(ShaderMaterial)
	if !ok || mat == nil {
		return nil
	}
	mat.SetShader(NewRefShader(shader))
	return mat
}

func (d *AudioVisualDemo) loadAudioStream(path string) RefAudioStream {
	loader := getResourceLoaderSingleton()
	if loader == nil {
		return nil
	}
	pathStr := NewStringWithUtf8Chars(path)
	defer pathStr.Destroy()
	typeHint := NewStringWithUtf8Chars("AudioStream")
	defer typeHint.Destroy()
	resource := loader.Load(pathStr, typeHint, RESOURCE_LOADER_CACHE_MODE_CACHE_MODE_REUSE)
	if resource == nil {
		return nil
	}
	resObj := resource.TypedPtr()
	stream, ok := ObjectCastTo(resObj, "AudioStream").(AudioStream)
	if !ok || stream == nil {
		return nil
	}
	return NewRefAudioStream(stream)
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

func getResourceLoaderSingleton() ResourceLoader {
	owner := (*GodotObject)(unsafe.Pointer(GetSingleton("ResourceLoader")))
	if owner == nil {
		printLine("AudioVisualDemo: ResourceLoader singleton not found")
		return nil
	}
	return NewResourceLoaderWithGodotOwnerObject(owner)
}

func getClassDBSingleton() ClassDB {
	owner := (*GodotObject)(unsafe.Pointer(GetSingleton("ClassDB")))
	if owner == nil {
		printLine("AudioVisualDemo: ClassDB singleton not found")
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

func shaderCode() string {
	return `shader_type canvas_item;

void fragment() {
	vec2 uv = UV * 2.0 - 1.0;
	float glow = 0.6 + 0.4 * sin(TIME * 3.5 + length(uv) * 6.0);
	vec3 inner = vec3(0.2, 0.35, 0.9);
	vec3 outer = vec3(1.0, 0.65, 0.2);
	COLOR = vec4(mix(inner, outer, glow), 1.0);
}
`
}

func NewAudioVisualDemoFromOwnerObject(owner *GodotObject) GDClass {
	obj := &AudioVisualDemo{}
	obj.SetGodotObjectOwner(owner)
	return obj
}

func RegisterClassAudioVisualDemo() {
	ClassDBRegisterClass(NewAudioVisualDemoFromOwnerObject, nil, nil, func(t *AudioVisualDemo) {
		ClassDBBindMethodVirtual(t, "V_Ready", "_ready", nil, nil)
		ClassDBBindMethodVirtual(t, "V_Process", "_process", []string{"delta"}, nil)
	})
}

func UnregisterClassAudioVisualDemo() {
	ClassDBUnregisterClass[*AudioVisualDemo]()
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
