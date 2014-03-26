package phpserialize

import (
	"testing"
)

func TestEncodeArrayValue2(t *testing.T) {
	data := make(map[interface{}]interface{})
	data2 := make(map[interface{}]interface{})
	data2["test"] = true
	data2["0"] = 5
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
	/*
		result, err := Encode(data)
		t.Errorf("data %v => %v\n", err, result)
		result, err = Encode(data3)
		t.Errorf("data %v => %v\n", err, result)
		result, err = Encode(nil)
		t.Errorf("data %v => %v\n", err, result)
		result, err = Encode("s\"tr'}e")
		t.Errorf("data %v => %v\n", err, result)*/
}
