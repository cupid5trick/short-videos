package controller

import (
	"fmt"
	"image/jpeg"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"douyin/douyin/common"
	"douyin/douyin/model"
	"douyin/douyin/service"
)

type ReturnAuthor struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type ReturnMyself struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
}

type ReturnVideo struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnAuthor `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}
type ReturnVideo2 struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnMyself `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

type VideoListResponse struct {
	common.Response
	VideoList []ReturnVideo `json:"video_list"`
}
type VideoListResponse2 struct {
	common.Response
	VideoList []ReturnVideo2 `json:"video_list"`
}

func Publish(c *gin.Context) { // 上传视频方法
	var err error
	// 1.中间件验证token后，获取userId
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}

	// 2.接收请求参数信息
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		_ = c.Error(err)
	}
	log.Printf("Received media (%s:%s) from uid:%s", title, data.Filename, userId)
	// 4.保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	video := model.Video{
		Model:         gorm.Model{},
		AuthorId:      userId,
		PlayUrl:       "",
		CoverUrl:      "",
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	err = publishToLocal0(data, &video, c.SaveUploadedFile)
	if err != nil {
		_ = c.Error(err)
	}
	log.Printf("play_url: %s, cover_url: %s", video.PlayUrl, video.CoverUrl)
	service.CreateVideo(&video)
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  video.Title + "published successfully at" + video.PlayUrl,
	})
}

// 在本地调用 ffmpeg 获取第 3 帧作为封面
func publishToLocal0(data *multipart.FileHeader, video *model.Video, save func(*multipart.FileHeader, string) error) error {
	var err error
	// 先存储到本地文件夹, 文件名可能含有非法字符
	videoRoot := common.APPConfig.Fs.VideoRoot
	imageRoot := common.APPConfig.Fs.ImageRoot
	MediaURI := common.APPConfig.Fs.URI
	fileName := (filepath.Base(data.Filename))
	ext := filepath.Ext(data.Filename)
	userId := video.AuthorId
	finalName := fmt.Sprintf("%d_%s", userId, fileName)
	saveFile := filepath.Join(videoRoot, finalName)

	err = save(data, saveFile)
	if err != nil {
		return err
	}
	playUrl := filepath.Join(MediaURI+videoRoot, finalName)
	log.Printf("Saved video (%s:%s) to %s, accessiable at %s", video.Title, data.Filename, saveFile, playUrl)

	if err != nil {
		return err
	}

	// 存储封面图片到本地
	coverName := strings.Replace(finalName, ext, ".jpeg", 1)
	img := service.ExampleReadFrameAsJpeg(saveFile, 3) // 获取第3帧封面
	saveImage := filepath.Join(imageRoot, coverName)

	img2, _ := jpeg.Decode(img)     // 保存到本地时要用到
	imgw, _ := os.Create(saveImage) // 先创建，后写入
	err = jpeg.Encode(imgw, img2, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	coverUrl := filepath.Join(MediaURI+imageRoot, coverName)
	log.Printf("Extracted cover for video (%s:%s) to %s, accessiable at %s", video.Title, data.Filename, saveImage, coverUrl)

	// 保存完成后记录下播放地址和封面地址
	video.CoverUrl = coverUrl
	video.PlayUrl = playUrl

	return nil
}

// func publishToLocal1(data *multipart.FileHeader, video *model.Video, save func(*multipart.FileHeader, string) error) error {
//
// }

// PublishList 获取列表的方法
func PublishList(c *gin.Context) {
	// 1.中间件鉴权token
	getHostId, _ := c.Get("user_id")
	var HostId uint
	if v, ok := getHostId.(uint); ok {
		HostId = v
	}
	// 2.查询要查看用户的id的所有视频，返回页面
	getGuestId := c.Query("user_id")
	id, _ := strconv.Atoi(getGuestId)
	GuestId := uint(id)

	if GuestId == 0 || GuestId == HostId {
		// 根据token-id查找用户
		getUser, err := service.GetUser(HostId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "Not find this person.",
			})
			c.Abort()
			return
		}

		returnMyself := ReturnMyself{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
		}
		// 根据用户id查找 所有相关视频信息

		videoList := service.GetVideoList(HostId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "null",
				},
				VideoList: nil,
			})
		} else { // 需要展示的列表信息
			var returnVideoList2 []ReturnVideo2
			for i := 0; i < len(videoList); i++ {
				returnVideo2 := ReturnVideo2{
					VideoId:       videoList[i].ID,
					Author:        returnMyself,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList2 = append(returnVideoList2, returnVideo2)
			}
			c.JSON(http.StatusOK, VideoListResponse2{
				Response: common.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
				VideoList: returnVideoList2,
			})
		}
	} else {
		// 根据传入id查找用户
		getUser, err := service.GetUser(GuestId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "Not find this person.",
			})
			c.Abort()
			return
		}

		returnAuthor := ReturnAuthor{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
			IsFollow:      service.IsFollowing(HostId, GuestId),
		}

		// 根据用户id查找 所有相关视频信息
		videoList := service.GetVideoList(GuestId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 1,
					StatusMsg:  "null",
				},
				VideoList: nil,
			})
		} else { // 需要展示的列表信息
			var returnVideoList []ReturnVideo
			for i := 0; i < len(videoList); i++ {
				returnVideo := ReturnVideo{
					VideoId:       videoList[i].ID,
					Author:        returnAuthor,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList = append(returnVideoList, returnVideo)
			}
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					StatusCode: 0,
					StatusMsg:  "success",
				},
				VideoList: returnVideoList,
			})
		}
	}
}
