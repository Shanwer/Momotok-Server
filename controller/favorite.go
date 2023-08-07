package controller

import (
	"Momotok-Server/rpc"
	"Momotok-Server/system"
	"Momotok-Server/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

// FavoriteAction handles user's like action
func FavoriteAction(c *gin.Context) {
	tokenString := c.Query("token")
	videoID := c.Query("video_id")
	actionType := c.Query("action_type")
	if uid, err := utils.GetUID(tokenString); err == nil {
		db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
		if err != nil {
			fmt.Println("Database connected failed: ", err)
			return
		}
		likeID := 0
		_ = db.QueryRow("SELECT id FROM likes WHERE video_id = ? AND user_id = ?", videoID, uid).Scan(&likeID)
		authorID := 0
		err = db.QueryRow("SELECT author_id FROM video WHERE id = ?", videoID).Scan(&authorID)
		if likeID == 0 && actionType == "1" {
			//likeID never exists when user likes a video
			_, err = db.Exec("INSERT INTO likes (video_id, user_id) VALUES (?, ?)", videoID, uid)
			_, err = db.Exec("UPDATE user SET total_likes = total_likes + 1 where id = ?", uid)
			_, err = db.Exec("UPDATE user SET total_received_likes = total_received_likes + 1 where id = ?", authorID)
			_, err = db.Exec("UPDATE video SET favourite_count = video.favourite_count + 1 where id = ?", videoID)
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 0},
			})
			return
		} else if likeID != 0 && actionType == "2" {
			//likeID must exist to unlike a video
			_, err = db.Exec("DELETE FROM likes where id = ?", likeID)
			_, err = db.Exec("UPDATE user SET total_likes = total_likes - 1 where id = ?", uid)
			_, err = db.Exec("UPDATE user SET total_received_likes = total_received_likes - 1 where id = ?", authorID)
			_, err = db.Exec("UPDATE video SET favourite_count = video.favourite_count - 1 where id = ?", videoID)
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 0},
			})
			return
		}
	} else {
		//check token or getUID returns err
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Invalid request"},
		})
		return

	}
}

// FavoriteList shows user's liked videos
func FavoriteList(c *gin.Context) {
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Unauthorized request"})
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
	videoList := make([]Video, 0)
	signature, _ := io.ReadAll(resp.Body)
	var user = User{
		Id:              parseIntID,
		Signature:       string(signature),
		Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
		BackgroundImage: "https://acg.suyanw.cn/api.php",
	}
	rows, err := db.Query("select video_id, play_url, cover_url, favourite_count, comment_count, title, user.username, user.total_received_likes, user.work_count, user.total_likes FROM video JOIN likes on video.id = likes.video_id JOIN user ON video.author_id = user.id where user_id = ?", userId)
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
		} //find liked videos
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
