package img

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/teohen/ascii-render/commom"
)

const IMG_SCALE_X = 1
const IMG_SCALE_Y = 2

func LoadImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	image, format, err := image.Decode(f)
	if err != nil {
		fmt.Println("error loading image", format, err)
	}
	return image, err
}
func RenderImgFile(filePath string, resizeWindow bool) {
	img, err := LoadImage(filePath)

	imageMax := img.Bounds().Max
	if resizeWindow {
		err := commom.ResizeTerminal(imageMax.X+1, imageMax.Y/2+1)
		if err != nil {
			log.Fatal("Error resizing window: ", err.Error())
		}
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	commom.Render(img, IMG_SCALE_X, IMG_SCALE_Y)
}
