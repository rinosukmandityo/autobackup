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
	helper.BackupFileToS3(configs["file"].(map[string]interface{}), configs["s3"].(map[string]interface{}))
}