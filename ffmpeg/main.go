package main

import (
	"fmt"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	stream := ffmpeg.Input("/home/cupid5trick/data/LearnProgramming/go/code/douyin-main/build/media/video/1_录屏取证操作示例.mp4")
	println(strings.Join(stream.Compile().Args, " "))
	cmd := stream.
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 3)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		Compile()
	// if cmd.Err != nil {
	// 	println("error: ", cmd.Err.Error())
	// }
	println(cmd.Args)
	script := strings.Join(cmd.Args, " ")
	println(script)
}
