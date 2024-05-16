package main

import (
	"net/http"
	"server/engine"
	"server/types"
	ws "server/websocket"
)

func main() {

	err := types.LoadConfig("config.ini")
	if err != nil {
		panic(err)
	}

	//创建引擎
	engineServer := engine.NewEngine()
	// 创建ws
	wsServer := ws.NewWebSocketServer()
	// 启动引擎
	engineServer.StartEngine(wsServer.GetWebSocketOutChan())
	// 启动ws程序
	wsServer.StartWs(engineServer.GetEngineOutChan())
	port, _ := types.GetConfigValue("websocket", "port")
	http.ListenAndServe(port, nil)

}
