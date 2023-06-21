package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"workshoptdd/config"
	"workshoptdd/entity"
	"workshoptdd/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	app *gin.Engine
)

func TestMain(m *testing.M) {
	db = config.InitDatabase("root:root@tcp(127.0.0.1:3306)/go_tdd_test?charset=utf8mb4&parseTime=True&loc=Local")
	app = routes.InitRoutes(db)

	m.Run()
	//	TODO cleanup database
}

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)

	app.ServeHTTP(w, req)

	response := w.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "OK", responseBody["message"])
}

func TestCreateTask(t *testing.T) {
	reqBody := strings.NewReader(`{
		"title": "example",
		"description": "description"
	}`)

	var beforeCount int64
	db.Find(&entity.Task{}).Count(&beforeCount)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", reqBody)

	app.ServeHTTP(w, req)
	response := w.Result()

	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]string
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "success create task", responseBody["message"])

	var afterCount int64
	db.Find(&entity.Task{}).Count(&afterCount)
	assert.Equal(t, afterCount, beforeCount+1)
}

func TestCreateTask_Fail(t *testing.T) {
	reqBody := strings.NewReader(`{
		"description": "description"
	}`)

	var beforeCount int64
	db.Find(&entity.Task{}).Count(&beforeCount)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", reqBody)

	app.ServeHTTP(w, req)
	response := w.Result()

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]string
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, "failed create task", responseBody["message"])

	var afterCount int64
	db.Find(&entity.Task{}).Count(&afterCount)
	assert.Equal(t, afterCount, beforeCount)
}

func TestGetTasks(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

	app.ServeHTTP(w, req)
	response := w.Result()

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)

	var responseBody []entity.Task
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		return
	}

	sampleTask := responseBody[0]
	if assert.NotEmpty(t, sampleTask) {
		assert.NotEmpty(t, sampleTask.Description)
		assert.NotEmpty(t, sampleTask.Title)
		assert.NotEmpty(t, sampleTask.ID)
	}
}

func TestDeleteTasks(t *testing.T) {
	dataToDelete := entity.Task{
		Title:       "delete title",
		Description: "delete description",
	}

	var beforeCount int64
	db.Find(&entity.Task{}).Count(&beforeCount)

	db.Create(&dataToDelete)
	url := "/tasks/" + strconv.Itoa(int(dataToDelete.ID))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, url, nil)

	app.ServeHTTP(w, req)
	response := w.Result()

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]string
	err := json.Unmarshal(body, &responseBody)
	if err != nil {
		return
	}
	assert.Equal(t, "success delete task", responseBody["message"])

	var afterCount int64
	db.Find(&entity.Task{}).Count(&afterCount)
	assert.Equal(t, beforeCount-1, afterCount)
}
