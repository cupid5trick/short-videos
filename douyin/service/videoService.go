package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"douyin/douyin/dao"
	"douyin/douyin/model"

	"gorm.io/gorm"
)

const videoNum = 2 // feed每次返回的视频数量

// FeedGet 获得视频列表
func FeedGet(lastTime int64) ([]model.Video, error) {
	// t := time.Now()
	// fmt.Println(t)
	if lastTime == 0 { // 没有传入参数或者视屏已经刷完
		lastTime = time.Now().Unix()
	}
	strTime := fmt.Sprint(time.Unix(lastTime, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("查询的时间", strTime)

	VideoList := make([]model.Video, 0)
	err := dao.Db.Table("videos").Where("created_at < ?", strTime).Order("created_at desc").Limit(videoNum).Find(&VideoList).Error
	return VideoList, err
}

// AddCommentCount add comment_count
func AddCommentCount(videoId uint) error {

	if err := dao.Db.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

// ReduceCommentCount reduce comment_count
func ReduceCommentCount(videoId uint) error {

	if err := dao.Db.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}

// GetVideoAuthor get video author
func GetVideoAuthor(videoId uint) (uint, error) {
	var video model.Video
	if err := dao.Db.Table("videos").Where("id = ?", videoId).Find(&video).Error; err != nil {
		return video.ID, err
	}
	return video.AuthorId, nil
}

// CreateVideo 添加一条视频信息
func CreateVideo(video *model.Video) {
	dao.Db.Table("videos").Create(&video)
}

// GetVideoList 根据用户id查找 所有与该用户相关视频信息
func GetVideoList(userId uint) []model.Video {
	var videoList []model.Video
	dao.Db.Table("videos").Where("author_id=?", userId).Find(&videoList)
	return videoList
}

// ExampleReadFrameAsJpeg 获取封面
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}
