package opengl

import (
	"unsafe"

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
	Dimesion  int32
	created   bool
	ePtr      unsafe.Pointer
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
	vao.buffers[id] = buffer
}

func (vao *VAO) UpdateBuffers() {
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
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(buffer.Elements), buffer.ePtr, gl.DYNAMIC_DRAW)

	//Setup attribute pointer
	attributeId := buffer.vao.shader.EnableAttribute(buffer.attribute)
	gl.VertexAttribPointer(attributeId, buffer.Dimesion, gl.FLOAT, false, 0, nil)
}

func (buffer *Buffer) Update() {
	if !buffer.created {
		buffer.Create()
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.ID)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(buffer.Elements), buffer.ePtr)
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

// Rendering logic
func (vao *VAO) PrepRender() {
	vao.shader.Use()
	vao.PrepUniforms()
	gl.BindVertexArray(vao.id)
	vao.texture.Use()
}

func (vao *VAO) VertNum() int32 {
	for _, b := range vao.buffers {
		return int32(len(b.Elements)) / b.Dimesion
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
