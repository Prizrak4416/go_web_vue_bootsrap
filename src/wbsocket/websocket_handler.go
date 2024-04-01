package wbsocket

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	"work/src/comand"
)

type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (wsh *WebSocketHandler) Handle(c *gin.Context) {
	conn, err := wsh.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
	  log.Println("Ошибка при обновлении WebSocket:", err)
	  return
	}
	defer conn.Close()
  
	for {
	  _, msg, err := conn.ReadMessage()
	  if err != nil {
		log.Println("Ошибка чтения сообщения:", err)
		break
	  }
  
	  switch string(msg) {
	  case "get_journal":
		journal, err := comand.GetJournal()
		if err != nil {
		  log.Println("Ошибка при получении журнала:", err)
		  continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(journal)); err != nil {
		  log.Println("Ошибка при отправке журнала:", err)
		  continue
		}
	  default:
		conn.WriteMessage(websocket.TextMessage, msg)
	  }
	}
  }