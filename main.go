package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math"
	"os"
	"strings"
)

func loadImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func main() {
	img, err := loadImage("teteo.png")
	if err != nil {
		panic(err)
	}
	mappedChars := " .:-=+*#%@"
	ramp := strings.Split(mappedChars, "")
	max := img.Bounds().Max
	scaleX, scaleY := 1, 2

	for i := 0; i < max.Y; i += scaleY {
		for j := 0; j < max.X; j += scaleX {
			pixelColor := img.At(j, i)
			r, g, b, _ := pixelColor.RGBA()
			avg := int((math.Floor(float64(r)/7281) + math.Floor(float64(g)/7281) + math.Floor(float64(b)/7281)) / 3)
			char := ramp[avg]
			fmt.Print(char)
		}
		fmt.Println()
	}

}
