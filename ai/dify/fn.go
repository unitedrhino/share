package dify

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"strings"
)

func SendGetRequest(forConsole bool, dc *DifyClient, api string) (httpCode int, bodyText []byte, err error) {
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		return -1, nil, err
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	bodyText, err = io.ReadAll(resp.Body)
	return resp.StatusCode, bodyText, err
}

func SendPostRequest(forConsole bool, dc *DifyClient, api string, postBody interface{}) (httpCode int, bodyText []byte, err error) {
	var payload *strings.Reader
	if postBody != nil {
		buf, err := json.Marshal(postBody)
		if err != nil {
			return -1, nil, err
		}
		payload = strings.NewReader(string(buf))
	} else {
		payload = nil
	}

	req, err := http.NewRequest("POST", api, payload)
	if err != nil {
		return -1, nil, err
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()

	bodyText, err = io.ReadAll(resp.Body)
	return resp.StatusCode, bodyText, err
}

func SendPostRequestSse[rsp any](ctx context.Context, forConsole bool, dc *DifyClient, api string, postBody interface{}) (httpCode int, bodyText chan rsp, err error) {
	var payload *strings.Reader
	if postBody != nil {
		buf, err := json.Marshal(postBody)
		if err != nil {
			return -1, nil, err
		}
		payload = strings.NewReader(string(buf))
	} else {
		payload = nil
	}

	req, err := http.NewRequest("POST", api, payload)
	if err != nil {
		return -1, nil, err
	}

	if forConsole {
		setConsoleAuthorization(dc, req)
	} else {
		setAPIAuthorization(dc, req)
	}

	resp, err := dc.Client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	if resp.StatusCode != 200 {
		return resp.StatusCode, nil, errors.System
	}
	c := make(chan rsp, 10)
	utils.Go(ctx, func() {
		defer resp.Body.Close()
		defer close(c)
		// 创建一个缓冲读取器来逐行读取响应
		reader := bufio.NewReader(resp.Body)
		for {
			// 读取一行数据
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				logx.WithContext(ctx).Errorf("读取数据时出错: %v", err)
				return
			}

			// 去除行尾的换行符
			line = strings.TrimSpace(line)
			// 忽略空行
			if line == "" {
				continue
			}

			// 解析 SSE 事件
			if strings.HasPrefix(line, "data:") {
				// 提取事件数据
				data := strings.TrimPrefix(line, "data:")
				data = strings.TrimSpace(data)
				var ret rsp
				err = json.Unmarshal([]byte(data), &ret)
				c <- ret
				logx.WithContext(ctx).Debugf("接收到事件数据: %s", data)
			}
		}
	})
	return resp.StatusCode, c, nil
}

func CommonRiskForSendRequest(code int, err error) error {
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return fmt.Errorf("status code: %d", code)
	}

	return nil
}

func CommonRiskForSendRequestWithCode(code int, err error, targetCode int) error {
	if err != nil {
		return err
	}

	if code != targetCode {
		return fmt.Errorf("status code: %d", code)
	}

	return nil
}
