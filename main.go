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
	"sync"
)

// Image size
const imageWidth = 3072
const imageHeight = 1920

// Types
type iterations uint32

// Maximum iterations
const maxIterations = 2047

func main() {
	// Image
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Shard image dimensions
	shardsInX, shardsInY := shardDimensions(runtime.NumCPU())
	imgShardW := imageWidth / shardsInX
	imgShardH := imageHeight / shardsInY

	// Run shards
	fmt.Printf("Parallel grid %d x %d of %d x %d shards\n", shardsInX, shardsInY, imgShardW, imgShardH)
	fmt.Printf("Remainder: %d x %d\n", imageWidth%shardsInX, imageHeight%shardsInY)
	var i1, j1 int
	var wg sync.WaitGroup
	for i0 := 0; i0 < imageWidth; i0 += imgShardW {
		for j0 := 0; j0 < imageHeight; j0 += imgShardH {
			i1 = i0 + imgShardW
			j1 = j0 + imgShardH
			wg.Add(1)
			go shardWorker(i0, j0, i1, j1, img, &wg)
		}
	}
	wg.Wait()

	// Write image
	fmt.Printf("Saving image...")
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
	fmt.Printf("Job's done!\n")
}

func shardWorker(i0, j0, i1, j1 int, img *image.RGBA, wg *sync.WaitGroup) {
	var iters iterations
	var color color.RGBA
	for i := i0; i < i1; i++ {
		for j := j0; j < j1; j++ {
			iters = mandelbrot(i, j)
			color = iterationToColor(iters)
			img.SetRGBA(i, j, color)
		}
	}
	(*wg).Done()
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

func mandelbrot(i, j int) iterations {
	x := float64(i)
	y := float64(j)
	rescale := float64(math.Min(imageWidth, imageHeight))
	cX := 4.0*(x/rescale) - 2.0*imageWidth/rescale
	cY := 2.0 - 4.0*(y/rescale)
	c := complex(cX, cY)
	z := 0 + 0i
	var n iterations
	for n = 0; n < maxIterations && cmplx.Abs(z) < 2; n++ {
		z = z*z + c
	}
	return n
}

func iterationToColor(iter iterations) color.RGBA {
	rescale := uint8(iter * 255 / maxIterations)
	return color.RGBA{R: rescale, G: rescale, B: rescale, A: 255}
}
