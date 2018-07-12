package repository

import (
	"database/sql"
	"os"
	"testing"

	"github.com/andriiginting/simple-crud-go/config"
	"github.com/andriiginting/simple-crud-go/domain"
	"github.com/stretchr/testify/assert"
)

var (
	foodRepo FoodRepository
)

func TestMain(m *testing.M) {
	dbinfo := config.ConnectionString()
	db, err := sql.Open("postgres", dbinfo)
	checkError(err)

	foodRepo = InitializeFoodRepository(db)

	defer db.Close()

	testResult := m.Run()

	os.Exit(testResult)
}

func TestReadFood(t *testing.T) {
	food := foodRepo.ReadFood(1)
	assert.NotNil(t, food)
	assert.Equal(t, 1, food.Id, "Food ID 1 should be equal")
	food = foodRepo.ReadFood(3)
	assert.NotNil(t, food)
	assert.Equal(t, 3, food.Id, "Food ID 3 should be equal")
	food = foodRepo.ReadFood(-1)
	assert.Nil(t, food)
}

func TestInsertFood(t *testing.T) {
	food := domain.Food{4, "Spagetthi", 12000, "La Fonte"}
	lastInsertId := foodRepo.InsertFood(food)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	insertedFood := foodRepo.ReadFood(lastInsertId)
	assert.NotNil(t, insertedFood)
}

func TestDeleteFood(t *testing.T) {
	foodInserted := domain.Food{0, "Kopi Aku Kamu", 18000, "Aku Kamu"}
	lastInsertId := foodRepo.InsertFood(foodInserted)
	assert.NotEqual(t, 0, lastInsertId, "Food should have been inserted")

	affectedRows := foodRepo.DeleteFood(lastInsertId)
	assert.Equal(t, 1, affectedRows, "There should be only 1 food affected")
}

func TestUpdateFoodPrice(t *testing.T) {
	affectedRows := foodRepo.UpdateFoodPrice(1, 28000)
	assert.Equal(t, 1, affectedRows, "There should be only 1 food affected")
	food := foodRepo.ReadFood(1)
	assert.Equal(t, 1, food.Id, "Effects of update should have applied only to food ID 1")
	assert.Equal(t, 28000, food.Price, "Food price should have been updated")
}
