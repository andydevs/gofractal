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
	"time"
)

// Image size
const imageWidth = 3072
const imageHeight = 1920

// Types
type iterations uint32

// Maximum iterations
const maxIterations = 2047

// Maximum threads
const maxThreads = 8

func main() {
	// Image
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	// Shard image dimensions
	threads := runtime.NumCPU()
	if threads > maxThreads {
		threads = maxThreads
	}
	fmt.Printf("Threads: %d\n", threads)
	shardsInX, shardsInY := shardDimensions(threads)
	imgShardW := imageWidth / shardsInX
	imgShardH := imageHeight / shardsInY
	fmt.Printf("Parallel grid %d x %d of %d x %d shards\n", shardsInX, shardsInY, imgShardW, imgShardH)

	// Run shards
	fmt.Printf("Running...")
	start := time.Now()
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
	elapsed := time.Since(start)
	fmt.Printf("Done in %s\n", elapsed.String())

	// Write image
	fmt.Printf("Saving image...\n")
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
