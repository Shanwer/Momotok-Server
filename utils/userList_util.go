package utils

import (
	"Momotok-Server/model"
	"Momotok-Server/rpc"
	"Momotok-Server/system"
	"database/sql"
	"fmt"
	"io"
)

func GetUserList(uid int64) ([]model.User, error) {
	resp, _ := rpc.HttpRequest("GET", "https://v1.hitokoto.cn/?c=a&c=d&c=i&c=k&encode=text", nil)
	var userList []model.User
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)
	}
	signature, _ := io.ReadAll(resp.Body)

	userInfo := model.User{
		Id:              uid,
		Signature:       string(signature),
		Avatar:          "https://acg.suyanw.cn/sjtx/random.php",
		BackgroundImage: "https://acg.suyanw.cn/api.php",
		IsFollow:        false,
		FollowerCount:   0,
		Name:            "",
	}
	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	err = db.QueryRow("SELECT follow_count, follower_count, username, total_likes, work_count, total_received_likes FROM user WHERE id = ?", uid).Scan(&userInfo.FollowCount, &userInfo.FollowerCount, &userInfo.Name, &userInfo.FavoriteCount, &userInfo.WorkCount, &userInfo.TotalFavorited)

	if err != nil {
		return nil, err
	}

	userList = append(userList, userInfo)
	return userList, nil
}

func GetUserStruct(uid int64) (model.User, error) {
	var user model.User
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

	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress) //连接数据库
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return user, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT follow_count, follower_count, username, total_likes, work_count, total_received_likes FROM user WHERE id = ?", uid)
	err = row.Scan(&user.FollowCount, &user.FollowerCount, &user.Name, &user.FavoriteCount, &user.WorkCount, &user.TotalFavorited)
	user.Id = uid

	if err != nil {
		fmt.Println("Failed to scan row:", err)
		return user, err
	}
	return user, nil
}
