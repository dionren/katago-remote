package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/types"
	"server/utils"
)

type WebSocketServer struct {
	wsChan     chan types.MsgClient //websocket输出chain
	engineChan chan types.MsgServer //引擎输出chain
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		wsChan: make(chan types.MsgClient, 100),
	}
}

func (o *WebSocketServer) GetWebSocketOutChan() chan types.MsgClient {
	return o.wsChan
}

func (o *WebSocketServer) StartWs(engineChan chan types.MsgServer) {
	o.engineChan = engineChan
	http.HandleFunc("/", o.handleWebSocket)
}

// 把客户端的命令推送到引擎中执行
func (o *WebSocketServer) sendCommandToCmd(msgClient types.MsgClient) {
	log.Printf(msgClient.Cmd)
	o.wsChan <- msgClient
}

func (o *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	for {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading to websocket: %v", err)
			return
		}
		defer conn.Close()
		// 连接后先回复一个ready
		var msgServer = types.MsgServer{
			Str:      "GTP ready, beginning main protocol loop\n",
			Code:     1,
			Category: "gtp",
		}
		jsonData, err := json.Marshal(msgServer)
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Printf("Error writing to websocket: %v", err)
			return
		}
		// 从命令行输出通道读取并发送到 WebSocket 客户端
		go func() {
			for output := range o.engineChan {
				if output.Zip == 1 {
					output.Str, _ = utils.GzipBase64(output.Str)
				}
				jsonData, err := json.Marshal(output)
				err = conn.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					log.Printf("Error writing to websocket: %v", err)
					return
				}
			}
		}()

		// 从 WebSocket 客户端读取消息并发送到命令行输入通道
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading from websocket: %v", err)
				break
			}
			data := string(message)
			var msgClient types.MsgClient
			//data  转type.MsgServer
			err = json.Unmarshal([]byte(data), &msgClient)
			if err != nil {
				log.Fatal(err)
			}
			o.sendCommandToCmd(msgClient)
		}
	}
}
