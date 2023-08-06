package controller

import (
	"Momotok-Server/rpc"
	"Momotok-Server/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if utils.CheckToken(token) {
		Uid, err := utils.GetUID(token)
		if err != nil {
			fmt.Println("Get UID Error:", err)
		}
		var user User
		user, err = getUserFromdb(Uid)

		if err != nil {
			fmt.Println("Get User Struct Error:", err)
		}

		videoID := c.Query("video_id")

		if actionType == "1" {
			text := c.Query("comment_text")
			// 获取当前时间
			currentTime := time.Now()
			// 将时间格式化为 "MM-DD" 格式
			currentDate := currentTime.Format("01-02")
			db, err := sql.Open("mysql", DatabaseAddress) //连接数据库
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

			fmt.Println("last insert id :", commentID)
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         commentID,
					User:       user,
					Content:    text,
					CreateDate: currentDate,
				}})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoID := c.Query("video_id")
	Comments, err := makeCommentList(videoID)
	if err != nil {
		fmt.Println("Make CommentList Error:", err)
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: Comments,
	})
}

func makeCommentList(videoID string) ([]Comment, error) {
	db, err := sql.Open("mysql", DatabaseAddress) //连接数据库
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM comments WHERE video_id = ? ORDER BY create_date DESC", videoID) //写入sql指令，按倒序查找列                                                                           //执行上述指令
	if err != nil {
		fmt.Println("Failed to execute query:", err)
	}

	defer rows.Close()
	Comments := make([]Comment, 0) //创建视频列表
	for rows.Next() {
		//循环读取直到列结束
		var id int64
		var video_id int64
		var commenter_id int64
		var content string
		var create_date string
		err := rows.Scan(&id, &video_id, &commenter_id, &content, &create_date)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}

		t, err := time.Parse("2006-01-02 15:04:05", create_date)
		if err != nil {
			fmt.Println("Failed to parse timestamp:", err)
			continue
		}

		// 格式化时间为 MM-DD 的字符串
		formattedDate := t.Format("01-02")

		user, err := getUserFromdb(commenter_id)

		comment := Comment{ //载入评论结构
			Id:         id,
			User:       user,
			Content:    content,
			CreateDate: formattedDate,
		}
		Comments = append(Comments, comment)
	}
	return Comments, nil
}

func getUserFromdb(uid int64) (User, error) {
	var user User
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
	user.Signature = string(signature)
	user.Avatar = "https://acg.suyanw.cn/sjtx/random.php"
	user.BackgroundImage = "https://acg.suyanw.cn/api.php"

	db, err := sql.Open("mysql", DatabaseAddress) //连接数据库
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, username, work_count  FROM user WHERE id = ?", uid)
	err = row.Scan(&user.Id, &user.Name, &user.WorkCount)
	if err != nil {
		fmt.Println("Failed to scan row:", err)
	}
	return user, nil
}
