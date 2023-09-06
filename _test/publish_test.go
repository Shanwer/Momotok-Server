package _test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

type Response1 struct {
	StatusCode int `json:"status_code"`
	//StatusMsg  string `json:"status_msg"`
}

type Response2 struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	NextTime   int    `json:"next_time"`
	VideoList  []struct {
		ID     int `json:"id"`
		Author struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			FollowCount     int    `json:"follow_count"`
			FollowerCount   int    `json:"follower_count"`
			IsFollow        bool   `json:"is_follow"`
			Avatar          string `json:"avatar"`
			BackgroundImage string `json:"background_image"`
			Signature       string `json:"signature"`
			TotalFavorited  string `json:"total_favorited"`
			WorkCount       int    `json:"work_count"`
			FavoriteCount   int    `json:"favorite_count"`
		} `json:"author"`
		PlayURL       string `json:"play_url"`
		CoverURL      string `json:"cover_url"`
		FavoriteCount int    `json:"favorite_count"`
		CommentCount  int    `json:"comment_count"`
		IsFavorite    bool   `json:"is_favorite"`
		Title         string `json:"title"`
	} `json:"video_list"`
}

func TestAction(t *testing.T) {

	// 打开视频文件
	file, err := os.Open("bear.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加 token 字段到多部分表单
	tokenField, err := writer.CreateFormField("token")
	if err != nil {
		log.Fatal(err)
	}
	tokenField.Write([]byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQxNDI3NzcsInVzZXJpZCI6OSwidXNlcm5hbWUiOiJ0ZXN0In0.B_sZ6EaoR76irdZwR2vsaazShyPlGP30KZhJeCYSesM"))

	titleField, err := writer.CreateFormField("title")
	if err != nil {
		log.Fatal(err)
	}
	titleField.Write([]byte("video"))

	dataField, err := writer.CreateFormFile("data", "bear.mp4")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(dataField, file) // 将文件内容拷贝到表单字段中
	if err != nil {
		log.Fatal(err)
	}

	// 必须在writer.Close()之前调用，以便写入最后的boundary
	writer.Close()

	// 创建一个新的HTTP请求，并将body设置为multipart.Writer的内容
	res, err := http.Post("http://0.0.0.0:8080/douyin/publish/action/", writer.FormDataContentType(), body)
	assert.NoError(t, err)
	var response Response1
	returnedJson, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	err = json.Unmarshal(returnedJson, &response)
	assert.Nil(t, err)
	if response.StatusCode != 0 {
		t.Errorf("Expected status code %d, but got %d", 0, response.StatusCode)
	}
}

func TestList(t *testing.T) {
	res, err := http.Get("http://localhost:8080/douyin/publish/list/?user_id=9&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQxMDA5NzUsInVzZXJpZCI6OSwidXNlcm5hbWUiOiJ0ZXN0In0.16RpXZQAGYgy5XGULMRlEufFq_ruPlsGjpQYxRT29DY")
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
