package opengl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type DefaultVAO struct {
	*BaseVAO
	cam, position     *mgl32.Vec2
	bounds            mgl32.Vec4
	checkBounds       bool
	position_pointers []*float32
	pointers_updated  bool
}

func CreateDefaultVao(window *Window, textureSource string, elements int) *DefaultVAO {
	vao := CreateVAO(window, textureSource)

	vBuff := Buffer{
		Dimension: 2,
	}
	tBuff := Buffer{
		Dimension: 2,
	}

	vElements := make([]float32, elements*2*3) // 2 points per vertex, 3 per triangle
	tElements := make([]float32, elements*2*3)

	vBuff.Elements = vElements
	tBuff.Elements = tElements

	vao.AddBuffer("vert", &vBuff)
	vao.AddBuffer("verttexcoord", &tBuff)

	var x, y float32
	defaultVAO := DefaultVAO{vao, &mgl32.Vec2{}, &mgl32.Vec2{}, mgl32.Vec4{}, false, []*float32{&x, &y}, false}

	defaultVAO.AttachDefaultShader()
	defaultVAO.Init()

	return &defaultVAO
}

/*
Pointers are set beforehand and updated, theadsafe? No. Works? Yes.
For actual stuff that requires thread safe operators (eg not just setting rotations) we can use the jobBlocks.
*/

// Called from non opengl thread
func (vao *DefaultVAO) SetTranslation(x, y *float32) {
	vao.position_pointers[0] = x
	vao.position_pointers[1] = y
	vao.UpdateUniforms()
}

func (vao *DefaultVAO) UpdatePointers() {
	vao.pointers_updated = true
}

// Called in the opengl thread
func (vao *DefaultVAO) updatePointers() {
	vao.position[0] = *vao.position_pointers[0]
	vao.position[1] = *vao.position_pointers[1]
	vao.UpdateUniforms()
	vao.pointers_updated = false
}

// Rendering logic
func (vao *DefaultVAO) PrepRender() {
	if vao.pointers_updated {
		vao.updatePointers()
	}
	vao.BaseVAO.PrepRender()
}

func (vao *DefaultVAO) AttachDefaultShader() {
	program := CreateProgram(0)
	vao.AttachShader(program)

	program.LoadVertShader("./resources/shaders/vertex.vert")
	program.LoadFragShader("./resources/shaders/fragment.frag")
	program.Link()

	program.AddAttribute("vert")
	// Currently unusued, optimized out by the shader compiler so will fail
	// program.AddAttribute("rotgroup")
	program.AddAttribute("verttexcoord")

	// Add and set rotation uniform
	vao.AddUniform("rot", &mgl32.Mat2{1, 0, 0, 1})
	vao.AddUniform("rotcenter", &mgl32.Vec2{0, 0})

	// Other uniforms can use default values.
	var zoom float32 = 1
	vao.AddUniform("trans", vao.position)
	vao.AddUniform("dim", &mgl32.Mat2{2. / float32(vao.window.Width), 0., 0., 2. / float32(vao.window.Height)})
	vao.AddUniform("cam", &mgl32.Vec2{})
	vao.AddUniform("zoom", &zoom)
}

// CPU side culling

func (vao *DefaultVAO) SetRenderBounds(bounds mgl32.Vec4) {
	vao.checkBounds = true
	vao.bounds = bounds
}

func (vao *DefaultVAO) ShouldRender() bool {
	if !vao.checkBounds {
		return true
	}

	// Check if we are within bounds
	windowWidth := vao.window.Width
	windowHeight := vao.window.Height

	bPos := mgl32.Vec2{vao.bounds.X(), vao.bounds.Y()}.Add(*vao.position)
	bPosBoundary := mgl32.Vec2{vao.bounds.Z(), vao.bounds.W()}.Add(bPos)

	if bPos.X() > float32(windowWidth) && bPosBoundary.X() > float32(windowWidth) {
		return false
	}

	if bPos.X() < 0 && bPosBoundary.X() < 0 {
		return false
	}

	if bPos.Y() > float32(windowHeight) && bPosBoundary.Y() > float32(windowHeight) {
		return false
	}

	if bPos.Y() < 0 && bPosBoundary.Y() < 0 {
		return false
	}

	return true
}
