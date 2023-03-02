package configs

import (
	"encoding/json"
	"gitlab.com/gbh007/gojlog"
	"os"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	ConnStr            string `json:"connect_string"`
	Port               string `json:"port"`
	SecretKeyReCaptcha string `json:"secretKeyReCaptcha"`
	Auth               struct {
		Mongo struct {
			URI      string `json:"uri"`
			User     string `json:"user,omitempty"`
			Password string `json:"password,omitempty"`
			DBName   string `json:"dbname"`
		} `json:"mongo"`
	} `json:"auth"`
}

var (
	_config *Config
)

// GetConfig получение объекта конфига
func GetConfig() *Config {
	return _config
}

func Load() error {
	confFile, err := os.Open("configs/config.json")
	if err != nil {
		gojlog.Critical(err)
		return err
	}
	defer confFile.Close()
	dc := json.NewDecoder(confFile)
	if err := dc.Decode(&_config); err != nil {
		gojlog.Critical("Чтение конфиг файла: ", err)
		return err
	}
	if _config.ConnStr == "" {
		gojlog.Critical("Не могу прочитать конфиг файл: ", _config.ConnStr)
		return err
	}
	return nil
}
