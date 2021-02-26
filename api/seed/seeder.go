package seed

import (
	"log"

	"github.com/dharlequin/restful-crud-golang-app/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Nickname: "Steven Victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

//Load loads initial values into DB
func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
