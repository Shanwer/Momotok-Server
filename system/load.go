package system

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ServerInfo *configModel = &configModel{} // server config information
// LoadConfigInformation load config information for application
func LoadConfigInformation() error {
	const configPath string = "system\\config.yaml"
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("unable to read config file: %v", err)
	}

	err = yaml.Unmarshal(configData, &ServerInfo)
	if err != nil {
		return fmt.Errorf("config parse failed: %v", err)
	}

	return nil
}
