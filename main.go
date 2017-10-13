package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d"
)

var LIGHT_COLORS = false
var CYCLE_COLORS = true
var PRESET_FACTORS = true
var START_FACTOR = 2
var END_FACTOR = 2

var lightColors []color.Color = []color.Color{
	color.RGBA{0xF4, 0x9A, 0xC2, 0xFF},
	color.RGBA{0xB1, 0x9C, 0xD9, 0xFF},
	color.RGBA{0x77, 0x9E, 0xCB, 0xFF},
	color.RGBA{0xAE, 0xC6, 0xCF, 0xFF},
	color.RGBA{0x77, 0xDD, 0x77, 0xFF},
	color.RGBA{0xDD, 0xDD, 0x76, 0xFF},
	color.RGBA{0xFF, 0xB3, 0x47, 0xFF},
	color.RGBA{0xFF, 0x69, 0x61, 0xFF},
	color.RGBA{0x00, 0x00, 0x00, 0xFF},
}
var darkColors []color.Color = []color.Color{
	color.RGBA{0xF4, 0x9A, 0xC2, 0xFF},
	color.RGBA{0xB1, 0x9C, 0xD9, 0xFF},
	color.RGBA{0x77, 0x9E, 0xCB, 0xFF},
	color.RGBA{0xAE, 0xC6, 0xCF, 0xFF},
	color.RGBA{0x77, 0xDD, 0x77, 0xFF},
	color.RGBA{0xDD, 0xDD, 0x76, 0xFF},
	color.RGBA{0xFF, 0xB3, 0x47, 0xFF},
	color.RGBA{0xFF, 0x69, 0x61, 0xFF},
	color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
}

var presetFactors []int = []int{
	2, 3, 4, 5, 6, 7, 8, 9, 11, 21,
	22, 26, 28, 31, 34, 36, 41, 42,
	45, 46, 48, 51, 52, 53, 56, 57,
	58, 59, 63, 65, 67, 71, 72, 73,
	74, 75, 76, 77, 78, 79, 81, 83,
	84, 85, 86, 89, 91, 92, 94, 99,
}

func main() {
	// Set up colors and factors
	var colors []color.Color
	var background color.RGBA
	var factors []int

	if LIGHT_COLORS {
		colors = lightColors
		background = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	} else {
		colors = darkColors
		background = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}

	if PRESET_FACTORS {
		factors = presetFactors
	} else {
	    factors = make([]int, END_FACTOR - START_FACTOR + 1)
	    for i := range factors {
	        factors[i] = START_FACTOR + i
	    }
	}

	// Go through all the factors
	var wg sync.WaitGroup
	counter := 0
	for _, factor := range factors {
		if CYCLE_COLORS {
			// Use a different color for each image
			counter %= len(colors)
			colour := colors[counter]
			name := fmt.Sprintf("%v.png", factor)
			wg.Add(1)
			go makeImage(factor, 100, 500, name, colour, background, &wg)
			counter++
		} else {
			// Generate every color of an image
			for _, colour := range colors {
				r, g, b, _ := colour.RGBA()
				name := fmt.Sprintf("%v-%02X%02X%02X.png", factor, byte(r), byte(g), byte(b))
				wg.Add(1)
				go makeImage(factor, 100, 500, name, colour, background, &wg)
			}
		}
	}
	wg.Wait()
}

func makeImage(factor, modulo int, radius float64, name string, fg color.Color, bg color.Color, wg *sync.WaitGroup) {
	m := radius / 5
	// Initialize the graphic context on an RGBA image
	dest := image.NewRGBA(image.Rect(0, 0, int(((radius + m) * 2)), int((radius+m)*2)))
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetLineWidth(5)

	// Set Background color
	gc.SetStrokeColor(bg)
	gc.SetFillColor(bg)

	// Background
	gc.MoveTo(0, 0)
	gc.LineTo(((radius+m)*2)*1.75, 0)
	gc.LineTo(((radius+m)*2)*1.75, (radius+m)*2)
	gc.LineTo(0, (radius+m)*2)
	gc.LineTo(0, 0)
	gc.FillStroke()

	// label
	gc.SetFillColor(fg)
	gc.SetFontData(draw2d.FontData{Name: "arial", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	gc.SetFontSize(36)
	gc.FillStringAt(fmt.Sprintf("%v", factor), radius * 2, radius * 2 + m)
	gc.FillStroke()

	// Set foreground Color
	gc.SetStrokeColor(fg)
	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0x00})

	// Outline
	gc.ArcTo(radius+m, radius+m, radius, radius, 0, 2.0*math.Pi)
	gc.FillStroke()

	// Lines
	for i := 0; i < modulo/2; i++ {
		val := (i * factor) % modulo
		x := math.Cos(2.0*math.Pi*(float64(i)/float64(modulo)))*radius + radius + m
		y := math.Sin(2.0*math.Pi*(float64(i)/float64(modulo)))*radius + radius + m
		gc.MoveTo(x, y)
		x = math.Cos(2.0*math.Pi*(float64(val)/float64(modulo)))*radius + radius + m
		y = math.Sin(2.0*math.Pi*(float64(val)/float64(modulo)))*radius + radius + m
		gc.LineTo(x, y)
	}
	gc.FillStroke()

	gc.Close()
	// Save to file
	draw2dimg.SaveToPngFile(name, dest)
	wg.Done()
}
