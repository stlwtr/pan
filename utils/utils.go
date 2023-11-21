package utils

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

func InterfaceToString(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case int:
		return strconv.FormatInt(int64(val.(int)), 10)
	case int64:
		return strconv.FormatInt(val.(int64), 10)
	case uint64:
		return strconv.FormatUint(val.(uint64), 10)
	case float32:
		return strconv.FormatFloat(float64(val.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	default:
		bytes, _ := MarshalJSON(val)
		return string(bytes)
	}
}

// 将struct转成url的参数，并去除有json omitempty且值为空的参数，类似param1=1&param2=2
func StructToUrlQuery(m interface{}) (string, error) {
	query := ""

	b, err := MarshalJSON(m)
	if err != nil {
		return query, err
	}

	var f map[string]interface{}
	err = UnmarshalJSON(b, &f)
	if err != nil {
		return query, err
	}

	v := url.Values{}
	for key, value := range f {
		v.Set(key, InterfaceToString(value))
	}
	query = v.Encode()
	return query, nil
}

func UnescapeUnicode(raw []byte) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}

type any = interface{}

func MarshalJSON(v any) ([]byte, error) {
	var resultBytes bytes.Buffer
	enc := json.NewEncoder(&resultBytes)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	return resultBytes.Bytes(), err
}

func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
