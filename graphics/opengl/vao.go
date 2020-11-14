package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VAO interface {
	Init()
	Delete()
	BindVao()
	AddBuffer(string, *Buffer)
	GetBuffer(string) *Buffer
	UpdateBuffers()
	UpdateBuffer(string)
	AttachShader(*Program)
	GetShader() *Program
	AddUniform(string, interface{})
	PixToTex(int, int) (float32, float32)
	PrepRender()
	VertNum() int32
	Render()
}

type BaseVAO struct {
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
	vao       VAO
}

// VAO creation and destruction
func CreateVAO(window *Window, textureSource string) *BaseVAO {
	var vaoID uint32

	gl.GenVertexArrays(1, &vaoID)

	texture := LoadTexture(textureSource)

	vao := BaseVAO{
		id:       vaoID,
		window:   window,
		texture:  texture,
		buffers:  make(map[string]*Buffer),
		uniforms: make(map[string]interface{}),
	}

	return &vao
}

func (vao *BaseVAO) Init() {
	vao.BindVao()

	for _, b := range vao.buffers {
		b.Create()
	}
}

func (vao *BaseVAO) Delete() {
	for _, b := range vao.buffers {
		b.Delete()
	}

	gl.DeleteVertexArrays(1, &vao.id)
}

func (vao *BaseVAO) BindVao() {
	gl.BindVertexArray(vao.id)
}

func (vao *BaseVAO) AddBuffer(id string, buffer *Buffer) {
	buffer.vao = vao
	buffer.attribute = id
	vao.buffers[id] = buffer
}

func (vao *BaseVAO) GetBuffer(id string) *Buffer {
	return vao.buffers[id]
}

func (vao *BaseVAO) UpdateBuffers() {
	vao.BindVao()

	for _, b := range vao.buffers {
		b.Update()
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (vao *BaseVAO) UpdateBuffer(name string) {
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
	attributeId := buffer.vao.GetShader().EnableAttribute(buffer.attribute)
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
func (vao *BaseVAO) AttachShader(shader *Program) {
	vao.shader = shader
}

// Uniforms should be pointers to the variable that is being set
func (vao *BaseVAO) AddUniform(name string, value interface{}) {
	vao.shader.AddUniform(name, value)

	vao.uniforms[name] = value
}

func (vao *BaseVAO) GetShader() *Program {
	return vao.shader
}

// should be moved into the shader
func (vao *BaseVAO) UpdateUniforms() {
	vao.shader.UpdateUniforms()
}

func (vao *BaseVAO) PixToTex(texX, texY int) (float32, float32) {
	return vao.texture.PixToTex(texX, texY)
}

// Rendering logic
func (vao *BaseVAO) PrepRender() {
	vao.shader.Use()
	vao.BindVao()
	vao.texture.Use()
}

func (vao *BaseVAO) VertNum() int32 {
	for _, b := range vao.buffers {
		a := int32(len(b.Elements)) / b.Dimension
		return a
	}

	return 0
}

func (vao *BaseVAO) Render() {
	vao.PrepRender()
	gl.DrawArrays(gl.TRIANGLES, 0, vao.VertNum())
}

func RenderVaos(vaos []VAO) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, vao := range vaos {
		vao.Render()
	}
}
