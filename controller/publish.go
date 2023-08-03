package controller

import (
	"Momotok-Server/rpc"
	"Momotok-Server/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

var staticUrl = "http://192.168.31.224:8080/static/" //TODO:be stored in the server config, it should be accessed freely by the front end application

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish function that check token then save upload file to public directory
func Publish(c *gin.Context) {
	tokenString := c.PostForm("token")
	uid, err := utils.GetUID(tokenString)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "invalid token",
		})
		return
	}

	title := c.PostForm("title")
	db, err := sql.Open("mysql", DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	file, err := c.FormFile("data")
	if err != nil {
		return
	}
	file.Filename = hashFileName(file.Filename)
	err = c.SaveUploadedFile(file, "./public/"+file.Filename)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(), //possible issue: repeated filename
		})
		return
	}
	cover_url, err := utils.GetSnapshot("./public/"+file.Filename, "./public/snapshot/"+file.Filename, 1) //get the first frame of the video
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	cover_url = staticUrl + "snapshot/" + file.Filename + ".png"

	play_url := staticUrl + file.Filename
	_, err = db.Exec("INSERT INTO video (author_id,title,favourite_count,comment_count,play_url,cover_url) VALUES (?,?,?,?,?,?)", uid, title, 0, 0, play_url, cover_url)
	var workCount int
	err = db.QueryRow("select work_count from user").Scan(&workCount)
	if err != nil {
		return
	}
	workCount++
	_, err = db.Exec("UPDATE user SET work_count = ? where id = ?", workCount, uid) //TODO:寻找一种更好的办法

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

// PublishList shows user's published videos
func PublishList(c *gin.Context) {
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unauthorized request"})
		return
	}
	db, err := sql.Open("mysql", DatabaseAddress)
	userId := c.Query("user_id")
	parseIntID, _ := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	resp, _ := rpc.HttpRequest("GET", "https://v1.hitokoto.cn/?c=a&c=d&c=i&c=k&encode=text", nil)
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)
	}
	signature, _ := io.ReadAll(resp.Body)
	var user = User{
		Id:              parseIntID,
		Signature:       string(signature),
		Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
		BackgroundImage: "https://acg.suyanw.cn/api.php",
	}
	rows, err := db.Query("select video.id, play_url, cover_url, favourite_count, comment_count, title, user.username, user.total_received_likes, user.work_count, user.total_likes FROM video JOIN user ON video.author_id = user.id where author_id = ? ", parseIntID)
	videoList := make([]Video, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var videoId int64
		var playUrl string
		var coverUrl string
		var favoriteCount int64
		var commentCount int64
		var title string
		err := rows.Scan(&videoId, &playUrl, &coverUrl, &favoriteCount, &commentCount, &title, &user.Name, &user.TotalReceivedLikes, &user.WorkCount, &user.TotalLikes)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		} //find published videos
		var likedID int
		isFavourite := false
		db.QueryRow("select id FROM likes where user_id = ? AND video_id = ?", userId, videoId).Scan(&likedID)
		if likedID != 0 {
			isFavourite = true
		}
		video := Video{ //载入视频结构
			Id:            videoId,
			Author:        user,
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavourite,
		}
		videoList = append(videoList, video) //视频切片加入视频列表
	}
	defer rows.Close()

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

func hashFileName(fileName string) string {
	// 创建SHA256哈希对象
	hash := sha256.New()

	// 将文件名转换为字节数组并进行哈希计算
	hash.Write([]byte(fileName))

	// 获取哈希值并转换为十六进制字符串
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
