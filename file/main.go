package main

import (
	"flag"

	"github.com/rinosukmandityo/autobackup/helper"
)

func main() {

	var configLoc string
	flag.StringVar(&configLoc, "config", "configs/configs.json", "config file location")
	flag.Parse()

	configs := helper.ReadJsonFile(configLoc)
	fileconfig := configs["file"].(map[string]interface{})
	helper.PutObjectsToS3(fileconfig, configs["s3"].(map[string]interface{}))
	if fileconfig["initialrun"].(bool) {
		fileconfig["initialrun"] = false
		configs["file"] = fileconfig
		helper.WriteJsonFile(configs, configLoc)
	}
}
