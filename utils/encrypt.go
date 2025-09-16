package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"gitee.com/unitedrhino/share/errors"
)

type HmacType = string

var (
	HmacTypeSha256 HmacType = "hmacsha256"
	HmacTypeSha1   HmacType = "hmacsha1"
	HmacTypeMd5    HmacType = "hmacmd5"
)

func Hmac(sign HmacType, data string, secret []byte) string {
	sign = strings.ToLower(sign)
	switch sign {
	case HmacTypeSha1:
		return HmacSha1(data, secret)
	case HmacTypeSha256:
		return HmacSha256(data, secret)
	default:
		return HmacMd5(data, secret)
	}
}

func HmacSha256(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha1(data string, secret []byte) string {
	h := hmac.New(sha1.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// HexToBase64 将16进制字符串转换为Base64编码
func HexToBase64(hexStr string) (string, error) {
	// 首先将16进制字符串解码为字节数组
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	// 然后将字节数组编码为Base64字符串
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HmacMd5(data string, secret []byte) string {
	h := hmac.New(md5.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padLen := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}
func AesCbcBase64(src, productSecret string) (string, error) {
	if src == "" || productSecret == "" {
		return "", errors.Default.AddMsg("加密参数错误")
	}
	// 截取 productSecret 前 16 位作为密钥
	key := []byte(productSecret)[:16]
	// 以长度 16 的字符 "0" 作为偏移量
	iv := bytes.Repeat([]byte("0"), 16)

	data := []byte(src)

	// 使用 AES-CBC 加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 对补全后的数据进行加密
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	cryptData := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cryptData, data)

	// 进行 base64 编码
	return base64.StdEncoding.EncodeToString(cryptData), nil
}

func Md5Map(params map[string]any) string {
	// 排序
	keys := make([]string, len(params))
	i := 0
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	stringBuf := bytes.Buffer{}
	for _, k := range keys {
		stringBuf.WriteString(fmt.Sprintf("%s%v", k, params[k]))
	}
	md := md5.Sum(stringBuf.Bytes())
	return hex.EncodeToString(md[:])
}
