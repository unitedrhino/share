package ai

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	wsURL     = "wss://dashscope.aliyuncs.com/api-ws/v1/inference/" // WebSocket服务器地址
	audioFile = "asr_example.wav"                                   // 替换为您的音频文件路径
)

var dialer = websocket.DefaultDialer

type Ali struct {
	*websocket.Conn
	apiKey string
}

func NewAli(apiKey string) *Ali {
	conn, err := connectWebSocket(apiKey)
	if err != nil {
		log.Fatal("连接WebSocket失败：", err)
	}
	// 启动一个goroutine来接收结果
	taskStarted := make(chan bool)
	taskDone := make(chan bool)
	startResultReceiver(conn, taskStarted, taskDone)
	// 发送run-task指令
	//_, err = sendRunTaskCmd(conn)
	//if err != nil {
	//	log.Fatal("发送run-task指令失败：", err)
	//}

	// 等待task-started事件
	//waitForTaskStarted(taskStarted)
	return &Ali{
		Conn:   conn,
		apiKey: apiKey,
	}
}

func (t *Ali) SendAudio(data []byte) error {
	err := t.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return err
	}
	return nil
}

func (t *Ali) SendText(taskID string, texts ...string) error {
	for _, text := range texts {
		runTaskCmd, err := generateContinueTaskCmd(text, taskID)
		if err != nil {
			return err
		}

		err = t.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
		if err != nil {
			return err
		}
	}
	return nil
}

//func test() {
//	// 若没有将API Key配置到环境变量，可将下行替换为：apiKey := "your_api_key"。不建议在生产环境中直接将API Key硬编码到代码中，以减少API Key泄露风险。
//	apiKey := os.Getenv("DASHSCOPE_API_KEY")
//
//	// 连接WebSocket服务
//	conn, err := connectWebSocket(apiKey)
//	if err != nil {
//		log.Fatal("连接WebSocket失败：", err)
//	}
//	defer closeConnection(conn)
//
//	// 启动一个goroutine来接收结果
//	taskStarted := make(chan bool)
//	taskDone := make(chan bool)
//	startResultReceiver(conn, taskStarted, taskDone)
//
//	// 发送run-task指令
//	taskID, err := sendRunTaskCmd(conn)
//	if err != nil {
//		log.Fatal("发送run-task指令失败：", err)
//	}
//
//	// 等待task-started事件
//	waitForTaskStarted(taskStarted)
//
//	// 发送待识别音频文件流
//	if err := sendAudioData(conn); err != nil {
//		log.Fatal("发送音频失败：", err)
//	}
//
//	// 发送finish-task指令
//	if err := sendFinishTaskCmd(conn, taskID); err != nil {
//		log.Fatal("发送finish-task指令失败：", err)
//	}
//
//	// 等待任务完成或失败
//	<-taskDone
//}

// 定义结构体来表示JSON数据
type Header struct {
	Action       string                 `json:"action"`
	TaskID       string                 `json:"task_id"`
	Streaming    string                 `json:"streaming"`
	Event        string                 `json:"event"`
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Attributes   map[string]interface{} `json:"attributes"`
}

type Output struct {
	Sentence struct {
		BeginTime int64  `json:"begin_time"`
		EndTime   *int64 `json:"end_time"`
		Text      string `json:"text"`
		Words     []struct {
			BeginTime   int64  `json:"begin_time"`
			EndTime     *int64 `json:"end_time"`
			Text        string `json:"text"`
			Punctuation string `json:"punctuation"`
		} `json:"words"`
	} `json:"sentence"`
	Usage interface{} `json:"usage"`
}

type Payload struct {
	TaskGroup  string `json:"task_group"`
	Task       string `json:"task"`
	Function   string `json:"function"`
	Model      string `json:"model"`
	Parameters Params `json:"parameters"`
	// 不使用热词功能时，不要传递resources参数
	// Resources  []Resource `json:"resources"`
	Input  Input  `json:"input"`
	Output Output `json:"output,omitempty"`
}

type Params struct {
	Format                   string   `json:"format"`
	SampleRate               int      `json:"sample_rate"`
	VocabularyID             string   `json:"vocabulary_id"`
	DisfluencyRemovalEnabled bool     `json:"disfluency_removal_enabled"`
	LanguageHints            []string `json:"language_hints"`
	TextType                 string   `json:"text_type"`
	Voice                    string   `json:"voice"`
	Volume                   int      `json:"volume"`
	Rate                     int      `json:"rate"`
	Pitch                    int      `json:"pitch"`
}

// 不使用热词功能时，不要传递resources参数
type Resource struct {
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
}

type Input struct {
	Text string `json:"text"`
}

type Event struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

