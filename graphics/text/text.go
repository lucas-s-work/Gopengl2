package text

import (
	"github.com/lucass-work/Gopengl2/graphics"
)

type Font struct {
	letterString string
	fontLocation string
	letterMap    map[rune]fontCoord
	letterWidth  int
	letterHeight int
}

type fontCoord struct {
	x, y int
}

type Text struct {
	font              Font
	R                 *graphics.DefaultRenderObject
	currentText       string
	currentTextIndexs []int
}

var (
	defaultFontLocation = "./resources/sprites/font.png"
	defaultLetterString = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"
	defaultFont         Font
)

func LoadFont(location, letters string) Font {
	lettersPerRow := 15
	letterMap := make(map[rune]fontCoord)

	i := 0
	j := 0
	width := 16
	height := 16
	for _, c := range letters {
		if i > lettersPerRow {
			i = 0
			j++
		}

		letterMap[c] = fontCoord{i * width, j * height}
	}

	return Font{location, letters, letterMap, 16, 16}
}

func LoadDefaultFont() {
	defaultFont = LoadFont("./resources/sprites/font.png", defaultLetterString)
}

// Create and initialize the text RO
func CreateText(text string, x, y int, font Font) Text {
	if &font == nil {
		font = defaultFont
	}

	if &font == nil {
		panic("font set to nil and default font not loaded.")
	}

	ro := graphics.CreateDefaultRenderObject(defaultFontLocation, 1000)
	indexs := make([]int, 1000)
	// Initialize the positions used for the render object
	for i := 0; i < 1000; i++ {
		indexs[i] = ro.CreateSquare(0, 0, 0, 0, 0, 0)
	}

	font.renderText(x, y, text, ro, indexs, 100)

	return Text{font, ro, text, indexs}
}

// Update the text RO
func (t Text) UpdateText(text string, x, y int) {
	t.font.renderText(x, y, text, t.R, t.currentTextIndexs, 100)
}

func (f Font) renderText(x, y int, text string, ro *graphics.DefaultRenderObject, indexs []int, wrap int) {
	// Perform this job asynchronously
	graphics.AddJobBlock(ro, func(r graphics.RenderObject) {
		ro := r.(*graphics.DefaultRenderObject)
		// remove all previous text
		for _, index := range indexs {
			ro.ModifyRect(index, 0, 0, 0, 0, 0, 0, 0, 0)
		}

		j := 0
		k := 0
		for i, c := range text {
			k++
			if k > wrap {
				k = 0
				j++
			}

			coord := f.letterMap[c]

			if &coord == nil {
				panic("Unable to find coord")
			}

			// Add current text with letter wrapping
			ro.ModifyRect(indexs[i], x+k*f.letterWidth, y+j*f.letterHeight, f.letterWidth, f.letterHeight, coord.x, coord.y, f.letterWidth, f.letterHeight)
		}
	})
}
