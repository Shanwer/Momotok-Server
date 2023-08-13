package controller

import (
	"Momotok-Server/model"
	"Momotok-Server/system"
	"Momotok-Server/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// usersLoginInfo use map to store user info, token is created by jwt
var usersLoginInfo = map[string]model.User{
	//TODO:其他模块需要改写为从数据库中获取
}

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	ip := c.ClientIP()
	hashedPassword, err := hashPassword(password)

	if err != nil {
		fmt.Println("Hash password failed：", err)
		return
	}
	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}
	_, err = db.Exec("INSERT INTO user (username, ip, password) VALUES (?, ?, ?)", username, ip, hashedPassword)
	if err != nil { //duplicate username check
		mysqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			fmt.Println("Register failed：", err)
		}
		if mysqlErr.Number == 1062 {
			c.JSON(http.StatusOK, UserResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: "User:" + username + "already exists!"},
			})
		} else {
			fmt.Println("Register failed：", err)
		}
		return
	}
	id := int64(0)
	err = db.QueryRow("SELECT ID FROM user WHERE username = ?", username).Scan(&id)
	token := utils.GenerateToken(username, id)

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	db, err := sql.Open("mysql", system.ServerInfo.Server.DatabaseAddress)
	if err != nil {
		fmt.Println("Database connected failed: ", err)
	}
	storedPassword := ""
	id := int64(0)
	err = db.QueryRow("SELECT password, id FROM user WHERE username = ?", username).Scan(&storedPassword, &id)
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		//fmt.Println("Wrong username or password: ", err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "Wrong username or password"},
		})
		return
	}
	token := utils.GenerateToken(username, id)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   id,
		Token:    token,
	})
}

func UserInfo(c *gin.Context) {
	uid := c.Query("user_id")
	if !utils.CheckToken(c.Query("token")) {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
		return
	}
	id, _ := strconv.ParseInt(uid, 10, 64)
	user, err := utils.GetUserStruct(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "Success",
		"user":        user,
	})
	return
}

// generate hashed password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
