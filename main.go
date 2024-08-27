package main

import (
	"flag"
	_ "image/png"
	"log"

	"github.com/teohen/ascii-render/img"
	"github.com/teohen/ascii-render/video"
)

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
