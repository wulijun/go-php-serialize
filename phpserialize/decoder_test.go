package phpserialize

import (
	"testing"
)

func TestDecodeArrayValue2(t *testing.T) {
	data := make(map[interface{}]interface{})
	data2 := make(map[interface{}]interface{})
	data2["test"] = true
	data2[int64(0)] = int64(5)
	data2["flt32"] = float32(2.3)
	data2["int64"] = int64(45)
	data3 := NewPhpObject()
	data3.SetClassName("A")
	data3.SetPrivateMemberValue("a", 1)
	data3.SetProtectedMemberValue("b", 3.14)
	data3.SetPublicMemberValue("c", data2)
	data["arr"] = data2
	data["3"] = "s\"tr'}e"
	data["g"] = nil
	data["object"] = data3

	var (
		result    string
		decodeRes interface{}
		err       error
	)
	if result, err = Encode(data); err != nil {
		t.Errorf("encode data fail %v, %v", err, data)
		return
	}
	if decodeRes, err = Decode(result); err != nil {
		t.Errorf("decode data fail %v, %v", err, result)
		return
	}
	decodeData, ok := decodeRes.(map[interface{}]interface{})
	if !ok {
		t.Errorf("decode data type error, expected: map[interface{}]interface{}, get: %T", decodeRes)
		return
	}
	obj, _ := decodeData["object"].(*PhpObject)
	if v, _ := obj.GetPrivateMemberValue("a"); v != int64(1) {
		t.Errorf("object private value expected 1, get %v, %T", v, v)
	}
	if v := obj.GetClassName(); v != "A" {
		t.Errorf("object class name expected A, get %v", v)
	}
	if decodeData["g"] != nil {
		t.Errorf("key g value expeted nil, get %v", decodeData["g"])
	}

	decodeData2, ok := decodeData["arr"].(map[interface{}]interface{})
	if !ok {
		t.Errorf("key arr value type expeted map, get %T", decodeData["arr"])
	}
	for k, v := range decodeData2 {
		if v != data2[k] {
			t.Errorf("key arr[%v] value expeted %v, get %v", k, v, data2[k])
		}
	}

	/*	t.Errorf("decode %v %v %v %v", err, decodeRes, obj.GetClassName(), privateMemberValue)

		result, err = Encode(data3)
		t.Errorf("data %v => %v\n", err, result)
		decodeRes, err = Decode(result)
		t.Errorf("decode %v %v", err, decodeRes)

		result, err = Encode(nil)
		t.Errorf("data %v => %v\n", err, result)
		decodeRes, err = Decode(result)
		t.Errorf("decode %v %v", err, decodeRes)

		result, err = Encode("s\"tr'}e")
		t.Errorf("data %v => %v\n", err, result)
		decodeRes, err = Decode(result)
		t.Errorf("decode %v %v", err, decodeRes)*/
}
