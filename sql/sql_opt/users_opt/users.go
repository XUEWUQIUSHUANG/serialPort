package useropt

import (
	"fmt"
	datastruct "go_code/serialPort/data_struct"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Login(db *gorm.DB, data datastruct.UserLogin) string {

	var user datastruct.User
	result := db.Where("username = ? and password = ?", data.Username, data.Password).First(&user)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return "error"
	} else {
		cookie := uuid.New().String()
		db.Model(&user).Where("username = ? and password = ?", data.Username, data.Password).Update("cookie", cookie)
		return cookie
	}
}

func Register(db *gorm.DB, data datastruct.User) bool {
	if err := db.Create(&data).Error; err != nil {
		return false
	} else {
		return true
	}
}

func Check_Users(db *gorm.DB, cookie string) bool {
	var user datastruct.User
	result := db.Where("cookie = ?", cookie).First(&user)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	} else {
		return true
	}
}

func Change_Password(db *gorm.DB, data datastruct.User) bool {
	var user datastruct.User
	cookie := uuid.New().String()
	result :=db.Model(&user).Where("username = ?", data.Username).Updates(datastruct.User{Password: data.Password, Cookie: cookie})
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return false
	} else {
		return true
	}
}
