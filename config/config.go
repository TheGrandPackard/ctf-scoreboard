package config

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config - config
type Config struct {
	StorageType   string `json:"storage_type"`
	StorageConfig string `json:"storage_config"`
	AuthSecret    string `json:"auth_secret"`
}

var configuration *Config

// LoadConfig - laodConfig
func LoadConfig() (c *Config) {
	if configuration != nil {
		return configuration
	}

	c = &Config{}

	loadedConfigFile := false
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err == os.ErrNotExist {
		// Use default values
	} else {
		// Use existing file
		loadedConfigFile = true
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &c)
	}

	// Assign default values
	if c.StorageType == "" {
		c.StorageType = "mysql"
	}
	if c.StorageConfig == "" {
		c.StorageConfig = "ctf-scoreboard:qwerasdf@tcp(127.0.0.1:3306)/ctf-scoreboard"
	}
	if c.AuthSecret == "" {
		c.AuthSecret = generateAuthSecret()
	}

	configuration = c

	// Close config file before writing
	if loadedConfigFile {
		jsonFile.Close()
	}

	configJSON, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Error opening config file for write: %s", err.Error())
	}
	err = ioutil.WriteFile("config.json", configJSON, 0700)
	if err != nil {
		log.Fatalf("Error writing config file: %s", err.Error())
	}

	return
}

func generateAuthSecret() (auth string) {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
