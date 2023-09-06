package _test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
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
	res, err := http.Post("http://localhost:8080/douyin/message/action?to_user_id=1&action_type=1&content=hello&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQxMDA5NzUsInVzZXJpZCI6OSwidXNlcm5hbWUiOiJ0ZXN0In0.16RpXZQAGYgy5XGULMRlEufFq_ruPlsGjpQYxRT29DY", "POST", nil)
	assert.NoError(t, err)
	var response Response
	returnedJson, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(returnedJson, &response)
	assert.Nil(t, err)
	if response.StatusCode != 0 {
		t.Errorf("Expected status code %d, but got %d", 0, response.StatusCode)
	}
}

func TestChat(t *testing.T) {
	res, err := http.Get("http://localhost:8080/douyin/message/chat?to_user_id=1&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQxMDA5NzUsInVzZXJpZCI6OSwidXNlcm5hbWUiOiJ0ZXN0In0.16RpXZQAGYgy5XGULMRlEufFq_ruPlsGjpQYxRT29DY")
	assert.NoError(t, err)
	var response Response2
	returnedJson, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(returnedJson, &response)
	assert.Nil(t, err)
	if response.StatusCode != 0 {
		t.Errorf("Expected status code %d, but got %d", 0, response.StatusCode)
	}
}
