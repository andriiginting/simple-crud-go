package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/andriiginting/simple-crud-go/domain"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	POSTGRE_USERNAME_KEY = "PGUSERNAME"
	POSTGRE_PORT_KEY     = "PGPORT"
	POSTGRE_PASSWORD_KEY = "PGPASSWORD"
	POSTGRE_DB_NAME_KEY  = "DBNAME"
)

var (
	db  *sql.DB
	err error
)

func main() {
	viperConfig := GetConfig()
	username, port, _, dbname := GetParsedConfig(viperConfig)
	dbinfo := fmt.Sprintf("user=%s port=%d dbname=%s sslmode=disable", username, port, dbname)
	db, err = sql.Open("postgres", dbinfo)
	checkError(err)
	defer db.Close()

	router := CreateRouter()
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigFile("application.yml")
	err := config.ReadInConfig()
	checkError(err)
	return config
}

func GetParsedConfig(viper *viper.Viper) (string, int, string, string) {
	deployEnv := os.Getenv("DEPLOYENV")
	return viper.GetString(deployEnv + "." + POSTGRE_USERNAME_KEY), viper.GetInt(deployEnv + "." + POSTGRE_PORT_KEY), viper.GetString(deployEnv + "." + POSTGRE_PASSWORD_KEY), viper.GetString(deployEnv + "." + POSTGRE_DB_NAME_KEY)
}

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/food/{id}", GetFood).Methods("GET")
	router.HandleFunc("/food", GetAllFood).Methods("GET")
	router.HandleFunc("/food", InsertNewFood).Methods("POST")
	router.HandleFunc("/food/{id}", DeleteExistingFood).Methods("DELETE")
	router.HandleFunc("/food/{id}", UpdateExistingFoodPrice).Methods("PUT")
	return router
}

func readFood(id int) *domain.Food {
	row := db.QueryRow("SELECT * FROM food WHERE ID=$1", id)

	var foodId int
	var name string
	var price int
	var owner string

	err = row.Scan(&foodId, &name, &price, &owner)
	if err == sql.ErrNoRows {
		return nil
	} else {
		checkError(err)
		return &domain.Food{foodId, name, price, owner}
	}
}

func insertFood(food domain.Food) int {
	var lastInsertId int
	err = db.QueryRow("INSERT INTO food(name,price,owner) VALUES($1,$2,$3) returning id;", food.Name, food.Price, food.Owner).Scan(&lastInsertId)
	checkError(err)
	return lastInsertId
}

func deleteFood(foodId int) int {
	stmt, err := db.Prepare("DELETE FROM food where id=$1")
	checkError(err)

	result, err := stmt.Exec(foodId)
	checkError(err)

	affectedRows, err := result.RowsAffected()
	checkError(err)

	return int(affectedRows)
}

func updateFoodPrice(id int, price int) int {
	stmt, err := db.Prepare("UPDATE food set price=$1 where id=$2")
	checkError(err)

	res, err := stmt.Exec(price, id)
	checkError(err)

	affectedRows, err := res.RowsAffected()
	checkError(err)

	return int(affectedRows)
}

func GetFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	foodId, err := strconv.Atoi(params["id"])
	checkError(err)
	food := readFood(foodId)
	json.NewEncoder(w).Encode(food)
}

func InsertNewFood(w http.ResponseWriter, r *http.Request) {
	foodName := r.FormValue("name")
	foodPrice, err := strconv.Atoi(r.FormValue("price"))
	checkError(err)
	foodOwner := r.FormValue("owner")
	food := domain.Food{0, foodName, foodPrice, foodOwner}
	insertedFoodId := insertFood(food)
	insertedFood := readFood(insertedFoodId)
	json.NewEncoder(w).Encode(insertedFood)
}

func DeleteExistingFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	foodId, err := strconv.Atoi(params["id"])
	checkError(err)

	deleteResponse := deleteFood(foodId)
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

	updatedResponse := updateFoodPrice(foodId, foodPrice)
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

func GetAllFood(w http.ResponseWriter, r *http.Request) {
	var foods []domain.Food

	rows, err := db.Query("SELECT * FROM food order by id")
	checkError(err)
	for rows.Next() {
		var id int
		var name string
		var price int
		var owner string
		err = rows.Scan(&id, &name, &price, &owner)
		checkError(err)

		foods = append(foods, domain.Food{id, name, price, owner})
	}
	json.NewEncoder(w).Encode(foods)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