// 连接WebSocket服务
func connectWebSocket(apiKey string) (*websocket.Conn, error) {
	header := make(http.Header)
	header.Add("X-DashScope-DataInspection", "enable")
	header.Add("Authorization", fmt.Sprintf("bearer %s", apiKey))
	conn, _, err := dialer.Dial(wsURL, header)
	return conn, err
}

// 启动一个goroutine异步接收WebSocket消息
func startResultReceiver(conn *websocket.Conn, taskStarted chan<- bool, taskDone chan<- bool) {
	go func() {
		for {
			msgType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("解析服务器消息失败：", err)
				return
			}
			if msgType == websocket.BinaryMessage {
				// 处理二进制音频流
				if err := writeBinaryDataToFile(message, outputFile); err != nil {
					fmt.Println("写入二进制数据失败：", err)
					return
				}
			} else {
				// 处理文本消息
				var event Event
				err = json.Unmarshal(message, &event)
				if err != nil {
					fmt.Println("解析事件失败：", err)
					continue
				}
				if handleEvent(conn, event, taskStarted, taskDone) {
					return
				}
			}
		}
	}()
}

// 发送run-task指令
func sendRunTaskCmd(conn *websocket.Conn) (string, error) {
	runTaskCmd, taskID, err := generateAsrTaskCmd()
	if err != nil {
		return "", err
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return taskID, err
}

func (t *Ali) StartAsr() (string, error) {
	runTaskCmd, taskID, err := generateAsrTaskCmd()
	if err != nil {
		return "", err
	}
	err = t.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return taskID, err
}

func (t *Ali) StartTts() (string, error) {
	runTaskCmd, taskID, err := generateTtsTaskCmd()
	if err != nil {
		return "", err
	}
	err = t.WriteMessage(websocket.TextMessage, []byte(runTaskCmd))
	return taskID, err
}

// 生成run-task指令
func generateAsrTaskCmd() (string, string, error) {
	taskID := uuid.New().String()
	runTaskCmd := Event{
		Header: Header{
			Action:    "run-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			TaskGroup: "audio",
			Task:      "asr",
			Function:  "recognition",
			Model:     "paraformer-realtime-v2",
			Parameters: Params{
				Format:     "opus",
				SampleRate: 16000,
			},
			Input: Input{},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), taskID, err
}

// 生成run-task指令
func generateTtsTaskCmd() (string, string, error) {
	taskID := uuid.New().String()
	runTaskCmd := Event{
		Header: Header{
			Action:    "run-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			TaskGroup: "audio",
			Task:      "tts",
			Function:  "SpeechSynthesizer",
			Model:     "cosyvoice-v1",
			Parameters: Params{
				TextType:   "PlainText",
				Voice:      "longxiaochun",
				Format:     "mp3",
				SampleRate: 22050,
				Volume:     50,
				Rate:       1,
				Pitch:      1,
			},
			Input: Input{},
		},
	}
	runTaskCmdJSON, err := json.Marshal(runTaskCmd)
	return string(runTaskCmdJSON), taskID, err
}

// 等待task-started事件
func waitForTaskStarted(taskStarted chan bool) {
	select {
	case <-taskStarted:
		fmt.Println("任务开启成功")
	case <-time.After(10 * time.Second):
		log.Fatal("等待task-started超时，任务开启失败")
	}
}

// 生成finish-task指令
func generateFinishTaskCmd(taskID string) (string, error) {
	finishTaskCmd := Event{
		Header: Header{
			Action:    "finish-task",
			TaskID:    taskID,
			Streaming: "duplex",
		},
		Payload: Payload{
			Input: Input{},
		},
	}
	finishTaskCmdJSON, err := json.Marshal(finishTaskCmd)
	return string(finishTaskCmdJSON), err
}

// 处理事件
func handleEvent(conn *websocket.Conn, event Event, taskStarted chan<- bool, taskDone chan<- bool) bool {
	fmt.Println(event)
	switch event.Header.Event {
	case "task-started":
		fmt.Println("收到task-started事件")
		taskStarted <- true
	case "result-generated":
		if event.Payload.Output.Sentence.Text != "" {
			fmt.Println("识别结果：", event.Payload.Output.Sentence.Text)
		}
	case "task-finished":
		fmt.Println("任务完成")
		taskDone <- true
		return true
	case "task-failed":
		handleTaskFailed(event, conn)
		taskDone <- true
		return true
	default:
		log.Printf("预料之外的事件：%v", event)
	}
	return false
}

// 处理任务失败事件
func handleTaskFailed(event Event, conn *websocket.Conn) {
	if event.Header.ErrorMessage != "" {
		log.Fatalf("任务失败：%s", event.Header.ErrorMessage)
	} else {
		log.Fatal("未知原因导致任务失败")
	}
}

// 关闭连接
func closeConnection(conn *websocket.Conn) {
	if conn != nil {
		conn.Close()
	}
}
