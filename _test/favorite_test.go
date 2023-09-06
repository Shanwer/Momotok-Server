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
	res, err := http.Post("http://localhost:8080/douyin/favorite/action/?video_id=1&action_type=1&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQxMDA5NzUsInVzZXJpZCI6OSwidXNlcm5hbWUiOiJ0ZXN0In0.16RpXZQAGYgy5XGULMRlEufFq_ruPlsGjpQYxRT29DY", "POST", nil)
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

func TestList(t *testing.T) {
	res, err := http.Get("http://localhost:8080/douyin/favorite/list/?user_id=9&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQxMDA5NzUsInVzZXJpZCI6OSwidXNlcm5hbWUiOiJ0ZXN0In0.16RpXZQAGYgy5XGULMRlEufFq_ruPlsGjpQYxRT29DY")
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
