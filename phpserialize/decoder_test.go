package phpserialize

import (
	"math"
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
		if k == "flt32" {
			if math.Abs(v.(float64)-float64(data2["flt32"].(float32))) > 0.001 {
				t.Errorf("key arr[%v] value expeted %v, get %v", k, v, data2[k])
			}
		} else if v != data2[k] {
			t.Errorf("key arr[%v] value expeted %v, get %v", k, v, data2[k])
		}
	}
	/*
		var (
			stringKey string
			//normalKey interface{}
		)
		stringKey = "flt32"
		//normalKey = stringKey
		t.Errorf("%T key %v %v", decodeData2, decodeData2[stringKey], decodeData2["int64"])
		result = "a:33:{s:4:\"cost\";d:457.00000000000011;s:3:\"uid\";s:9:\"119983495\";i:2332002;d:15;i:2735232;d:15;i:66324;d:14.999999999999998;i:55107;d:14;i:84875;d:14;i:65145;d:15;i:437688;d:15;i:40022;d:14;i:2538364;d:15;i:34701;d:15.000000000000002;i:33490;d:15;i:2813734;d:14;i:1878927;d:15;i:1837284;d:15;i:2735238;d:15;i:437385;d:15;i:44723;d:14;i:2538367;d:15;i:84014;d:15;i:427239;d:15.000000000000002;i:76055;d:15;i:67652;d:14.000000000000002;i:33482;d:15;i:45859;d:14;i:44248;d:14;i:44724;d:15;i:45861;d:15;i:37995;d:15;i:39985;d:15;i:40994;d:15;i:42360;d:15;}"
		decodeRes, _ = Decode(result)
		decodeData, _ = decodeRes.(map[interface{}]interface{})
		for k, v := range decodeData {
			t.Errorf("%T %v => %T %v", k, k, v, v)
		}

		result = "a:34:{i:66324;d:15;i:67652;d:15;i:427239;d:15;i:437688;d:15;i:1837284;d:15;s:7:\"balance\";d:1617.05;s:5:\"quota\";d:0;i:55107;d:15;i:65145;d:15;i:1878927;d:15;i:2332002;d:15;i:33490;d:15;i:34701;d:15;i:37995;d:15;i:39985;d:15;i:40022;d:15;i:40994;d:15;i:42360;d:15;i:44248;d:15;i:44723;d:15;i:44724;d:15;i:45859;d:15;i:45861;d:15;i:2439794;d:10;i:33482;d:15;i:76055;d:15;i:84014;d:15;i:84875;d:15;i:437385;d:15;i:2538364;d:15;i:2538367;d:15;i:2735232;d:15;i:2735238;d:15;i:2813734;d:15;}"
		decodeRes, _ = Decode(result)
		decodeData, _ = decodeRes.(map[interface{}]interface{})
		for k, v := range decodeData {
			t.Errorf("%T %v => %T %v", k, k, v, v)
		}

		t.Errorf("66324 => %v %v balance => %v", decodeData[66324], decodeData[int64(66324)], decodeData["balance"])
		//output  66324 => <nil> 15 balance => 1617.05

		t.Errorf("decode %v %v %v %v", err, decodeRes, obj.GetClassName(), privateMemberValue)

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
		t.Errorf("decode %v %v", err, decodeRes)
	*/
}
