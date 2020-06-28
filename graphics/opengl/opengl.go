package opengl

import "github.com/go-gl/gl/v4.1-core/gl"

func GlInit() {
	err := gl.Init()

	if err != nil {
		panic(err)
	}
}
