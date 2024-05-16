package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

func UnGzip(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))

	if err != nil {
		var out []byte
		return out, err
	}

	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// UnGzipBase64 从Base64数据解析成byte数组，再解压缩，再转换成字符串
func UnGzipBase64(in string) (string, error) {
	bytesOut, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}

	bytesUnGzip, err := UnGzip(bytesOut)
	if err != nil {
		return "", err
	}

	return string(bytesUnGzip), nil
}
