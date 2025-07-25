package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

// 编码转换函数
func convertEncoding(input []byte, encoder encoding.Encoding, toUTF8 bool) ([]byte, error) {
	var transformer transform.Transformer
	if toUTF8 {
		// 转换为 UTF-8
		transformer = encoder.NewDecoder()
	} else {
		// 从 UTF-8 转换为目标编码
		transformer = encoder.NewEncoder()
	}

	reader := transform.NewReader(bytes.NewReader(input), transformer)
	output, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("编码转换失败: %v", err)
	}
	return output, nil
}

// GBK 转 UTF-8
func GBKToUTF8(input []byte) ([]byte, error) {
	return convertEncoding(input, simplifiedchinese.GBK, true)
}

// UTF-8 转 GBK
func UTF8ToGBK(input []byte) ([]byte, error) {
	return convertEncoding(input, simplifiedchinese.GBK, false)
}
