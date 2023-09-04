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

	router.POST("/douyin/favorite/action/", controller.FavoriteAction)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:8080/douyin/favorite/action/?video_id=1&action_type=1&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
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

func TestList(t *testing.T) {
	router := gin.Default()

	router.GET("/douyin/favorite/list/", controller.FavoriteAction)

	req, _ := http.NewRequest("GET", "http://0.0.0.0:8080/douyin/favorite/list/?user_id=1&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", nil)
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
