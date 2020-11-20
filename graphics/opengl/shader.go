package opengl

import (
	"Gopengl2/util"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader types

const (
	INVALID    = 0
	VERTSHADER = gl.VERTEX_SHADER
	FRAGSHADER = gl.FRAGMENT_SHADER
	GEOMSHADER = gl.GEOMETRY_SHADER
)

var (
	loadedShaders []*shader
)

type shader struct {
	Id   uint32
	file string
}

type Program struct {
	Id         uint32
	attributes map[string]uint32
	uniforms   map[string]uniform
}

// Shader program loading and creation

func CreateProgram(Id uint32) *Program {
	if Id == 0 {
		Id = gl.CreateProgram()
	}

	return &Program{
		Id,
		make(map[string]uint32),
		make(map[string]uniform),
	}
}

func (program *Program) AttachShader(s *shader) {
	gl.AttachShader(program.Id, s.Id)
}

/*
Load and attach shaders, if the shader has already been loaded it is not re-created.
*/

func ReadFile(source string) (string, error) {
	data, err := ioutil.ReadFile(util.RelativePath(source))

	if err != nil {
		return "", err
	}

	return string(data[:]) + "\x00", nil
}

func (program *Program) LoadVertShader(file string) {
	program.loadShader(file, VERTSHADER)
}

func (program *Program) LoadFragShader(file string) {
	program.loadShader(file, FRAGSHADER)
}

func (program *Program) loadShader(file string, shaderType uint32) {
	existingShader := findShader(file)

	if existingShader != nil {
		program.AttachShader(existingShader)

		return
	}

	rawData, err := ReadFile(file)

	if err != nil {
		panic(fmt.Errorf("Unable to find vertex shader file: %s", file))
	}

	program.loadShader(rawData, FRAGSHADER)

	shaderId := gl.CreateShader(shaderType)
	source, free := gl.Strs(rawData)

	gl.ShaderSource(shaderId, 1, source, nil)
	free()
	gl.CompileShader(shaderId)

	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("failed to compile %v: %v", source, log))
	}

	loadedShaders = append(loadedShaders, &shader{
		shaderId, file,
	})

	gl.AttachShader(program.Id, shaderId)
}

// Determine if a shader has already been created

func findShader(file string) *shader {
	for _, s := range loadedShaders {
		if s.file == file {
			return s
		}
	}

	return nil
}

// Shader binding and linking

func (p *Program) Use() {
	gl.UseProgram(p.Id)
}

func (p *Program) UnUse() {
	gl.UseProgram(0)
}

func (p *Program) Link() {
	gl.LinkProgram(p.Id)
}

// Attribute handling

func (p *Program) AddAttribute(attribute string) {
	attrib := gl.GetAttribLocation(p.Id, gl.Str(attribute+"\x00"))

	if attrib == -1 {
		panic("Invalid Attribute given")
	}

	p.attributes[attribute] = uint32(attrib)
}

func (p *Program) EnableAttribute(attribute string) uint32 {
	attributeValue := p.attributes[attribute]

	gl.EnableVertexAttribArray(attributeValue)

	return attributeValue
}

func (p *Program) DisableAttribute(attribute string) {
	gl.DisableVertexAttribArray(p.attributes[attribute])
}

// Uniform handling

type uniform struct {
	id      uint32
	value   interface{}
	updated bool
}

func (uni *uniform) ID() uint32 {
	return uni.id
}

func (uni *uniform) Value() interface{} {
	return uni.value
}

/*
Each Render call Attach is called for each uniform, Attach only sets the value if
Update set
*/

func (uni *uniform) Update() {
	uni.updated = true
}

func (uni *uniform) Attach() {
	// Updating uniforms is expensive, avoid if possible
	if !uni.updated {
		return
	}

	uni.updated = false

	switch uni.value.(type) {
	case *float32:
		value := *(uni.value).(*float32)
		gl.Uniform1f(int32(uni.id), value)
	case *mgl32.Vec2:
		value := *(uni.value).(*mgl32.Vec2)
		gl.Uniform2f(int32(uni.id), value.X(), value.Y())
	case *mgl32.Vec3:
		value := *(uni.value).(*mgl32.Vec3)
		gl.Uniform3f(int32(uni.id), value.X(), value.Y(), value.Z())
	case *mgl32.Vec4:
		value := *(uni.value).(*mgl32.Vec4)
		gl.Uniform4f(int32(uni.id), value.X(), value.Y(), value.Z(), value.W())
	default:
		fmt.Println(uni.Value())
		panic("Unsupported uniform type, these should be pointers")
	}
}

func (p *Program) AddUniform(name string, value interface{}) {
	uni := uniform{
		uint32(gl.GetUniformLocation(p.Id, gl.Str(name+"\x00"))),
		value,
		true,
	}
	uni.Attach()
	p.uniforms[name] = uni
}

/*
Set uniform value, this should be set to pointers as this is an expensive way to update,
if set to a pointer then call Update()
*/

func (p *Program) SetUniform(name string, value interface{}) {
	if _, exists := p.uniforms[name]; !exists {
		panic("Attempting to set non existent uniform")
	}

	uni := uniform{
		p.uniforms[name].id,
		value,
		true,
	}

	uni.Attach()
	p.uniforms[name] = uni
}

func (p *Program) UpdateUniform(name string) {
	u, exists := p.uniforms[name]
	if !exists {
		panic("Attempting to set non existent uniform")
	}

	u.Update()
}

func (p *Program) UpdateUniforms() {
	for _, uni := range p.uniforms {
		uni.Update()
		uni.Attach()
	}
}
