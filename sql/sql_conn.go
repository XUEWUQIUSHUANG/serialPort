package sql

import (
	datastruct "go_code/serialPort/data_struct"
	parameter "go_code/serialPort/sql/sql_opt/parameter_opt"
	useropt "go_code/serialPort/sql/sql_opt/users_opt"

	"github.com/jinzhu/gorm"
)

var database *gorm.DB

func SqlInit(db *gorm.DB) {
	database = db
}

func Register(data datastruct.User) bool {
	return useropt.Register(database, data)
}

func Login(data datastruct.UserLogin) string {
	return useropt.Login(database, data)
}

func Change_Password(data datastruct.User) bool {
	return useropt.Change_Password(database, data)
}

func Check_Users(cookie string) bool {
	return useropt.Check_Users(database, cookie)
}

func Get_Parameter() []datastruct.EnvironmentData {
	return parameter.Get_Parameter(database)
}

func Update_Parameter(data datastruct.EnvironmentData) {
	parameter.Set_Parameter(database, data)
}
