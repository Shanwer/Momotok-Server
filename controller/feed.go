package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

var currentPage = 1 //全局变量记录当前page

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	perPage := 30                                 // 默认每页加载 30 个视频
	videos := makeVideoList(currentPage, perPage) //调用方法返回视频列表
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
	currentPage++
}

func makeVideoList(page, perPage int) []Video {
	db, err := sql.Open("mysql", DatabaseAddress) //连接数据库
	//数据库表名为video，字段为id, author_id, play_url, cover_url, favorite_count, comment_count, is_favorite, title，&publish_time具体类型见下述定义
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
	defer db.Close() //设置页数

	offSet := (page - 1) * perPage //offSet:视频开始位置

	rows, err := db.Query("SELECT * FROM video ORDER BY publish_time DESC LIMIT ? OFFSET ?", perPage, offSet) //写入sql指令，按倒序查找列                                                                           //执行上述指令
	if err != nil {
		fmt.Println("Failed to execute query:", err)
		return nil
	}
	defer rows.Close()
	videos := make([]Video, 0) //创建视频列表
	isLast := 0
	for rows.Next() {
		//循环读取直到列结束
		var id int64
		var author_id int64
		var play_url string
		var cover_url string
		var favorite_count int64
		var comment_count int64
		var title string
		var published_time []uint8 //TODO:未使用的变量
		err := rows.Scan(&id, &author_id, &play_url, &cover_url, &favorite_count, &comment_count, &title, &published_time)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		var likedID int
		isFavourite := false
		db.QueryRow("select id FROM likes where user_id = ? AND video_id = ?", author_id, id).Scan(&likedID)
		if likedID != 0 {
			isFavourite = true
		}

		user, err := getUserFromdb(author_id)

		video := Video{ //载入视频结构
			Id:            id,
			Author:        user,
			PlayUrl:       play_url,
			CoverUrl:      cover_url,
			FavoriteCount: favorite_count,
			CommentCount:  comment_count,
			IsFavorite:    isFavourite,
			Title:         title,
		}
		videos = append(videos, video) //视频切片加入视频列表
		isLast++
	}
	if isLast < perPage {
		currentPage = 0
	}
	return videos //返回视频列表
}
