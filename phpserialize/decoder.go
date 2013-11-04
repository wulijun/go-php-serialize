package phpserialize

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PhpDecoder struct {
	source *strings.Reader
}

func Decode(value string) (result interface{}, err error) {
	decoder := &PhpDecoder{
		source: strings.NewReader(value),
	}
	result, err = decoder.DecodeValue()
	return
}

//all integer is int64，float number is float64
func (decoder *PhpDecoder) DecodeValue() (value interface{}, err error) {
	if token, _, err := decoder.source.ReadRune(); err == nil {
		if token == 'N' {
			err = decoder.expect(VALUES_SEPARATOR)
			return nil, err
		}
		decoder.expect(TYPE_VALUE_SEPARATOR)
		switch token {
		case 'b':
			if rawValue, _, _err := decoder.source.ReadRune(); _err == nil {
				value = rawValue == '1'
			} else {
				err = errors.New("Can not read boolean value")
			}
			if err != nil {
				return nil, err
			}
			err = decoder.expect(VALUES_SEPARATOR)
		case 'i':
			if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
				if tmpv, _err := strconv.Atoi(rawValue); _err != nil {
					err = fmt.Errorf("Can not convert %v to Int:%v", rawValue, _err)
				} else {
					value = int64(tmpv)
				}
			} else {
				err = errors.New("Can not read int value")
			}
		case 'd':
			if rawValue, _err := decoder.readUntil(VALUES_SEPARATOR); _err == nil {
				if value, _err = strconv.ParseFloat(rawValue, 64); _err != nil {
					err = fmt.Errorf("Can not convert %v to Float:%v", rawValue, _err)
				}
			} else {
				err = errors.New("Can not read float value")
			}
		case 's':
			value, err = decoder.decodeString()
			if err != nil {
				return nil, err
			}
			err = decoder.expect(VALUES_SEPARATOR)
		case 'a':
			value, err = decoder.decodeArray()
		case 'O':
			value, err = decoder.decodeObject()
		}
	}
	return value, err
}

func (decoder *PhpDecoder) decodeObject() (*PhpObject, error) {
	value := &PhpObject{}
	var err error

	if value.className, err = decoder.decodeString(); err != nil {
		return nil, err
	}
	if err = decoder.expect(TYPE_VALUE_SEPARATOR); err != nil {
		return nil, err
	}
	if value.members, err = decoder.decodeArray(); err != nil {
		return nil, err
	}

	return value, err
}

func (decoder *PhpDecoder) decodeArray() (value map[interface{}]interface{}, err error) {
	value = make(map[interface{}]interface{})
	if rawArrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
		if arrLen, _err := strconv.Atoi(rawArrlen); _err != nil {
			err = fmt.Errorf("Can not convert array length %v to int:%v", rawArrlen, _err)
		} else {
			decoder.expect('{')
			for i := 0; i < arrLen; i++ {
				if k, _err := decoder.DecodeValue(); _err != nil {
					err = fmt.Errorf("Can not read array key %v", _err)
				} else if v, _err := decoder.DecodeValue(); _err != nil {
					err = fmt.Errorf("Can not read array value %v", _err)
				} else {
					switch t := k.(type) {
					default:
						err = fmt.Errorf("Unexpected key type %T", t)
					case string, int64, float64:
						value[k] = v
					}
				}
			}
			decoder.expect('}')
		}
	} else {
		err = errors.New("Can not read array length")
	}
	return value, err
}

func (decoder *PhpDecoder) decodeString() (value string, err error) {
	if rawStrlen, _err := decoder.readUntil(TYPE_VALUE_SEPARATOR); _err == nil {
		if strLen, _err := strconv.Atoi(rawStrlen); _err != nil {
			err = errors.New(fmt.Sprintf("Can not convert string length %v to int:%v", rawStrlen, _err))
		} else {
			if err = decoder.expect('"'); err != nil {
				return
			}
			tmpValue := make([]byte, strLen, strLen)
			if nRead, _err := decoder.source.Read(tmpValue); _err != nil || nRead != strLen {
				err = errors.New(fmt.Sprintf("Can not read string content %v. Read only: %v from %v", _err, nRead, strLen))
			} else {
				value = string(tmpValue)
				err = decoder.expect('"')
			}
		}
	} else {
		err = errors.New("Can not read string length")
	}
	return value, err
}

func (decoder *PhpDecoder) readUntil(stopByte byte) (string, error) {
	result := new(bytes.Buffer)
	var (
		token byte
		err   error
	)
	for {
		if token, err = decoder.source.ReadByte(); err != nil || token == stopByte {
			break
		} else {
			result.WriteByte(token)
		}
	}
	return result.String(), err
}

func (decoder *PhpDecoder) expect(expectRune rune) error {
	token, _, err := decoder.source.ReadRune()
	if err != nil {
		err = fmt.Errorf("Can not read expected: %v", expectRune)
	} else if token != expectRune {
		err = fmt.Errorf("Read %v, but expected: %v", token, expectRune)
	}
	return err
}
