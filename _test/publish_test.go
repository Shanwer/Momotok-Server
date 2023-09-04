package _test

import (
	"Momotok-Server/controller"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type Response1 struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
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
	router := gin.Default()

	router.POST("/douyin/publish/action/", controller.Publish)

	// Create a new multipart writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form field for the video file

	// Open the video file
	file, err := os.Open("video.avi")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// Create a buffer to store the file content
	fileContent := bytes.Buffer{}
	_, err = io.Copy(&fileContent, file)
	if err != nil {
		t.Fatal(err)
	}

	// Add other form fields
	_ = writer.WriteField("data", fileContent.String())
	_ = writer.WriteField("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM5ODYwMDYsInVzZXJpZCI6NiwidXNlcm5hbWUiOiJ0ZXN0In0.jwcpAC6mErbh3esZuYwS5qLRVxQTmQ3q9183RcTXbCI")
	_ = writer.WriteField("title", "video")

	// Close the multipart writer to finalize the body
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request with the body
	req, err := http.NewRequest("POST", "http://0.0.0.0:8080/douyin/publish/action/", body)
	if err != nil {
		t.Fatal(err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	var response Response1
	if res.Body.Len() > 0 {
		err := json.Unmarshal(res.Body.Bytes(), &response)
		if err != nil {
			t.Fatal(err)
		}
	}

	if response.StatusCode != 0 && res.Code == 200 {
		t.Errorf("Expected status code %d, but got %d", 0, response.StatusCode)
	}
}

func TestList(t *testing.T) {
	router := gin.Default()

	router.GET("/douyin/publish/list/", controller.PublishList)

	req, _ := http.NewRequest("GET", "http://0.0.0.0:8080/douyin/publish/list/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM5ODYwMDYsInVzZXJpZCI6NiwidXNlcm5hbWUiOiJ0ZXN0In0.jwcpAC6mErbh3esZuYwS5qLRVxQTmQ3q9183RcTXbCI&user_id=6", nil)
	req.Header.Set("Content-Type", "application/json")

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
