package controller

import (
	"Momotok-Server/model"
	"Momotok-Server/system"
	"Momotok-Server/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	model.Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentAction handles comment or delete comments action
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	Uid, err := utils.GetUID(token)
	if err != nil {
		fmt.Println("User doesn't exist", err)
	}
	var user model.User
	user, err = utils.GetUserStruct(Uid)

	if err != nil {
		fmt.Println("Get User Struct Error:", err)
		return
	}

	videoID := c.Query("video_id")

	if actionType == "1" {
		text := c.Query("comment_text")
		// 获取当前时间
		currentTime := time.Now()
		// 将时间格式化为 "MM-DD" 格式
		currentDate := currentTime.Format("01-02")
		db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress) //连接数据库
		if err != nil {
			fmt.Println("Failed to connect to database:", err)
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO comments (video_id,commenter_id,content) VALUES (?,?,?)", videoID, Uid, text)
		if err != nil {
			fmt.Println("Comment Upload failed：", err)
		}

		_, err = db.Exec("UPDATE video SET comment_count = comment_count + 1 WHERE id = ?", videoID)
		if err != nil {
			fmt.Println("Failed to update video comment count:", err)
		}

		var commentID int64
		row := db.QueryRow("SELECT LAST_INSERT_ID()")
		if err := row.Scan(&commentID); err != nil {
			fmt.Println("Failed to get comment ID:", err)
		}

		c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
			Comment: model.Comment{
				Id:         commentID,
				User:       user,
				Content:    text,
				CreateDate: currentDate,
			}})
		return
	} else if actionType == "2" {
		db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress) //连接数据库
		if err != nil {
			fmt.Println("Failed to connect to database:", err)
		}
		defer db.Close()
		commentID := c.Query("comment_id")
		_, err = db.Exec("DELETE FROM comments WHERE id = ?", commentID)
		if err != nil {
			fmt.Println("Comment Upload failed：", err)
		}
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
		return
	}
}

// CommentList provides comment list
func CommentList(c *gin.Context) {
	videoID := c.Query("video_id")
	Comments, err := makeCommentList(videoID)
	if err != nil {
		fmt.Println("Make CommentList Error:", err)
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: Comments,
	})
}

func makeCommentList(videoID string) ([]model.Comment, error) {
	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress) //连接数据库
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM comments WHERE video_id = ? ORDER BY create_date DESC", videoID) //写入sql指令，按倒序查找列
	if err != nil {
		fmt.Println("Failed to execute query:", err)
	}

	defer rows.Close()
	Comments := make([]model.Comment, 0) //创建视频列表
	for rows.Next() {
		//循环读取直到列结束
		var id int64
		var videoID int64
		var commenterID int64
		var content string
		var createDate string
		err := rows.Scan(&id, &videoID, &commenterID, &content, &createDate)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}

		t, err := time.Parse("2006-01-02 15:04:05", createDate)
		if err != nil {
			fmt.Println("Failed to parse timestamp:", err)
			continue
		}

		// 格式化时间为 MM-DD 的字符串
		formattedDate := t.Format("01-02")

		user, err := utils.GetUserStruct(commenterID)
		if err != nil {
			return nil, err
		}
		comment := model.Comment{ //载入评论结构
			Id:         id,
			User:       user,
			Content:    content,
			CreateDate: formattedDate,
		}
		Comments = append(Comments, comment)
	}
	return Comments, nil
}
