package parameteropt

import (
	"fmt"
	datastruct "go_code/serialPort/data_struct"
	"io"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

var multiWriter io.Writer

func init() {
	logPath := "./log/parameter/" + time.Now().Format("2006-01-02_15-04-05") + ".log"
	f, err := os.Create(logPath)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	multiWriter = io.MultiWriter(f)
}

func Get_Parameter(db *gorm.DB) []datastruct.EnvironmentData {
	var parameters []datastruct.EnvironmentData
	result := db.Find(&parameters)
	if result.Error != nil {
		fmt.Println("Error:", result.Error)
		return nil
	} else {
		return parameters
	}
}

func Set_Parameter(db *gorm.DB, parameter datastruct.EnvironmentData) {
	result := db.Create(&parameter)
	if result.Error != nil {
		fmt.Fprintln(multiWriter, "Set_Parameter Error:", result.Error, time.Now().Format("2006-01-02 15:04:05"))
	} else {
		fmt.Fprintln(multiWriter, "Set_Parameter Success", time.Now().Format("2006-01-02 15:04:05"))
	}
}
