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
	fileconfig := configs[helper.CONF_FOR_FILE].(map[string]interface{})
	helper.PutObjectsToS3(fileconfig, configs[helper.CONF_FOR_S3].(map[string]interface{}))
	if fileconfig[helper.CONF_INIT_RUN].(bool) {
		fileconfig[helper.CONF_INIT_RUN] = false
		configs[helper.CONF_FOR_FILE] = fileconfig
		helper.WriteJsonFile(configs, configLoc)
	}
}
