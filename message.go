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

type ChatResponse struct {
	model.Response
	MessageList []model.Message `json:"message_list"`
}

// MessageAction handles sending message
func MessageAction(c *gin.Context) {
	tokenString := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	actionType := c.Query("action_type")
	createAt := c.Query("create_at")
	if senderUID, err := utils.GetUID(tokenString); err == nil && actionType == "1" {
		db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
		if err != nil {
			fmt.Println("Database connected failed: ", err)
			return
		}
		db.Exec("INSERT INTO messages(sender_id, retriever_id, message, create_at) value(?, ?, ?, ?)", senderUID, toUserId, content, createAt)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
		return
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "Invalid request"})
		return
	}
}

// MessageChat provides user with message list
func MessageChat(c *gin.Context) {
	tokenString := c.Query("token")
	toUserId := c.Query("to_user_id")
	preMsgTime := c.Query("pre_msg_time")
	var msgList = make([]model.Message, 0)
	if utils.CheckToken(tokenString) {
		db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress) //连接数据库
		rows, err := db.Query("SELECT * FROM messages WHERE retriever_id = ? AND create_at > ?", toUserId, preMsgTime)
		if err != nil {
			fmt.Println("Failed to connect to database:", err)
		}
		for rows.Next() {
			var id int64
			var toUserID int64
			var fromUserID int64
			var content string
			var createdTime string
			err := rows.Scan(&id, &fromUserID, &toUserID, &createdTime, &content)
			if err != nil {
				fmt.Println(err)
				return
			}
			t, err := time.Parse("2006-01-02 15:04:05", createdTime)
			msgStruct := model.Message{
				Content:    content,
				CreateTime: t.Unix(),
				FromUserID: fromUserID,
				ID:         id,
				ToUserID:   toUserID,
			}
			msgList = append(msgList, msgStruct)
		}
		if err != nil {
			fmt.Println("Failed to mark messages as read:", err)
		}
		c.JSON(http.StatusOK, ChatResponse{
			Response:    model.Response{StatusCode: 0},
			MessageList: msgList,
		})
	} else {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
	}
}
