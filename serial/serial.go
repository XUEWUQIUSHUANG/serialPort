package serial

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	datastruct "go_code/serialPort/data_struct"
	"go_code/serialPort/sql"
	ws "go_code/serialPort/websocket"

	"github.com/gorilla/websocket"
	"github.com/jacobsa/go-serial/serial"
)

var Port io.ReadWriteCloser
var err error
var DataChan chan []byte

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func Init_serial() {
	DataChan = make(chan []byte, 2000)
	options := serial.OpenOptions{
		PortName:        "COM2",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	Port, err = serial.Open(options)
	if err != nil {
		log.Fatalf("Failed to open port: %v", err)
	}
	log.Println("Serial port opened successfully!")
	go listen(Port)
	go processData()
	go ws.Listen(Port)
}

func listen(port io.ReadWriteCloser) {
	buf := make([]byte, 4096)
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Println("Error reading from serial port:", err)
			return
		}
		if n > 0 {
			DataChan <- buf[:n]
		}
		time.Sleep(1 * time.Second)
	}
}

func processData() {

	var jsonBuffer bytes.Buffer

	for data := range DataChan {
		for _, b := range data {
			if b == '{' && jsonBuffer.Len() == 0 {
				jsonBuffer.WriteByte(b)
			} else if b == '}' && jsonBuffer.Len() > 0 {
				jsonBuffer.WriteByte(b)
				var data datastruct.EnvironmentData
				err := json.Unmarshal(jsonBuffer.Bytes(), &data)
				data.RecordedAt = time.Now()
				if err != nil {
					fmt.Println("Error unmarshalling JSON:", err)
				} else {
					sql.Update_Parameter(data)
					Try(func() {
						if ws.Server != nil {
							jsonData, err := json.Marshal(data)
							if err != nil {
								log.Println("Error while marshaling JSON:", err)
								return
							}
							ws.Server.WriteMessage(websocket.TextMessage, jsonData)
						}
					}, func(err interface{}) {
						fmt.Println(err)
					})
				}
				jsonBuffer.Reset()
			} else if jsonBuffer.Len() > 0 {
				jsonBuffer.WriteByte(b)
			}
		}
	}

}
