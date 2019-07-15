package helper

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "os"
)

func ReadJsonFile(fPath string) (res map[string]interface{}) {
    res = map[string]interface{}{}
    jsonFile, err := os.Open(fPath)
    if err != nil {
        log.Println(err)
        return
    }
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal([]byte(byteValue), &res)
    return
}
