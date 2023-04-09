package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	// Image
	img := image.NewGray(image.Rect(0, 0, 600, 400))

	// Cell calculation
	var iters uint8
	bounds := img.Bounds().Max
	for i := 0; i < bounds.X; i++ {
		for j := 0; j < bounds.Y; j++ {
			iters = mandelbrot(i, j)
			img.SetGray(i, j, color.Gray{Y: iters})
		}
	}

	// Write image
	var err error
	f, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}

func mandelbrot(x, y int) uint8 {
	return uint8(x * y % 256)
}
