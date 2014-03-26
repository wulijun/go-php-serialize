package phpserialize

import (
	"bytes"
	"fmt"
	"strconv"
)

func Encode(value interface{}) (result string, err error) {
	buf := new(bytes.Buffer)
	err = encodeValue(buf, value)
	if err == nil {
		result = buf.String()
	}
	return
}

func encodeValue(buf *bytes.Buffer, value interface{}) (err error) {
	switch t := value.(type) {
	default:
		err = fmt.Errorf("Unexpected type %T", t)
	case bool:
		buf.WriteString("b")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		if t {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
		buf.WriteRune(VALUES_SEPARATOR)
	case nil:
		buf.WriteString("N")
		buf.WriteRune(VALUES_SEPARATOR)
	case int, int64, int32, int16, int8:
		buf.WriteString("i")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := fmt.Sprintf("%v", t)
		buf.WriteString(strValue)
		buf.WriteRune(VALUES_SEPARATOR)
	case float32:
		buf.WriteString("d")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
		buf.WriteString(strValue)
		buf.WriteRune(VALUES_SEPARATOR)
	case float64:
		buf.WriteString("d")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
		buf.WriteString(strValue)
		buf.WriteRune(VALUES_SEPARATOR)
	case string:
		buf.WriteString("s")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		encodeString(buf, t)
		buf.WriteRune(VALUES_SEPARATOR)
	case map[string]interface{}, []interface{}, map[interface{}]interface{}:
		buf.WriteString("a")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		err = encodeArrayCore(buf, t)
	case *PhpObject:
		buf.WriteString("O")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		encodeString(buf, t.GetClassName())
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		err = encodeArrayCore(buf, t.GetMembers())
	}
	return
}

func encodeString(buf *bytes.Buffer, strValue string) {
	valLen := strconv.Itoa(len(strValue))
	buf.WriteString(valLen)
	buf.WriteRune(TYPE_VALUE_SEPARATOR)
	buf.WriteRune('"')
	buf.WriteString(strValue)
	buf.WriteRune('"')
}

func rangeArrayCore(buf *bytes.Buffer, arrValue interface{}) (err error) {

	if arrVal, ok := arrValue.([]interface{}); ok {
		for k, v := range arrVal {
			if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
				if err = encodeValue(buf, intKey); err != nil {
					break
				}
			} else {
				if err = encodeValue(buf, k); err != nil {
					break
				}
			}
			if err = encodeValue(buf, v); err != nil {
				break
			}
		}
	} else if arrVal, ok := arrValue.(map[interface{}]interface{}); ok {
		for k, v := range arrVal {
			if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
				if err = encodeValue(buf, intKey); err != nil {
					break
				}
			} else {
				if err = encodeValue(buf, k); err != nil {
					break
				}
			}
			if err = encodeValue(buf, v); err != nil {
				break
			}
		}
	} else if arrVal, ok := arrValue.(map[string]interface{}); ok {
		for k, v := range arrVal {
			if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
				if err = encodeValue(buf, intKey); err != nil {
					break
				}
			} else {
				if err = encodeValue(buf, k); err != nil {
					break
				}
			}
			if err = encodeValue(buf, v); err != nil {
				break
			}
		}
	}

	return err
}

func encodeArrayCore(buf *bytes.Buffer, arrValue interface{}) (err error) {
	valLen := ""
	if arr, ok := arrValue.([]interface{}); ok {
		valLen = strconv.Itoa(len(arr))
	} else if arr, ok := arrValue.(map[interface{}]interface{}); ok {
		valLen = strconv.Itoa(len(arr))
	} else if arr, ok := arrValue.(map[string]interface{}); ok {
		valLen = strconv.Itoa(len(arr))
	}

	if "" == valLen {
		return fmt.Errorf("Unexpected type of array %T", arrValue)
	}

	buf.WriteString(valLen)
	buf.WriteRune(TYPE_VALUE_SEPARATOR)

	buf.WriteRune('{')
	err = rangeArrayCore(buf, arrValue)
	buf.WriteRune('}')
	return err
}
