package dify

import (
	"context"
	"encoding/json"
	"fmt"
)

type ChatMessagesPayload struct {
	Inputs         any                       `json:"inputs"`
	Query          string                    `json:"query"`
	ResponseMode   string                    `json:"response_mode"`
	ConversationID string                    `json:"conversation_id,omitempty"`
	User           string                    `json:"user"`
	Files          []ChatMessagesPayloadFile `json:"files,omitempty"`
}

type ChatMessagesPayloadFile struct {
	Type           string `json:"type"`
	TransferMethod string `json:"transfer_method"`
	URL            string `json:"url,omitempty"`
	UploadFileID   string `json:"upload_file_id,omitempty"`
}

type ChatMessagesResponse struct {
	Event          string `json:"event"`
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       any    `json:"metadata"`
	CreatedAt      int    `json:"created_at"`
}

/*
event: message LLM 返回文本块事件，即：完整的文本以分块的方式输出。

	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	conversation_id (string) 会话 ID
	answer (string) LLM 返回文本块内容
	created_at (int) 创建时间戳，如：1705395332

event: agent_message Agent模式下返回文本块事件，即：在Agent模式下，文章的文本以分块的方式输出（仅Agent模式下使用）

	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	conversation_id (string) 会话 ID
	answer (string) LLM 返回文本块内容
	created_at (int) 创建时间戳，如：1705395332

event: agent_thought Agent模式下有关Agent思考步骤的相关内容，涉及到工具调用（仅Agent模式下使用）

	id (string) agent_thought ID，每一轮Agent迭代都会有一个唯一的id
	task_id (string) 任务ID，用于请求跟踪下方的停止响应接口
	message_id (string) 消息唯一ID
	position (int) agent_thought在消息中的位置，如第一轮迭代position为1
	thought (string) agent的思考内容
	observation (string) 工具调用的返回结果
	tool (string) 使用的工具列表，以 ; 分割多个工具
	tool_input (string) 工具的输入，JSON格式的字符串(object)。如：{"dalle3": {"prompt": "a cute cat"}}
	created_at (int) 创建时间戳，如：1705395332
	message_files (array[string]) 当前 agent_thought 关联的文件ID
	file_id (string) 文件ID
	conversation_id (string) 会话ID

event: message_file 文件事件，表示有新文件需要展示

	id (string) 文件唯一ID
	type (string) 文件类型，目前仅为image
	belongs_to (string) 文件归属，user或assistant，该接口返回仅为 assistant
	url (string) 文件访问地址
	conversation_id (string) 会话ID
	event: message_end 消息结束事件，收到此事件则代表流式返回结束。
	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	conversation_id (string) 会话 ID
	metadata (object) 元数据
	usage (Usage) 模型用量信息
	retriever_resources (array[RetrieverResource]) 引用和归属分段列表

event: tts_message TTS 音频流事件，即：语音合成输出。内容是Mp3格式的音频块，使用 base64 编码后的字符串，播放的时候直接解码即可。(开启自动播放才有此消息)

	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	audio (string) 语音合成之后的音频块使用 Base64 编码之后的文本内容，播放的时候直接 base64 解码送入播放器即可
	created_at (int) 创建时间戳，如：1705395332

event: tts_message_end TTS 音频流结束事件，收到这个事件表示音频流返回结束。

	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	audio (string) 结束事件是没有音频的，所以这里是空字符串
	created_at (int) 创建时间戳，如：1705395332

event: message_replace 消息内容替换事件。 开启内容审查和审查输出内容时，若命中了审查条件，则会通过此事件替换消息内容为预设回复。

	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	conversation_id (string) 会话 ID
	answer (string) 替换内容（直接替换 LLM 所有回复文本）
	created_at (int) 创建时间戳，如：1705395332

event: error 流式输出过程中出现的异常会以 stream event 形式输出，收到异常事件后即结束。

	task_id (string) 任务 ID，用于请求跟踪和下方的停止响应接口
	message_id (string) 消息唯一 ID
	status (int) HTTP 状态码
	code (string) 错误码
	message (string) 错误消息

event: ping 每 10s 一次的 ping 事件，保持连接存活。
*/
type ChatMessagesSseResponse struct {
	Event          string   `json:"event"`
	MessageID      string   `json:"message_id"`
	ConversationID string   `json:"conversation_id"`
	Mode           string   `json:"mode"`
	Answer         string   `json:"answer"`
	Metadata       any      `json:"metadata"`
	CreatedAt      int64    `json:"created_at"`
	Audio          string   `json:"audio"`
	ID             string   `json:"id"`
	TaskID         string   `json:"task_id"`
	Position       int      `json:"position"`
	Thought        string   `json:"thought"`
	Observation    string   `json:"observation"`
	Tool           string   `json:"tool"`
	ToolInput      string   `json:"tool_input"`
	MessageFiles   []string `json:"message_files"`
}

func PrepareChatPayload(payload map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (dc *DifyClient) ChatMessages(inputs map[string]interface{}, query string, conversation_id string, files []any) (result ChatMessagesResponse, err error) {
	var payload ChatMessagesPayload

	payload.Inputs = inputs
	payload.Query = query

	payload.ResponseMode = RESPONSE_MODE_BLOCKING
	payload.User = dc.User

	if conversation_id != "" {
		payload.ConversationID = conversation_id
	}

	if len(files) > 0 {
		// TODO TBD
		return result, fmt.Errorf("files are not supported")
	}

	api := dc.GetAPI(API_CHAT_MESSAGES)

	code, body, err := SendPostRequestToAPI(dc, api, payload)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal the response: %v", err)
	}
	return result, nil
}

func (dc *DifyClient) ChatMessagesStreaming(ctx context.Context, inputs map[string]interface{}, query string, conversation_id string, files []any) (result chan ChatMessagesSseResponse, err error) {
	var payload ChatMessagesPayload

	payload.Inputs = inputs
	payload.Query = query

	payload.ResponseMode = RESPONSE_MODE_STREAMING
	payload.User = dc.User

	if conversation_id != "" {
		payload.ConversationID = conversation_id
	}

	if len(files) > 0 {
		// TODO TBD
		return nil, fmt.Errorf("files are not supported")
	}

	api := dc.GetAPI(API_CHAT_MESSAGES)

	code, body, err := SendPostSseRequestToAPI[ChatMessagesSseResponse](ctx, dc, api, payload)

	err = CommonRiskForSendRequest(code, err)
	if err != nil {
		return result, err
	}

	// if !strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
	// 	return "", fmt.Errorf("response is not a streaming response")
	// }

	return body, nil
}
