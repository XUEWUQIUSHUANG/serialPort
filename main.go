package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	datastruct "go_code/serialPort/data_struct"
	"go_code/serialPort/router"
	"go_code/serialPort/serial"
	"go_code/serialPort/sql"
	"go_code/serialPort/websocket"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var config datastruct.Config

func init() {
	filePath := "./config.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件时出错:", err)
		return
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		fmt.Println("解码 JSON 数据时出错:", err)
		return
	}
}

func main() {

	// mysql
	fmt.Println("Connecting to database...")
	db, err := gorm.Open("mysql", "root:"+config.Password+"@("+config.IP+":"+config.Port+")/"+config.Database+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database.")
	sql.SqlInit(db)
	defer db.Close()

	// serial
	serial.Init_serial()
	defer serial.Port.Close()

	// log
	logPath := "./log/gin/" + time.Now().Local().GoString() + ".log"
	f, _ := os.Create(logPath)
	gin.DefaultWriter = io.MultiWriter(f)

	// router
	r := gin.Default()
	r.Use(staticFileMiddleware("./dist"))
	router.Init_listen(r)

	//ws
	websocket.Init_listen(r)
	r.Run(":3000")

	select {}

}

func staticFileMiddleware(staticPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := filepath.Join(staticPath, c.Request.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			c.File(filePath)
			c.Abort()
		} else {
			c.Next()
		}
	}
}
