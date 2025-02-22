package config

import (
	"encoding/json"
	"errors"
	"log"
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
	// do something to turn configs into input scanners
	if c.configMap == nil {
		return errors.New("config file not loaded")
	}

	var inputScanners []json.RawMessage
	err := json.Unmarshal(c.configMap["inputScanners"], &inputScanners)
	if err != nil {
		return err
	}
	// inputScanners, ok := c.configMap["inputScanners"]
	// if !ok {
	// 	return errors.New("No input scanners in config")
	// }

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

		log.Println(is)
		c.inputScanners = append(c.inputScanners, is)
	}

	// // asserts and converts interface{} -> []interface{}
	// for _, inputScanner := range inputScanners.([]interface{}) {
	// 	inputScannerMap := inputScanner.(map[string]interface{})

	// 	name, ok := inputScannerMap["name"]
	// 	if !ok {
	// 		name = "scanner"
	// 	}

	// 	scannerType, ok := inputScannerMap["type"]
	// 	if !ok {
	// 		return errors.New("empty scanner type")
	// 	}

	// 	scannerColor, ok := inputScannerMap["color"]
	// 	if !ok {
	// 		// TODO: assign random color
	// 		scannerColor = "black"
	// 	}
	// 	colorAttribute, ok := colorMap[scannerColor.(string)]
	// 	if !ok {
	// 		colorAttribute = color.FgBlack
	// 	}

	// 	var is scanner.InputScanner
	// 	if scannerType == "CMD" {
	// 		cmd, ok := inputScannerMap["command"]
	// 		if !ok {
	// 			return errors.New("no command provided for scanner of type CMD")
	// 		}
	// 		cis := &scanner.CmdInputScanner{}
	// 		cis.SetCmd(cmd.(string))
	// 		is = cis
	// 	} else if scannerType == "K8S" {
	// 		pod := inputScannerMap["pod"].(map[string]interface{})

	// 		podName := pod["name"].(string)
	// 		podSelector := pod["podSelector"].(map[string]interface{})
	// 		if podName == "" && podSelector == nil {
	// 			return errors.New("must provide one of pod name or podSelector for scanner of type K8S")
	// 		}

	// 		namespace, ok := pod["namespace"].(string)
	// 		if !ok {
	// 			log.Println("W! No namespace provided - using default namespace")
	// 			namespace = "default"
	// 		}
	// 		container := pod["container"].(string)

	// 		kis := &scanner.K8sInputScanner{}
	// 		kis.SetPodName(podName)
	// 		// kis.SetPodSelector(podSelector)
	// 		kis.SetNamespace(namespace)
	// 		kis.SetContainer(container)
	// 		is = kis
	// 	} else {
	// 		return fmt.Errorf("unknown scanner type: %s", scannerType)
	// 	}

	// 	// Set common fields
	// 	is.SetName(name.(string))
	// 	is.SetOutputColor(color.New(colorAttribute))
	// 	c.inputScanners = append(c.inputScanners, is)

	// }
	return nil
}
