package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	// Image
	image := [60][40]uint8{}

	// Cell calculation
	for i, row := range image {
		for j := range row {
			image[i][j] = uint8((j * i) % 256)
		}
	}

	// Write image
	f, err := os.Create("notimage.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, row := range image {
		for _, cell := range row {
			toPrint := uint8(math.Pow10(int(cell) / 100))
			fmt.Fprintf(f, "%3d ", toPrint)
		}
		fmt.Fprintln(f)
	}
}
