package main

import (
	"bufio"
	"client/types"
	"client/utils"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 1 {
		println("need wss url.")
		return
	}
	print(len(os.Args))
	print("args:", os.Args[1])

	wssUrl := os.Args[1]
	println("katavip-client v2.0.0")
	println("https://github.com/dionren/katavip-client")
	println(wssUrl)

	ws, _, err := websocket.DefaultDialer.Dial(wssUrl, nil)

	if err != nil {
		println("WebSocket connected error.")
		return
	} else {
		// 从WSS接收数据
		go func() {
			var msgServer types.MsgServer
			for {
				_, data, err := ws.ReadMessage()
				if err != nil {
					log.Fatal(err)
				}

				err = json.Unmarshal(data, &msgServer)

				if err != nil {
					return
				}

				if msgServer.Zip == 1 {
					msgServer.Str, _ = utils.UnGzipBase64(msgServer.Str)
				}

				_, err = os.Stdout.WriteString(msgServer.Str + "\n")
				if err != nil {
					return
				}

			}
		}()
	}

	var msgClient types.MsgClient

	// 发送压缩指令
	msgClient.Category = "ext"
	msgClient.Cmd = "zip"
	payload, _ := json.Marshal(msgClient)
	err = ws.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		log.Fatal(err)
	}

	// 从STD不断的轮询输入数据并通过WSS发送至服务器
	reader := bufio.NewReader(os.Stdin)
	for {
		byteArray, _, _ := reader.ReadLine()
		msgClient.Cmd = string(byteArray)
		msgClient.Category = "gtp"

		if msgClient.Cmd == "quit" {
			break
		}

		payload, _ := json.Marshal(msgClient)
		err = ws.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = ws.Close()
	if err != nil {
		log.Fatal(err)
	}
}
