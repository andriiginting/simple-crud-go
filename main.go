package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/andriiginting/simple-crud-go/config"
	"github.com/andriiginting/simple-crud-go/domain"
	"github.com/andriiginting/simple-crud-go/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	foodRepo repository.FoodRepository
	err error
)

func main() {
	dbinfo := config.ConnectionString()
	db, err := sql.Open("postgres", dbinfo)
	checkError(err)

	foodRepo = repository.InitializeFoodRepository(db)
	defer db.Close()

	router := CreateRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/food/{id}", GetFood).Methods("GET")
	router.HandleFunc("/food", GetAllFoods).Methods("GET")
	router.HandleFunc("/food", InsertNewFood).Methods("POST")
	router.HandleFunc("/food/{id}", DeleteExistingFood).Methods("DELETE")
	router.HandleFunc("/food/{id}", UpdateExistingFoodPrice).Methods("PUT")
	return router
}

func GetFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	foodId, err := strconv.Atoi(params["id"])
	checkError(err)
	food := foodRepo.ReadFood(foodId)
	json.NewEncoder(w).Encode(food)
}

func InsertNewFood(w http.ResponseWriter, r *http.Request) {
	foodName := r.FormValue("name")
	foodPrice, err := strconv.Atoi(r.FormValue("price"))
	checkError(err)
	foodOwner := r.FormValue("owner")
	food := domain.Food{0, foodName, foodPrice, foodOwner}
	insertedFoodId := foodRepo.InsertFood(food)
	insertedFood := foodRepo.ReadFood(insertedFoodId)
	json.NewEncoder(w).Encode(insertedFood)
}

func DeleteExistingFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	foodId, err := strconv.Atoi(params["id"])
	checkError(err)

	deleteResponse := foodRepo.DeleteFood(foodId)
	successResponse := url.Values{}
	successResponse.Add("status", "200")
	failedResponse := url.Values{}
	failedResponse.Add("status", "400")

	if deleteResponse == 1 {
		json.NewEncoder(w).Encode(successResponse)
	} else {
		json.NewEncoder(w).Encode(failedResponse)
	}
}

func UpdateExistingFoodPrice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	foodId, err := strconv.Atoi(params["id"])
	checkError(err)

	foodPrice, err := strconv.Atoi(r.FormValue("price"))
	checkError(err)

	updatedResponse := foodRepo.UpdateFoodPrice(foodId, foodPrice)
	successResponse := url.Values{}
	successResponse.Add("status", "200")
	failedResponse := url.Values{}
	failedResponse.Add("status", "400")

	if updatedResponse == 1 {
		json.NewEncoder(w).Encode(successResponse)
	} else {
		json.NewEncoder(w).Encode(failedResponse)
	}
}

func GetAllFoods(w http.ResponseWriter, r *http.Request) {
	foods := foodRepo.GetAllFoods()
	json.NewEncoder(w).Encode(foods)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
