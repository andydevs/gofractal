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
const imageWidth = 1920
const imageHeight = 1080

func main() {
	// Image
	img := image.NewGray(image.Rect(0, 0, imageWidth, imageHeight))

	// Shard image dimensions
	shardsInX, shardsInY := shardDimensions(runtime.NumCPU())
	imgShardW := imageWidth / shardsInX
	imgShardH := imageHeight / shardsInY

	// Run shards
	fmt.Printf("Parallel grid %d x %d of %d x %d shards\n", shardsInX, shardsInY, imgShardW, imgShardH)
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

func shardWorker(i0, j0, i1, j1 int, img *image.Gray, wg *sync.WaitGroup) {
	var iters uint8
	for i := i0; i < i1; i++ {
		for j := j0; j < j1; j++ {
			iters = mandelbrot(i, j)
			img.SetGray(i, j, color.Gray{Y: iters})
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
