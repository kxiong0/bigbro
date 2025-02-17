package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/kxiong0/bigbro/internal/scanner"
)

var colorMap = map[string]color.Attribute{
	"black":   color.FgBlack,
	"white":   color.FgWhite,
	"red":     color.FgRed,
	"green":   color.FgGreen,
	"yellow":  color.FgYellow,
	"blue":    color.FgBlue,
	"magenta": color.FgMagenta,
	"cyan":    color.FgCyan,
}

type Config struct {
	configMap     map[string]interface{}
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
	var config map[string]interface{} // Generic map to hold JSON data
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
		return errors.New("parsing a configMap not loaded")
	}
	inputScanners, ok := c.configMap["inputScanners"]
	if ok {
		// asserts and converts interface{} -> []interface{}
		for _, inputScanner := range inputScanners.([]interface{}) {
			inputScannerMap := inputScanner.(map[string]interface{})

			name, ok := inputScannerMap["name"]
			if !ok {
				name = "scanner"
			}

			scannerType, ok := inputScannerMap["type"]
			if !ok {
				log.Fatal("Empty scanner type")
			}

			scannerColor, ok := inputScannerMap["color"]
			if !ok {
				scannerColor = "black"
			}
			colorAttribute, ok := colorMap[scannerColor.(string)]
			if !ok {
				colorAttribute = color.FgBlack
			}

			if scannerType == "CMD" {
				cmd, ok := inputScannerMap["command"]
				if !ok {
					log.Fatal("No command provided for scanner of type CMD")
				}
				cis := &scanner.CmdInputScanner{}
				cis.SetName(name.(string))
				cis.SetOutputColor(color.New(colorAttribute))
				cis.SetCmd(cmd.(string))

				c.inputScanners = append(c.inputScanners, cis)
			} else {
				log.Fatal("Unknown scanner type:", scannerType)
			}
		}
	}
	return nil
}
