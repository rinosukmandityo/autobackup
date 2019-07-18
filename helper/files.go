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

func WriteJsonFile(data map[string]interface{}, fPath string) (e error) {
    jsonFile, e := os.Open(fPath)
    if e != nil {
        log.Println(e)
        return
    }
    defer jsonFile.Close()

    f, e := json.MarshalIndent(data, "", "\t")
    e = ioutil.WriteFile(fPath, f, 0777)
    return
}

func IsFileNotExist(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return true
    }

    return false
}

func IsFileExist(path string) bool {
    if _, err := os.Stat(path); err == nil {
        return true
    }

    return false
}
