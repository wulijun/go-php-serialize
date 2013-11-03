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
	case int, int64, int32, int16, int8:
		buf.WriteString("i")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := fmt.Sprintf("%v", t)
		buf.WriteString(strValue)
	case float32:
		buf.WriteString("d")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
		buf.WriteString(strValue)
	case float64:
		buf.WriteString("d")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
		buf.WriteString(strValue)
	case string:
		buf.WriteString("s")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		encodeString(buf, t)
	case map[interface{}]interface{}:
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

func encodeArrayCore(buf *bytes.Buffer, arrValue map[interface{}]interface{}) (err error) {
	valLen := strconv.Itoa(len(arrValue))
	buf.WriteString(valLen)
	buf.WriteRune(TYPE_VALUE_SEPARATOR)

	buf.WriteRune('{')
	for k, v := range arrValue {
		if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
			if err = encodeValue(buf, intKey); err != nil {
				break
			}
		} else {
			if err = encodeValue(buf, k); err != nil {
				break
			}
		}
		buf.WriteRune(VALUES_SEPARATOR)
		if err = encodeValue(buf, v); err != nil {
			break
		}
		switch v.(type) {
		case map[interface{}]interface{}, *PhpObject:
			//array and object can't end with VALUES_SEPARATOR
		default:
			buf.WriteRune(VALUES_SEPARATOR)
		}
	}
	buf.WriteRune('}')
	return err
}
