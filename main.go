package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"runtime"
)

// Image size
const imageWidth = 1920
const imageHeight = 1080

func main() {
	// Image
	img := image.NewGray(image.Rect(0, 0, imageWidth, imageHeight))

	// Shard image dimensions
	shardsInX, shardsInY := shardDimensions(runtime.NumCPU())

	// Show if we were to shard
	fmt.Println("If we were to shard")
	imgShardW := imageWidth / shardsInX
	imgShardH := imageHeight / shardsInY
	fmt.Printf("Image shard grid: %d by %d\n", shardsInX, shardsInY)
	fmt.Printf("Image shard size: %d by %d\n", imgShardW, imgShardH)

	// Cell calculation
	var iters uint8
	for i := 0; i < imageWidth; i++ {
		for j := 0; j < imageHeight; j++ {
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

func shardDimensions(n int) (int, int) {
	a := int(math.Sqrt(float64(n)))
	b := n / a
	for a*b != n {
		if a*b < n {
			a++
		} else if a*b > n {
			b--
		}
	}
	if a < b {
		return b, a
	} else {
		return a, b
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
