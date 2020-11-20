package graphics

import "Gopengl2/graphics/opengl"

type DefaultRenderObject struct {
	*BaseRenderObject
}

func CreateDefaultRenderObject(texture string, elements int) *DefaultRenderObject {
	vao := opengl.CreateDefaultVao(window, texture, elements)

	baseRo := &BaseRenderObject{
		vao,
		0,
		vao.GetBuffer("vert"),
		vao.GetBuffer("verttexcoord"),
		true,
		false,
		false,
		false,
	}

	ro := &DefaultRenderObject{baseRo}

	renderObjects = append(renderObjects, ro)

	return ro
}

func (ro *DefaultRenderObject) CreateRect(x, y, width, height, texX, texY, texWidth, texHeight int) int {
	index := ro.freeVert
	ro.freeVert += 6
	ro.ModifyRect(index, x, y, width, height, texX, texY, texWidth, texHeight)

	return index
}

func (ro *DefaultRenderObject) ModifyRect(index, x, y, width, height, texX, texY, texWidth, texHeight int) {
	ro.SetVertex(index, x, y, texX, texY+texHeight)
	ro.SetVertex(index+1, x+width, y, texX+texWidth, texY+texHeight)
	ro.SetVertex(index+2, x, y+height, texX, texY)
	ro.SetVertex(index+3, x, y+height, texX, texY)
	ro.SetVertex(index+4, x+width, y+height, texX+texWidth, texY)
	ro.SetVertex(index+5, x+width, y, texX+texWidth, texY+texHeight)
}

func (ro *DefaultRenderObject) RemoveSquare(index int) {
	ro.ModifyRect(index, 0, 0, 0, 0, 0, 0, 0, 0)
}

func (ro *DefaultRenderObject) CreateSquare(x, y, width, texX, texY, texWidth int) int {
	return ro.CreateRect(x, y, width, width, texX, texY, texWidth, texWidth)
}

func (ro *DefaultRenderObject) ModifySquare(index, x, y, width, texX, texY, texWidth int) {
	ro.ModifyRect(index, x, y, width, width, texX, texY, texWidth, texWidth)
}

func (ro *DefaultRenderObject) SetTranslation(x, y float32) {
	ro.vao.SetTranslation(x, y)
}
