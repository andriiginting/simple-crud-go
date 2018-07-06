package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Food struct {
	id    int
	name  string
	price int
	owner string
}

const (
	DB_USERNAME = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "food_test"
)

var (
	db  *sql.DB
	err error
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USERNAME, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)
	checkError(err)
	defer db.Close()
}

func readFood(id int) *Food {
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
		return &Food{foodId, name, price, owner}
	}
}

func insertFood(food Food) int {
	var lastInsertId int
	err = db.QueryRow("INSERT INTO food(name,price,owner) VALUES($1,$2,$3) returning id;", food.name, food.price, food.owner).Scan(&lastInsertId)
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

	res, err := stmt.Exec(price,id)
	checkError(err)

	affectedRows, err := res.RowsAffected()
	checkError(err)

	return int(affectedRows)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}