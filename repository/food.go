package repository

import (
	"database/sql"
	"log"

	"github.com/andriiginting/simple-crud-go/domain"
	_ "github.com/lib/pq"
)

var (
	err error
)

type FoodRepository struct {
	db *sql.DB
}

func (self FoodRepository) ReadFood(id int) *domain.Food {
	row := self.db.QueryRow("SELECT * FROM food WHERE ID=$1", id)

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

func (self FoodRepository) InsertFood(food domain.Food) int {
	var lastInsertId int
	err = self.db.QueryRow("INSERT INTO food(name,price,owner) VALUES($1,$2,$3) returning id;", food.Name, food.Price, food.Owner).Scan(&lastInsertId)
	checkError(err)
	return lastInsertId
}

func (self FoodRepository) DeleteFood(foodId int) int {
	stmt, err := self.db.Prepare("DELETE FROM food where id=$1")
	checkError(err)

	result, err := stmt.Exec(foodId)
	checkError(err)

	affectedRows, err := result.RowsAffected()
	checkError(err)

	return int(affectedRows)
}

func (self FoodRepository) UpdateFoodPrice(id int, price int) int {
	stmt, err := self.db.Prepare("UPDATE food set price=$1 where id=$2")
	checkError(err)

	res, err := stmt.Exec(price, id)
	checkError(err)

	affectedRows, err := res.RowsAffected()
	checkError(err)

	return int(affectedRows)
}

func InitializeFoodRepository(db *sql.DB) FoodRepository{
	return FoodRepository{db}
}

func (self FoodRepository) GetDB() *sql.DB {
	return self.db
}

func checkError(err error) {
	if (err != nil) {
		log.Fatalf(err.Error())
	}
}
