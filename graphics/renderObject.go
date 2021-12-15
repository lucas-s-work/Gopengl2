package graphics

import (
	"fmt"

	"github.com/lucas-s-work/gopengl2/graphics/opengl"
)

type RenderObject interface {
	SetVertex(int, int, int, int, int)
	SetAutoUpdate(bool)
	UpdateBuffers()
	PrepRender()
	SetWait(bool) chan WaitSignal
	Render()
	Delete()
	Created() bool
	GetVAO() opengl.VAO
}

/*
Currently uses DefaultVAO but this can be changed in the parent
object, this does not implement any of the attribute/uniform setting
*/

type WaitSignal struct{}

type BaseRenderObject struct {
	vao                        *opengl.DefaultVAO
	freeVert                   int
	vBuff                      *opengl.Buffer
	tBuff                      *opengl.Buffer
	ShouldRender               bool
	shouldWait                 bool
	waitChan                   chan WaitSignal
	updated, autoUpdate, async bool
}

func CreateBaseRenderObject(texture string, elements int) *BaseRenderObject {
	vao := opengl.CreateDefaultVao(window, texture, elements)

	ro := &BaseRenderObject{
		vao,
		0,
		vao.GetBuffer("vert"),
		vao.GetBuffer("verttexcoord"),
		true,
		false,
		make(chan WaitSignal),
		false,
		false,
		false,
	}

	renderObjects = append(renderObjects, ro)

	return ro
}

func (ro *BaseRenderObject) SetVertex(index, x, y, texX, texY int) {
	i := index * 2

	tX, tY := ro.vao.PixToTex(texX, texY)

	ro.vBuff.Elements[i] = float32(x)
	ro.vBuff.Elements[i+1] = float32(y)
	ro.tBuff.Elements[i] = tX
	ro.tBuff.Elements[i+1] = tY

	ro.updated = true
}

func (ro *BaseRenderObject) SetAutoUpdate(update bool) {
	if ro.async {
		panic("cannot auto-update async render object")
	}
	ro.autoUpdate = update
}

func (ro *BaseRenderObject) UpdateBuffers() {
	ro.vao.BindVao()
	ro.vao.UpdateBuffers()
	ro.updated = false
}

func (ro *BaseRenderObject) PrepRender() {
	if ro.updated && ro.autoUpdate {
		ro.vao.BindVao()
		ro.vao.UpdateBuffers()
		ro.updated = true
	}
	ro.vao.PrepRender()
}

func (ro *BaseRenderObject) SetWait(shouldWait bool) chan WaitSignal {
	ro.shouldWait = shouldWait

	if shouldWait {
		return ro.waitChan
	} else {
		return nil
	}
}

func (ro *BaseRenderObject) Render() {
	if ro.shouldWait {
		fmt.Println("waiting")
		<-ro.waitChan
	}

	if ro.ShouldRender {
		ro.PrepRender()
		ro.vao.Render()
	}
}

func (ro *BaseRenderObject) Delete() {
	ro.vao.Delete()
}

func (ro *BaseRenderObject) GetVAO() opengl.VAO {
	return ro.vao
}

func (ro *BaseRenderObject) Created() bool {
	if ro == nil {
		return false
	}

	return ro.vao != nil
}
