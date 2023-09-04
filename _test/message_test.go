package _test

import (
	"Momotok-Server/controller"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type Response2 struct {
	StatusCode  int       `json:"status_code"`
	StatusMsg   string    `json:"status_msg"`
	MessageList []Message `json:"message_list"`
}

type Message struct {
	ID         int    `json:"id"`
	ToUserID   int    `json:"to_user_id"`
	FromUserID int    `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
}

func TestAction(t *testing.T) {
	router := gin.Default()

	router.POST("/douyin/message/action/", controller.MessageAction)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:8080/douyin/message/action/?to_user_id=1&action_type=1&content=hello&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	var response Response
	if res.Body.Len() > 0 {
		err := json.Unmarshal(res.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
	}

	if response.StatusCode != 0 {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}
}

func TestChat(t *testing.T) {
	router := gin.Default()

	router.GET("/douyin/message/chat/", controller.MessageChat)

	req, _ := http.NewRequest("GET", "http://0.0.0.0:8080/douyin/message/chat/?to_user_id=1&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	var response Response2
	if res.Body.Len() > 0 {
		err := json.Unmarshal(res.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}
	}

	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
	}

	if response.StatusCode != 0 {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}
}
