package mysqlexample

import (
	"fmt"
	"ginapi/structs"

	tokentool "ginapi/utils/token"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initdb() {
	dsn := "ellis:ellis@tcp(192.168.214.134:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	CreateTable(&structs.User{})
}

func CreateTable(models interface{}) {
	db.AutoMigrate(models)
}

func Create(value interface{}) bool {
	d := db.Create(value)
	err := d.Error
	if err != nil {
		return false
	} else {
		fmt.Printf("d.RowsAffected: %v\n", d.RowsAffected)
		return true
	}
}

func CheckUserExist(username string, models interface{}) bool {
	d := db.Where("name=?", username).Limit(1).Find(models)
	if d.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func LoginCheck(username string) (bool, string) {
	result := false
	user := structs.User{}
	d := db.Where("name=?", username).Limit(1).Find(&user)
	if d.RowsAffected > 0 {
		result = true
	} else {
		result = false
	}

	if result {
		token, err := tokentool.GenerateToken(user.ID)
		if err != nil {
			return false, ""
		} else {
			return true, token
		}

	} else {
		return false, ""
	}
}

func GetUserByID(id int, models interface{}) {
	db.Find(models, id)
}
