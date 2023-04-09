package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

// Image size
const imageWidth = 600
const imageHeight = 400

func main() {
	// Image
	img := image.NewGray(image.Rect(0, 0, imageWidth, imageHeight))

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

func mandelbrot(i, j int) uint8 {
	x := float64(i)
	y := float64(j)
	rescale := float64(math.Min(imageWidth, imageHeight))
	cX := 4.0*(x/rescale) - 2.0*imageWidth/rescale
	cY := 2.0 - 4.0*(y/rescale)
	c := complex(cX, cY)
	z := 0 + 0i
	var n uint8
	for n = 0; n < 255 && cmplx.Abs(z) < 2; n++ {
		z = z*z + c
	}
	return n
}
