package main

import (
	"flag"
	"github.com/teohen/ascii-render/img"
	"github.com/teohen/ascii-render/video"
	_ "image/png"
	"log"
)

const IMG_SCALE_X = 1
const IMG_SCALE_Y = 2

func main() {
	var (
		filePath = flag.String("file", "", "path of the file to be rendered")
		fileType = flag.String("type", "img", "type of the file to be rendered (img, vid)")
	)

	flag.Parse()

	path := *filePath
	ftype := *fileType

	if path == "" {
		log.Fatal("a file path is required as argument '--file'")
	}

	if ftype == "vid" {
		video.RenderVideoFile(path, true)
	} else {
		img.RenderImgFile(path, true)
	}
}
