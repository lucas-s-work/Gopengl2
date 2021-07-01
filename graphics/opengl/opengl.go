package opengl

import "github.com/go-gl/gl/v4.1-core/gl"

var (
	freeVaos []uint32
	vaoFree  []bool
)

const (
	MaxVAO = 128
)

func GlInit() {
	err := gl.Init()

	if err != nil {
		panic(err)
	}

	// Workaround for non-uniqueness on MacOS, halves GPU usage.
	freeVaos = make([]uint32, MaxVAO)
	vaoFree = make([]bool, MaxVAO)
	gl.GenVertexArrays(MaxVAO, &freeVaos[0])
	for i := range vaoFree {
		vaoFree[i] = true
	}
}

func GetVAOId() uint32 {
	for i, free := range vaoFree {
		if free {
			vaoFree[i] = false
			return freeVaos[i]
		}
	}

	panic("No free VAO id's remain")
}
