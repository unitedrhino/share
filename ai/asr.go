package ai

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"os"
)

const (
	outputFile = "output.mp3" // 输出文件路径
)

// 发送待合成文本
func sendContinueTaskCmd(conn *websocket.Conn, taskID string) error {
	texts := []string{"床前明月光", "疑是地上霜", "举头望明月", "低头思故乡"}

	for _, text := range texts {
		runTaskCmd, err := generateContinueTaskCmd(text, taskID)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
		if err != nil {
			return err
		}
	}

	return nil
}

// 生成continue-task指令
func generateContinueTaskCmd(text string, taskID string) (string, error) {
	runTaskCmd := Event{
		Header: Header{
			Action:    "continue-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{
				Text: text,
			},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), err
}

// 写入二进制数据到文件
func writeBinaryDataToFile(data []byte, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// 清空输出文件
func clearOutputFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
