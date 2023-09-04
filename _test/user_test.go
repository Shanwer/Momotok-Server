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
	UserID     int    `json:"user_id"`
	Token      string `json:"token"`
}

func TestRegister(t *testing.T) {
	router := gin.Default()

	router.POST("/douyin/user/register", controller.Register)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:8080/douyin/user/register?username=test&password=test", nil)
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

	if response.StatusCode != 0 {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.StatusCode)
	}
}

func TestLogin(t *testing.T) {
	t.Helper()

	router := gin.Default()

	router.POST("/douyin/user/login", controller.Login)

	req, _ := http.NewRequest("POST", "http://0.0.0.0:8080/douyin/user/login?username=test&password=test", nil)
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
		t.Errorf("Expected status code %d, but got %d", 0, response.StatusCode)
	}
}

func TestUser(t *testing.T) {
	router := gin.Default()

	router.GET("/douyin/user/", controller.UserInfo)

	req, _ := http.NewRequest("GET", "http://0.0.0.0:8080/douyin/user/?user_id=6&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTM5MjYzOTIsInVzZXJpZCI6NiwidXNlcm5hbWUiOiJ0ZXN0In0.7w-mygkH6JUjtsWD9GNYyjYn0LGYaAtAJKEKpNafu14", nil)
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
		t.Errorf("Expected status code %d, but got %d", 0, response.StatusCode)
	}
}
