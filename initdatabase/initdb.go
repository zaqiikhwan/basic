package initdatabase

import (
	"main.go/Struct"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/project_1?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	err = db.AutoMigrate(&Struct.User{}, &Struct.Biodata{}, &Struct.Doctor{}, &Struct.Transaction{})
	if err != nil {
		return err
	}
	db.Exec("INSERT INTO biodata (id) VALUES (?)", 1)
	return nil
}