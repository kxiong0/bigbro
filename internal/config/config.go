package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/kxiong0/bigbro/internal/scanner"
)

type Config struct {
	configMap     map[string]json.RawMessage
	inputScanners []scanner.InputScanner
}

func (c *Config) LoadConfigFile(filename string) error {
	err := c.loadConfigMap(filename)
	if err != nil {
		return err
	}

	err = c.parseConfigMap()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) GetInputScanners() []scanner.InputScanner {
	return c.inputScanners
}

func (c *Config) loadConfigMap(filename string) error {
	var config map[string]json.RawMessage // Generic map to hold JSON data
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &config)
	c.configMap = config
	return err
}

func (c *Config) parseConfigMap() error {
	if c.configMap == nil {
		return errors.New("config file not loaded")
	}

	var inputScanners []json.RawMessage
	err := json.Unmarshal(c.configMap["inputScanners"], &inputScanners)
	if err != nil {
		return err
	}

	for _, inputScanner := range inputScanners {
		var scannerMap map[string]interface{}
		err = json.Unmarshal(inputScanner, &scannerMap)
		if err != nil {
			return err
		}

		var is scanner.InputScanner
		scannerType := scannerMap["type"].(string)
		switch scannerType {
		case "CMD":
			is = &scanner.CmdInputScanner{}
			err = json.Unmarshal(inputScanner, &is)
		case "K8S":
			is = &scanner.K8sInputScanner{}
			err = json.Unmarshal(inputScanner, &is)
		default:
			return errors.New("invalid scanner type")
		}

		if err != nil {
			return err
		}
		c.inputScanners = append(c.inputScanners, is)
	}
	return nil
}
