package engine

import (
	"bufio"
	"log"
	"os/exec"
	"server/types"
	"sync"
)

type Engine struct {
	engineChan       chan types.MsgServer // 命令行输出
	wsChan           chan types.MsgClient // 输出信息的队列
	MGtpIndex        int                  // MGtpIndex 正在执行的GTP的index, 主要用于野狐
	MOperatorOffline int64                // MOperatorOffline operator offline，操作端断开的时间，用于处理长时间不用自动退租，null表示从未连过，或者连接正常
	MZip             int                  // MZip 是否压缩
	cmd              *exec.Cmd            // Command instance to control the running process
	mu               sync.Mutex           // Mutex to protect concurrent access to shared variables
}

func NewEngine() *Engine {
	return &Engine{
		engineChan:       make(chan types.MsgServer, 100),
		MGtpIndex:        0,
		MOperatorOffline: 0,
		MZip:             0,
	}
}
func (w *Engine) GetEngineOutChan() chan types.MsgServer {
	return w.engineChan
}

func (w *Engine) StartEngine(outChan chan types.MsgClient) {
	w.wsChan = outChan
	go w.runCommand()
}

func (w *Engine) runCommand() {
	w.mu.Lock()
	katago, err := types.GetConfigValue("engine", "katago")
	gtpConfig, err := types.GetConfigValue("engine", "gtpConfig")
	model, err := types.GetConfigValue("engine", "model")

	w.cmd = exec.Command(katago, "gtp", "-config", gtpConfig, "-model", model)
	w.mu.Unlock()
	stdout, err := w.cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating StdoutPipe for Cmd: %v", err)
	}
	stdin, err := w.cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error creating StdinPipe for Cmd: %v", err)
	}

	if err := w.cmd.Start(); err != nil {
		log.Fatalf("Error starting Cmd: %v", err)
	}

	// 从命令行程序读取输出并发送到通道
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			data := scanner.Text()
			log.Print(data)
			var msgServer = types.MsgServer{
				Str:  data,
				Code: 1,
				Zip:  0,
			}
			if w.MZip == 1 {
				msgServer.Zip = 1
			}
			// 引擎数据解析成消息推送给websocket
			w.engineChan <- msgServer

		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from StdoutPipe: %v", err)
		}
		close(w.engineChan)
	}()

	// 从通道读取输入并发送到命令行程序
	go func() {
		for input := range w.wsChan {
			//读取websocket消息并发送给引擎
			if input.Category == "ext" && input.Cmd == "zip" {
				w.MZip = 1
				continue
			}
			_, err := stdin.Write([]byte(input.Cmd + "\n"))
			if err != nil {
				log.Printf("Error writing to StdinPipe: %v", err)
			}

		}
	}()

	// 等待命令结束
	if err := w.cmd.Wait(); err != nil {
		log.Printf("Cmd finished with error: %v", err)
	}
}
