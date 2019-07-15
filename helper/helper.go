package helper

import (
	"time"
)

func DateEqual(date1, date2 time.Time) bool {
	return date1.Format("20060102") == date2.Format("20060102")
}

func ToMapString(data interface{}) (res map[string]string) {
	res = map[string]string{}

	for k, v := range data.(map[string]interface{}) {
		res[k] = v.(string)
	}
	return
}
