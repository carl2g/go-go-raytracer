package gogoraytracer

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
)

type PixelPosition struct {
	X, Y int32
}

type Color struct {
	R, G, B uint8
}

var blackColor = &Color{
	R: 0,
	G: 0,
	B: 0,
}

var whiteColor = &Color{
	R: 255,
	G: 255,
	B: 255,
}

func (c *Color) Copy() *Color {
	return &Color{
		R: c.R,
		G: c.G,
		B: c.B,
	}
}

func (c *Color) Scale(scale float64) *Color {
	return &Color{
		R: uint8(math.Trunc(float64(c.R) * scale)),
		G: uint8(math.Trunc(float64(c.G) * scale)),
		B: uint8(math.Trunc(float64(c.B) * scale)),
	}
}

func (c *Color) Add(otherColor *Color) *Color {
	return &Color{
		R: uint8(math.Min(math.Trunc(float64(c.R)+float64(otherColor.R)), 255)),
		G: uint8(math.Min(math.Trunc(float64(c.G)+float64(otherColor.G)), 255)),
		B: uint8(math.Min(math.Trunc(float64(c.B)+float64(otherColor.B)), 255)),
	}
}

type Pixel struct {
	PixelPosition *PixelPosition
	Color         *Color
}

type Frame struct {
	Width, Height int
	Img           *image.RGBA
}

func NewFrame(height, width int) (frame *Frame) {
	frame = &Frame{
		Width:  width,
		Height: height,
		Img: image.NewRGBA(
			image.Rect(0, 0, width, height),
		),
	}
	return frame
}

func (frame *Frame) getHeight() int {
	return frame.Height
}

func (frame *Frame) getWidth() int {
	return frame.Width
}

func (frame *Frame) Draw() {
	outFile, err := os.Create("output.png")
	if err != nil {
		log.Printf("failed to create output file: %v", err)
		return
	}
	defer outFile.Close()

	if err := png.Encode(outFile, frame.Img); err != nil {
		log.Printf("failed to encode png: %v", err)
		return
	}

	if err := exec.Command("xdg-open", "output.png").Start(); err != nil {
		log.Printf("failed to open image: %v", err)
	}
}

func (frame *Frame) SetPixel(pixel *Pixel) {
	x := int(pixel.PixelPosition.X)
	y := int(pixel.PixelPosition.Y)
	if x < 0 || y < 0 || x >= frame.Width || y >= frame.Height {
		return
	}
	frame.Img.Set(x, y, color.RGBA{pixel.Color.R, pixel.Color.G, pixel.Color.B, 255})
}
