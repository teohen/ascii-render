package commom

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os/exec"
)

func Render(img image.Image, x, y int) {
	mappedChars := " .:-=+*#%@"
	ramp := strings.Split(mappedChars, "")
	max := img.Bounds().Max

	frame := ""
	for i := 0; i < max.Y; i += y {
		for j := 0; j < max.X; j += x {
			pixelColor := img.At(j, i)
			r, g, b, _ := pixelColor.RGBA()
			avg := int((math.Floor(float64(r)/7281) + math.Floor(float64(g)/7281) + math.Floor(float64(b)/7281)) / 3)
			char := ramp[avg]

			frame += char
		}
		frame += "\n"
	}
	fmt.Println(frame)
}

func FrameToImage(frame []byte, width, height, index int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := (y*width + x) * 3
			r := frame[index]
			g := frame[index+1]
			b := frame[index+2]
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

func ResizeTerminal(cols, rows int) error {
	cmd := exec.Command("resize", "-s ", fmt.Sprintf("%d", rows), fmt.Sprintf("%d", cols))
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("Error resizing window")
		return err
	}
	return nil
}
