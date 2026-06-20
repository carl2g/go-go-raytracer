package main

import (
	gogoraytracer "go-go-raytracer/mod/internal"
)

func main() {
	pixelsChan := make(chan gogoraytracer.PixelPosition, 4096)
	frame := gogoraytracer.NewFrame(800, 1200)

	gogoraytracer.RunEngine(frame, pixelsChan)
}
