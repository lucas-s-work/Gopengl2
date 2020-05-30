package opengl

import (
	"Gopengl2/util"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	GlWindow      *glfw.Window
	Width, Height float64
	Name          string
	// Mouse and keyboard
	KeyMap                 map[string]bool
	MouseX, MouseY         int
	Mouse1, Mouse2, Mouse3 bool
}

// Window Creation and destruction

func CreateWindow(width, height int, name string) *Window {
	err := glfw.Init()

	util.CheckErr(err)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, name, nil, nil)

	util.CheckErr(err)

	window.MakeContextCurrent()

	w := Window{
		GlWindow: window,
		Width:    float64(width),
		Height:   float64(height),
		Name:     name,
	}

	return &w
}

func DestroyWindow(window *glfw.Window) {
	window.Destroy()
}

func (w *Window) SwapBuffers() {
	w.GlWindow.SwapBuffers()
}

// Input handling
func (w *Window) PollInput() {
	glfw.PollEvents()
	window := w.GlWindow

	//Get Keyboard input
	w.KeyMap["w"] = window.GetKey(glfw.KeyW) == glfw.Press
	w.KeyMap["a"] = window.GetKey(glfw.KeyA) == glfw.Press
	w.KeyMap["s"] = window.GetKey(glfw.KeyS) == glfw.Press
	w.KeyMap["d"] = window.GetKey(glfw.KeyD) == glfw.Press

	//Get Mouse input
	w.Mouse1 = window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press
	w.Mouse2 = window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
	w.Mouse3 = window.GetMouseButton(glfw.MouseButtonMiddle) == glfw.Press

	mX, mY := window.GetCursorPos()

	w.MouseX, w.MouseY = w.ScreenToPix(mX, mY)
}

func (w *Window) Key(key string) bool {
	return w.KeyMap[key]
}

func (w *Window) KeyCombo(keys ...string) bool {
	for _, key := range keys {
		if !w.Key(key) {
			return false
		}
	}

	return true
}

// Pixel to Screen coordinate conversion
// We deal with 0,0 being the bottom left coordinates

func (w *Window) ScreenToPix(x, y float64) (int, int) {
	// Adjust to bottom left to be 0,0
	x += 1
	y += 1

	x *= (w.Width / 2)
	y *= (w.Height / 2)

	return int(x), int(y)
}
