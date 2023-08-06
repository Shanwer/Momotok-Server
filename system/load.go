package system

import (
	"Momotok-Server/common"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

//LoadConfigInformation load config information for application

func LoadConfigInformation() error {
	const configPath string = "common\\config.yaml"
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("unable to read config file: %v", err)
	}

	var config struct {
		Server struct {
			Mode                 string `yaml:"mode"`
			Host                 string `yaml:"host"`
			Port                 string `yaml:"port"`
			EnableHttps          bool   `yaml:"enable_https"`
			TokenExpireSecond    int    `yaml:"token_expire_second"`
			PreviewImagePath     string `yaml:"preview_image_path"`
			VideoUrlPath         string `yaml:"video_url_path"`
			MaxPerPage           int    `yaml:"max_per_page"`
			DatabaseAddress      string `yaml:"database_address"`
			DriverName           string `yaml:"driver_name"`
			SecretKey            string `yaml:"secret_key"`
		} `yaml:"server"`
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return fmt.Errorf("config parse failed: %v", err)
	}

	common.ServerInfo.Mode = config.Server.Mode
	common.ServerInfo.Host = config.Server.Host
	common.ServerInfo.Port = config.Server.Port
	common.ServerInfo.EnableHttps = config.Server.EnableHttps
	common.ServerInfo.TokenExpireSecond = config.Server.TokenExpireSecond
	common.ServerInfo.PreviewImagePath = config.Server.PreviewImagePath
	common.ServerInfo.VideoUrlPath = config.Server.VideoUrlPath
	common.ServerInfo.MaxPerPage = config.Server.MaxPerPage
	common.ServerInfo.DatabaseAddress = config.Server.DatabaseAddress
	common.ServerInfo.DriverName = config.Server.DriverName
	common.ServerInfo.SecretKey = config.Server.SecretKey

	return nil
}
