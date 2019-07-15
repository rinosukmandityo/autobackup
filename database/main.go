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

	helper.BackupDBToS3(configs["database"].(map[string]interface{}), configs["s3"].(map[string]interface{}))

}
