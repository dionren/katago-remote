package types

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Config struct {
	Port           string `json:"port"`
	KataGoPath     string `json:"katago_path"`
	ConfigFilePath string `json:"config_file"`
	ModelFilePath  string `json:"model_file"`
}

// GlobalConfig 全局配置变量
var GlobalConfig *ini.File

// LoadConfig 从INI配置文件加载配置的函数
func LoadConfig(configPath string) error {
	config, err := ini.Load(configPath)
	if err != nil {
		return err
	}
	GlobalConfig = config
	return nil
}

// GetConfigValue 根据section和key获取配置的函数
func GetConfigValue(section string, key string) (string, error) {
	if GlobalConfig == nil {
		return "", fmt.Errorf("global config.ini not initialized")
	}
	getKey, err := GlobalConfig.Section(section).GetKey(key)
	if err != nil {
		return "", err
	}
	return getKey.String(), nil
}
