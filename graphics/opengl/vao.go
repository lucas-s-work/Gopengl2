package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VAO struct {
	id       uint32
	buffers  map[string]*Buffer
	uniforms map[string]interface{}
	window   *Window
	shader   *Program
	texture  *Texture
}

type Buffer struct {
	ID        uint32
	Elements  []float32
	Dimension int32
	created   bool
	attribute string
	vao       *VAO
}

// VAO creation and destruction
func CreateVAO(window *Window, textureSource string) *VAO {
	var vaoID uint32

	gl.GenVertexArrays(1, &vaoID)

	texture := LoadTexture(textureSource)

	vao := VAO{
		id:       vaoID,
		window:   window,
		texture:  texture,
		buffers:  make(map[string]*Buffer),
		uniforms: make(map[string]interface{}),
	}

	return &vao
}

func (vao *VAO) Init() {
	vao.BindVao()

	for _, b := range vao.buffers {
		b.Create()
	}
}

func (vao *VAO) Delete() {
	for _, b := range vao.buffers {
		b.Delete()
	}

	gl.DeleteVertexArrays(1, &vao.id)
}

func (vao *VAO) BindVao() {
	gl.BindVertexArray(vao.id)
}

func (vao *VAO) AddBuffer(id string, buffer *Buffer) {
	buffer.vao = vao
	buffer.attribute = id
	vao.buffers[id] = buffer
}

func (vao *VAO) GetBuffer(id string) *Buffer {
	return vao.buffers[id]
}

func (vao *VAO) UpdateBuffers() {
	vao.BindVao()

	for _, b := range vao.buffers {
		b.Update()
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vao *VAO) UpdateBuffer(name string) {
	vao.buffers[name].Update()
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

// Buffer Updating

func (buffer *Buffer) Create() {
	if buffer.created {
		panic("Attempting to re-create created Buffer")
	}

	// Generate buffer
	buffer.created = true
	gl.GenBuffers(1, &buffer.ID)

	// Set buffer data
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.ID)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(buffer.Elements), gl.Ptr(buffer.Elements), gl.DYNAMIC_DRAW)

	//Setup attribute pointer
	attributeId := buffer.vao.shader.EnableAttribute(buffer.attribute)
	gl.VertexAttribPointer(attributeId, buffer.Dimension, gl.FLOAT, false, 0, nil)
}

func (buffer *Buffer) Update() {
	if !buffer.created {
		buffer.Create()
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.ID)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(buffer.Elements), gl.Ptr(buffer.Elements))
}

func (buffer *Buffer) Delete() {
	gl.DeleteBuffers(1, &buffer.ID)
	buffer.created = false
}

// Shader Logic
func (vao *VAO) AttachShader(shader *Program) {
	vao.shader = shader
}

func (vao *VAO) AttachProgram(program *Program) {
	vao.shader = program
}

func (vao *VAO) AddUniform(name string, value interface{}) {
	vao.shader.AddUniform(name, value)

	vao.uniforms[name] = value
}

func (vao *VAO) PrepUniforms() {
	for id, uni := range vao.shader.uniforms {
		vao.shader.SetUniform(id, uni.Value())
	}
}

func (vao *VAO) PixToTex(texX, texY int) (float32, float32) {
	return vao.texture.PixToTex(texX, texY)
}

// Rendering logic
func (vao *VAO) PrepRender() {
	vao.shader.Use()
	vao.PrepUniforms()
	vao.BindVao()
	vao.texture.Use()
}

func (vao *VAO) VertNum() int32 {
	for _, b := range vao.buffers {
		a := int32(len(b.Elements)) / b.Dimension
		return a
	}

	return 0
}

func (vao *VAO) Render() {
	vao.PrepRender()
	gl.DrawArrays(gl.TRIANGLES, 0, vao.VertNum())
}

func RenderVaos(vaos []*VAO) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, vao := range vaos {
		vao.Render()
	}
}
