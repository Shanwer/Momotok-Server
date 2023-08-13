package controller

import (
	"Momotok-Server/model"
	"Momotok-Server/rpc"
	"Momotok-Server/system"
	"Momotok-Server/utils"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish function that check token then save upload file to public directory
func Publish(c *gin.Context) {
	tokenString := c.PostForm("token")
	uid, err := utils.GetUID(tokenString)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "invalid token",
		})
		return
	}

	title := c.PostForm("title")
	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
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
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	_, err = utils.GetSnapshot("./public/"+file.Filename, "./public/snapshot/"+file.Filename, 1) //get the first frame of the video
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	cover_url := system.ServerInfo.Server.StaticFileUrl + "snapshot/" + file.Filename + ".jpg"
	play_url := system.ServerInfo.Server.StaticFileUrl + file.Filename

	_, err = db.Exec("INSERT INTO video (author_id,title,favourite_count,comment_count,play_url,cover_url) VALUES (?,?,?,?,?,?)", uid, title, 0, 0, play_url, cover_url)
	if err != nil { //duplicate username check
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			fmt.Println("Upload failed：", err)
		}
		if mysqlErr.Number == 1062 {
			c.JSON(http.StatusOK, UserResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: "Video already exists!"},
			})
		} else {
			fmt.Println("Upload failed：", err)
		}
		return
	}

	_, err = db.Exec("UPDATE user SET work_count = work_count + 1 where id = ?", uid)

	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

// PublishList shows user's published videos
func PublishList(c *gin.Context) {
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "Unauthorized request"})
		return
	}
	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
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
	var user = model.User{
		Id:              parseIntID,
		Signature:       string(signature),
		Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
		BackgroundImage: "https://acg.suyanw.cn/api.php",
	}
	rows, err := db.Query("select video.id, play_url, cover_url, favourite_count, comment_count, title, user.username, user.total_received_likes, user.work_count, user.total_likes, user.follow_count, user.follower_count FROM video JOIN user ON video.author_id = user.id where author_id = ? ", parseIntID)
	videoList := make([]model.Video, 0)
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
		err := rows.Scan(&videoId, &playUrl, &coverUrl, &favoriteCount, &commentCount, &title, &user.Name, &user.TotalFavorited, &user.WorkCount, &user.FavoriteCount, &user.FollowCount, &user.FollowerCount)
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
		video := model.Video{ //载入视频结构
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
		Response: model.Response{
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
