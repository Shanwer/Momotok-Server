package controller

import (
	"Momotok-Server/system"
	"Momotok-Server/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	id := -1

	if action_type == "" || to_user_id == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Incomplete information",
		})
		return
	}

	uid, err := utils.GetUID(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "invalid token",
		})
		return
	}

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	if action_type == "1-关注" {
		_ = db.QueryRow("select  id from concern where `concern_uid` = ? AND `concern_at_uid` = ?", uid, to_user_id).Scan(&id)
		if id != -1 {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "The condition is not true",
			})
			return
		}

		_, err = db.Exec("INSERT INTO concern (concern_uid, concern_at_uid) VALUES (?, ?)", uid, to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Concern failed",
			})
			return
		}
	} else if action_type == "2-取消关注" {
		_ = db.QueryRow("select  id from concern where `concern_uid` = ? AND `concern_at_uid` = ?", uid, to_user_id).Scan(&id)
		if id == -1 {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "The condition is not true",
			})
			return
		}

		_, err = db.Exec("DELETE FROM `momotok`.`concern` WHERE`concern_uid` = ? AND `concern_at_uid` = ?", uid, to_user_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Unfollow failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
	return
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
