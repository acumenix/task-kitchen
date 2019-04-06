package kitchen_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type testConfig struct {
	TableName   string `json:"table_name"`
	TableRegion string `json:"table_region"`
}

var testCfg testConfig

func init() {
	confPath := "./test.json"
	if newPath := os.Getenv("TEST_CONFIG_PATH"); newPath != "" {
		confPath = newPath
	}

	raw, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatalf("Fail to read test config file: %s, %s", confPath, err)
	}

	if err := json.Unmarshal(raw, &testCfg); err != nil {
		log.Fatalf("Fail to unmarshal test config file: %s, %s", confPath, err)
	}
}
