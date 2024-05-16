package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"os"
	"server/types"
)

// GzipBase64 将字符串压缩成gzip格式，然后编码为Base64字符串
func GzipBase64(in string) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err := gz.Write([]byte(in))
	if err != nil {
		return "", err
	}

	if err := gz.Close(); err != nil {
		return "", err
	}

	// 获取压缩后的数据
	compressed := b.Bytes()

	// 将压缩后的数据编码为Base64
	return base64.StdEncoding.EncodeToString(compressed), nil
}

func LoadConfig(configPath string) (*types.Config, error) {
	config := &types.Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
