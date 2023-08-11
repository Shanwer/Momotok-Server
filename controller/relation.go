package controller

import (
	"Momotok-Server/rpc"
	"Momotok-Server/system"
	"Momotok-Server/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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

	if action_type == "关注" {
		err = db.QueryRow("select  id from followerlist where `follower_uid` = ? AND `followe_to_uid` = ?", uid, to_user_id).Scan(&id)
		if err != nil && err.Error() != "sql: no rows in result set" {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "The condition is not true",
			})
			return
		}

		tx, _ := db.Begin()
		_, err1 := tx.Exec("INSERT INTO followerlist (follower_uid, followe_to_uid) VALUES (?, ?)", uid, to_user_id)
		_, err2 := tx.Exec("UPDATE user SET follower_count = follower_count + 1 WHERE id = ?", to_user_id)
		_, err3 := tx.Exec("UPDATE user SET follow_count = follow_count + 1 WHERE id = ?", uid)
		if err1 != nil || err2 != nil || err3 != nil {
			tx.Rollback()

			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Concern failed",
			})
			return
		}
	} else if action_type == "取消关注" {
		err = db.QueryRow("select  id from followerlist where `follower_uid` = ? AND `followe_to_uid` = ?", uid, to_user_id).Scan(&id)
		if err != nil && err.Error() != "sql: no rows in result set" {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "The condition is not true",
			})
			return
		}

		tx, _ := db.Begin()
		_, err1 := tx.Exec("DELETE FROM followerlist WHERE`follower_uid` = ? AND `followe_to_uid` = ?", uid, to_user_id)
		_, err2 := tx.Exec("UPDATE user SET follower_count = follower_count - 1 WHERE id = ?", to_user_id)
		_, err3 := tx.Exec("UPDATE user SET follow_count = follow_count - 1 WHERE id = ?", uid)
		if err1 != nil || err2 != nil || err3 != nil {
			tx.Rollback()
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
	uid := c.Query("user_id")
	if !utils.CheckToken(c.Query("token")) {
		return
	}
	{
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "The condition is not true",
		})
		return
	}

	var userlist []User

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	follower_list, err := db.Query("select  followe_to_uid from followerlist where `follower_uid` = ?", uid)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Information cannot be obtained",
		})
		return
	}

	for follower_list.Next() {
		var followto int64

		err = follower_list.Scan(&followto)
		if err != nil {
			log.Fatal(err)
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

		userInfo := User{
			Id:              followto,
			Signature:       string(signature),
			Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
			BackgroundImage: "https://acg.suyanw.cn/api.php",
			IsFollow:        false,
			FollowerCount:   0,
			Name:            "",
		}

		temp_uid := 0
		err1 := db.QueryRow("SELECT followe_to_uid FROM followerlist WHERE follower_uid = ?", followto).Scan(&temp_uid)
		err2 := db.QueryRow("SELECT follow_count FROM user WHERE id = ?", followto).Scan(&userInfo.FollowCount)
		err3 := db.QueryRow("SELECT follower_count FROM user WHERE id = ?", followto).Scan(&userInfo.FollowerCount)
		err4 := db.QueryRow("SELECT username FROM user WHERE id = ?", followto).Scan(&userInfo.Name)
		err5 := db.QueryRow("SELECT total_likes FROM user WHERE id = ?", followto).Scan(&userInfo.TotalLikes)
		err6 := db.QueryRow("SELECT work_count FROM user WHERE id = ?", followto).Scan(&userInfo.WorkCount)
		err7 := db.QueryRow("SELECT total_received_likes FROM user WHERE id = ?", followto).Scan(&userInfo.TotalReceivedLikes)

		if (err1 != nil && err1.Error() != "sql: no rows in result set") || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Information cannot be obtained",
			})
			return
		}

		userlist = append(userlist, userInfo)

	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "Success",
		"user_list":   userlist,
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
