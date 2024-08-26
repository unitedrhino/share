package utils

import (
	"bytes"
	"encoding/csv"
	"gitee.com/i-Things/share/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"path"
	"unicode/utf8"
)

func ReadExcel(file io.Reader, fileName string, tablename ...string) ([][]string, error) {
	ext := path.Ext(fileName)
	switch ext {
	case ".csv":
		fb, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		//删除 BOM 字符
		bom := []byte{0xEF, 0xBB, 0xBF} // BOM 字符
		if bytes.HasPrefix(fb, bom) {
			fb = fb[len(bom):] // 删除前三个字节
		}
		fr := bytes.NewReader(fb)
		// 兼容 UTF-8 和 GBK/GB2312
		var reader *csv.Reader
		if utf8.Valid(fb) {
			reader = csv.NewReader(fr)
		} else {
			decoder := simplifiedchinese.GBK.NewDecoder()
			reader = csv.NewReader(transform.NewReader(fr, decoder))
		}
		rows, err := reader.ReadAll()
		return rows, err
	case ".xlsx":
		//读取文件路径
		f, err := excelize.OpenReader(file)
		if err != nil {
			return nil, err
		}
		firstSheet := ""
		if len(tablename) > 0 {
			firstSheet = tablename[0]
		} else {
			firstSheet = f.GetSheetName(0)
		}
		rows, err := f.GetRows(firstSheet)
		return rows, err
	}
	return nil, errors.NotRealize
}
