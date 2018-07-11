package main

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	viperConfig := GetConfig()
	username, port, _, dbname := GetParsedConfig(viperConfig)
	dbinfo := fmt.Sprintf("user=%s port=%d dbname=%s sslmode=disable", username, port, dbname)
	db, err = sql.Open("postgres", dbinfo)
	checkError(err)
	defer db.Close()

	testResult := m.Run()

	os.Exit(testResult)
}

func TestReadFood(t *testing.T) {
	food := readFood(1)
	assert.NotNil(t, food)
	assert.Equal(t, 1, food.Id, "Food ID 1 should be equal")
	food = readFood(3)
	assert.NotNil(t, food)
	assert.Equal(t, 3, food.Id, "Food ID 3 should be equal")
	food = readFood(-1)
	assert.Nil(t, food)
}

func TestInsertFood(t *testing.T) {
	food := Food{4, "Spagetthi", 12000, "La Fonte"}
	lastInsertId := insertFood(food)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	insertedFood := readFood(lastInsertId)
	assert.NotNil(t, insertedFood)
}

func TestDeleteFood(t *testing.T) {
	foodInserted := Food{0, "Kopi Aku Kamu", 18000, "Aku Kamu"}
	lastInsertId := insertFood(foodInserted)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	affectedRows := deleteFood(lastInsertId)
	assert.Equal(t, 1, affectedRows, "There should be only 1 food affected")
}

func TestUpdateFoodPrice(t *testing.T) {
	affectedRows := updateFoodPrice(1, 28000)
	assert.Equal(t, 1, affectedRows, "There should be only 1 food affected")
	food := readFood(1)
	assert.Equal(t, 1, food.Id, "Effects of update should have applied only to food ID 1")
	assert.Equal(t, 28000, food.Price, "Food price should have been updated")
}

func TestGetFood(t *testing.T) {
	req, err := http.NewRequest("GET", "/food/1", nil)
	checkError(err)

	reqAllFood, err := http.NewRequest("GET", "/food", nil)
	checkError(err)

	rr := httptest.NewRecorder()
	rrAllFood := httptest.NewRecorder()
	router := CreateRouter()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response get food should be 200 OK")

	router.ServeHTTP(rrAllFood, reqAllFood)
	assert.Equal(t, http.StatusOK, rrAllFood.Code, "Response for get all food should be 200 OK")
}

func TestInsertNewFood(t *testing.T) {
	foodData := url.Values{}
	foodData.Set("name", "Latte")
	foodData.Set("price", "10000")
	foodData.Set("owner", "Jumpstart")

	req, err := http.NewRequest("POST", "/food", strings.NewReader(foodData.Encode()))
	checkError(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	router := CreateRouter()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response should be 200 OK")
}

func TestDeleteExistingFood(t *testing.T) {
	food := Food{4, "Spagetthi", 12000, "La Fonte"}
	lastInsertId := insertFood(food)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	req, err := http.NewRequest("DELETE", "/food/"+strconv.Itoa(lastInsertId), nil)
	checkError(err)

	rr := httptest.NewRecorder()
	router := CreateRouter()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response should be 200 OK")
}

func TestUpdateExistingFoodPrice(t *testing.T) {
	food := Food{4, "Spagetthi", 12000, "La Fonte"}
	lastInsertId := insertFood(food)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	foodPrice := url.Values{}
	foodPrice.Set("price", "30000")

	req, err := http.NewRequest("PUT", "/food/"+strconv.Itoa(lastInsertId), strings.NewReader(foodPrice.Encode()))
	checkError(err)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	router := CreateRouter()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response should be 200 OK")
}

func TestGetAllFood(t *testing.T) {
	req, err := http.NewRequest("GET", "/food", nil)
	checkError(err)

	rr := httptest.NewRecorder()
	router := CreateRouter()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response should be 200 OK")
}
