package video

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/teohen/ascii-render/commom"
)

const VIDEO_SCALE_X = 2
const VIDEO_SCALE_Y = 4

type FFmpegLoader struct {
	Cmd    *exec.Cmd
	Output io.ReadCloser
	Height int
	Width  int
	FPS    int
}

func LoadVideo(inputFile string) (FFmpegLoader, error) {
	var ffmpg FFmpegLoader
	var err error

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height,r_frame_rate", "-of", "default=noprint_wrappers=1:nokey=1", inputFile)
	output, err := cmd.Output()

	if err != nil {
		return ffmpg, err
	}

	res := strings.Split(string(output), "\n")

	ffmpg.Width, err = strconv.Atoi(res[0])
	ffmpg.Height, err = strconv.Atoi(res[1])
	operands := strings.Split(res[2], "/")

	firstOp, error := strconv.Atoi(operands[0])
	secOp, error := strconv.Atoi(operands[1])

	if error != nil {
		log.Fatal("Fail to calculate fps")
	}

	ffmpg.FPS = firstOp / secOp

	if err != nil {
		return ffmpg, err
	}

	ffmpg.Cmd = exec.Command("ffmpeg", "-i", inputFile, "-f", "image2pipe", "-pix_fmt", "rgb24", "-vcodec", "rawvideo", "-")
	ffmpg.Output, err = ffmpg.Cmd.StdoutPipe()

	if err != nil {
		fmt.Println("Error getting stdout pipe:", err)
		return FFmpegLoader{}, errors.New("stdout pipe not found")
	}

	if err := ffmpg.Cmd.Start(); err != nil {
		fmt.Println("Error starting FFmpeg:", err)
		return FFmpegLoader{}, errors.New("start pipe not")
	}

	return ffmpg, nil
}

func CreateFrames(video FFmpegLoader, frameCh chan image.Image, wg *sync.WaitGroup) {
	width, height := video.Width, video.Height
	frameSize := width * height * 3

	buffer := make([]byte, frameSize)
	for frameIndex := 0; ; frameIndex++ {
		_, err := io.ReadFull(video.Output, buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading frame:", err)
			return
		}
		img := commom.FrameToImage(buffer, width, height, frameIndex)
		frameCh <- img
	}
	wg.Done()
}
func RenderFrames(frameCh chan image.Image, wg *sync.WaitGroup, fps, scaleX, scaleY int) {
	finalTicker := 1000000 / fps

	ticker := time.NewTicker(time.Duration(finalTicker) * time.Microsecond)
	defer ticker.Stop()

	for range ticker.C {
		frame := <-frameCh
		commom.Render(frame, scaleX, scaleY)
	}
	wg.Done()
}
func RenderVideoFile(filePath string, resizeWindow bool) {
	videoCmd, err := LoadVideo(filePath)

	if err != nil {
		log.Fatal("not able to load video")
	}

	if resizeWindow {
		err := commom.ResizeTerminal(videoCmd.Width/VIDEO_SCALE_X+1, videoCmd.Height/VIDEO_SCALE_Y+1)
		if err != nil {
			log.Fatal("Error resizing window: ", err.Error())
		}

	}

	var wg sync.WaitGroup
	frameCh := make(chan image.Image, 2000)

	go CreateFrames(videoCmd, frameCh, &wg)
	go RenderFrames(frameCh, &wg, videoCmd.FPS, VIDEO_SCALE_X, VIDEO_SCALE_Y)
	wg.Add(2)
	err = videoCmd.Cmd.Wait()
	wg.Wait()
	close(frameCh)

	if err != nil {
		fmt.Println("Error waiting for FFmpeg to finish:", err)
		return
	}
}
