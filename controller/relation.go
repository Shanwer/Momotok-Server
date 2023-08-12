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
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction handles follow and unfollow action
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	id := -1 //id variable is used for checks that prevent some special cases like unfollow an unfollowed person

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

	if strconv.FormatInt(uid, 10) == to_user_id {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "You cannot follow yourself",
		})
		return
	}

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	if action_type == "1" {
		err = db.QueryRow("select id from follow_list where follower_uid = ? AND following_uid = ?", uid, to_user_id).Scan(&id)
		if id != -1 {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "You have already followed the user!",
			})
			return
		}

		tx, _ := db.Begin()
		_, err1 := tx.Exec("INSERT INTO follow_list (follower_uid, following_uid) VALUES (?, ?)", uid, to_user_id)
		_, err2 := tx.Exec("UPDATE user SET follower_count = follower_count + 1 WHERE id = ?", to_user_id)
		_, err3 := tx.Exec("UPDATE user SET follow_count = follow_count + 1 WHERE id = ?", uid)
		if err1 != nil || err2 != nil || err3 != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Follow failed",
			})
			return
		}
		err := tx.Commit()
		if err != nil {
			return
		}
	} else if action_type == "2" {
		err = db.QueryRow("select id from follow_list where follower_uid = ? AND following_uid = ?", uid, to_user_id).Scan(&id)
		if err != nil && err.Error() != "sql: no rows in result set" {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "The condition is not true",
			})
			return
		}

		tx, _ := db.Begin()
		_, err1 := tx.Exec("DELETE FROM follow_list WHERE follower_uid = ? AND following_uid = ?", uid, to_user_id)
		_, err2 := tx.Exec("UPDATE user SET follower_count = follower_count - 1 WHERE id = ?", to_user_id)
		_, err3 := tx.Exec("UPDATE user SET follow_count = follow_count - 1 WHERE id = ?", uid)
		if err1 != nil || err2 != nil || err3 != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Unfollow failed",
			})
			return
		}
		err := tx.Commit()
		if err != nil {
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
	return
}

// FollowList provides user with follow list
func FollowList(c *gin.Context) {
	uid := c.Query("user_id")
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
		return
	}
	var userList []User

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	followingList, err := db.Query("select following_uid from follow_list where follower_uid = ?", uid)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Information cannot be obtained",
		})
		return
	}

	for followingList.Next() {
		var followingToUID int64 //user's id who is followed by requested user

		err = followingList.Scan(&followingToUID)
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
			Id:              followingToUID,
			Signature:       string(signature),
			Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
			BackgroundImage: "https://acg.suyanw.cn/api.php",
			IsFollow:        false,
			FollowerCount:   0,
			Name:            "",
		}

		err = db.QueryRow("SELECT follow_count, follower_count, username, total_likes, work_count, total_received_likes FROM user WHERE id = ?", followingToUID).Scan(&userInfo.FollowCount, &userInfo.FollowerCount, &userInfo.Name, &userInfo.TotalLikes, &userInfo.WorkCount, &userInfo.TotalReceivedLikes)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Information cannot be obtained",
			})
			return
		}
		userList = append(userList, userInfo)

	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "Success",
		"user_list":   userList,
	})
	return
}

// FollowerList provides the user with follower list
func FollowerList(c *gin.Context) {
	uid := c.Query("user_id")
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
		return
	}
	var userList []User

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	followerList, err := db.Query("select follower_uid from follow_list where following_uid = ?", uid)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Information cannot be obtained",
		})
		return
	}

	for followerList.Next() {
		var followerUID int64 //user's id who is following the requested user

		err = followerList.Scan(&followerUID)
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
			Id:              followerUID,
			Signature:       string(signature),
			Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
			BackgroundImage: "https://acg.suyanw.cn/api.php",
			IsFollow:        false,
			FollowerCount:   0,
			Name:            "",
		}

		err = db.QueryRow("SELECT follow_count, follower_count, username, total_likes, work_count, total_received_likes FROM user WHERE id = ?", followerUID).Scan(&userInfo.FollowCount, &userInfo.FollowerCount, &userInfo.Name, &userInfo.TotalLikes, &userInfo.WorkCount, &userInfo.TotalReceivedLikes)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Information cannot be obtained",
			})
			return
		}
		userList = append(userList, userInfo)

	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "Success",
		"user_list":   userList,
	})
	return
}

// FriendList provides the user with friend list
func FriendList(c *gin.Context) {
	uid := c.Query("user_id")
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
		return
	}
	var userList []User

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}

	friendList, err := db.Query("select t1.following_uid from follow_list as t1 join follow_list as t2 on t1.following_uid = t2.follower_uid and t1.follower_uid = t2.following_uid where t1.follower_uid = ?", uid)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Information cannot be obtained",
		})
		return
	}

	for friendList.Next() {
		var friendUID int64 //user's id who is followed by requested user

		err = friendList.Scan(&friendUID)
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
			Id:              friendUID,
			Signature:       string(signature),
			Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
			BackgroundImage: "https://acg.suyanw.cn/api.php",
			IsFollow:        false,
			FollowerCount:   0,
			Name:            "",
		}

		err = db.QueryRow("SELECT follow_count, follower_count, username, total_likes, work_count, total_received_likes FROM user WHERE id = ?", friendUID).Scan(&userInfo.FollowCount, &userInfo.FollowerCount, &userInfo.Name, &userInfo.TotalLikes, &userInfo.WorkCount, &userInfo.TotalReceivedLikes)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Information cannot be obtained",
			})
			return
		}
		userList = append(userList, userInfo)

	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "Success",
		"user_list":   userList,
	})
	return
}
