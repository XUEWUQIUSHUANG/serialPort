package websocket

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Server *websocket.Conn

func Init_listen(r *gin.Engine) {
	r.GET("/ws", func(context *gin.Context) {
		Server, _ = upgrader.Upgrade(context.Writer, context.Request, nil)
	})
}

func Listen(port io.ReadWriteCloser) {
	for {
		if Server != nil {
			_, msg, err := Server.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			port.Write(msg)
		}
	}
}
