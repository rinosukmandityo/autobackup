package main

import (
	"flag"
	"path/filepath"

	"github.com/rinosukmandityo/autobackup/helper"
)

func main() {

	var configLoc string
	flag.StringVar(&configLoc, "config", filepath.Join(helper.WD, "configs", "configs.json"), "config file location")
	flag.Parse()

	configs := helper.ReadJsonFile(configLoc)

	helper.BackupDBToS3(configs[helper.CONF_FOR_DB].(map[string]interface{}), configs[helper.CONF_FOR_S3].(map[string]interface{}))

}
