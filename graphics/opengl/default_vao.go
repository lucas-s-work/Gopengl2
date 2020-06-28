package opengl

import "github.com/go-gl/mathgl/mgl32"

func CreateDefaultVao(window *Window, textureSource string, elements int) *VAO {
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

	vao.AttachDefaultShader()
	vao.Init()

	return vao
}

func (vao *VAO) AttachDefaultShader() {
	program := CreateProgram(0)
	vao.AttachShader(program)

	program.LoadVertShader("Gopengl2/shaders/vertex.vert")
	program.LoadFragShader("Gopengl2/shaders/fragment.frag")
	program.Link()

	program.AddAttribute("vert")
	// Currently unusued, optimized out by the shader compiler so will fail
	// program.AddAttribute("rotgroup")
	program.AddAttribute("verttexcoord")

	// Add and set rotation uniform
	vao.AddUniform("rot", mgl32.Vec4{0, 0, 1, 0})

	// Other uniforms can use default values.
	var zoom float32 = 1

	vao.AddUniform("trans", mgl32.Vec2{})
	vao.AddUniform("dim", mgl32.Vec2{float32(vao.window.Width), float32(vao.window.Height)})
	vao.AddUniform("cam", mgl32.Vec2{})
	vao.AddUniform("zoom", zoom)
}
