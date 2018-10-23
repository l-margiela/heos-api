package heosapi

import (
	"fmt"
)

func paramsToStr(params map[string]string) string {
	var str string

	first := true
	for k, v := range params {
		if first {
			first = false
			str = fmt.Sprintf("%s=%s", k, v)
			continue
		}
		str += fmt.Sprintf("&%s=%s", k, v)
	}

	return str
}
