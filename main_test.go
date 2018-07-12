package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/andriiginting/simple-crud-go/config"
	"github.com/andriiginting/simple-crud-go/domain"
	"github.com/andriiginting/simple-crud-go/repository"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dbinfo := config.ConnectionString()
	db, err := sql.Open("postgres", dbinfo)
	checkError(err)

	foodRepo = repository.InitializeFoodRepository(db)
	defer db.Close()

	testResult := m.Run()

	os.Exit(testResult)
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
	food := domain.Food{4, "Spagetthi", 12000, "La Fonte"}
	lastInsertId := foodRepo.InsertFood(food)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	req, err := http.NewRequest("DELETE", "/food/"+strconv.Itoa(lastInsertId), nil)
	checkError(err)

	rr := httptest.NewRecorder()
	router := CreateRouter()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response should be 200 OK")
}

func TestUpdateExistingFoodPrice(t *testing.T) {
	food := domain.Food{4, "Spagetthi", 12000, "La Fonte"}
	lastInsertId := foodRepo.InsertFood(food)
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
