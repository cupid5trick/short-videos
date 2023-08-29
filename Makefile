APP=douyin
VERSION=0.0.1
FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
.PHONY: build clean
build:
	${FLAGS} go build -o build/douyin douyin/main.go
	${FLAGS} go build -o build/douyin-ffmpeg ffmpeg/main.go

clean:
	@find build/ -maxdepth 1 -type f -exec rm {} \;

