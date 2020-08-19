package gjson

import (
	"fmt"
	"strings"

	"github.com/valyala/fastjson/fastfloat"

	"github.com/valyala/fastjson"
)

//GJSON struct
type GJSON struct {
	data  *fastjson.Value
	arena fastjson.Arena
}

//////////////////////////////////////// Get* ////////////////////////////////

//GetBool  only fastjson.TypeBool Type  get bool value return correct
func (t *GJSON) GetBool(k string) bool {
	return t.data.GetBool(k)
}

func (t *GJSON) IsExist(k string) bool {
	return t.data.Exists(k)
}

//GetString get string value
func (t *GJSON) GetString(k string) string {
	v := t.data.Get(k)
	if v == nil {
		return ""
	}
	switch v.Type() {
	case fastjson.TypeNull, fastjson.TypeObject, fastjson.TypeArray:
		return ""
	case fastjson.TypeFalse:
		return "false"
	case fastjson.TypeNumber:
		return v.String()
	case fastjson.TypeTrue:
		return "true"
	case fastjson.TypeString:
		return string([]byte(v.String())[1 : len(v.String())-1])
	}
	return ""
}
func (t *GJSON) GetBytes(k string) []byte {
	v := t.data.Get(k)
	if v == nil {
		return nil
	}
	switch v.Type() {
	case fastjson.TypeNull, fastjson.TypeObject, fastjson.TypeArray:
		return nil
	case fastjson.TypeNumber, fastjson.TypeFalse, fastjson.TypeTrue:
		return v.MarshalTo(nil)
	case fastjson.TypeString:
		res, err := v.StringBytes()
		if err != nil {
			return nil
		}
		return res
	}
	return nil
}
func (t *GJSON) GetFloat64(k string) float64 {
	v := t.data.Get(k)
	if v == nil {
		return 0
	}
	switch v.Type() {
	case fastjson.TypeFalse, fastjson.TypeNull, fastjson.TypeObject, fastjson.TypeArray:
		return 0
	case fastjson.TypeNumber:
		r, err := v.Float64()
		if err != nil {
			fmt.Println(err.Error())
			return 0
		}
		return r
	case fastjson.TypeTrue:
		return 1
	case fastjson.TypeString:

		val := strings.Replace(v.String(), ",", "", -1)
		val = string([]byte(val)[1 : len(val)-1])
		res := fastfloat.ParseBestEffort(val)
		return res
	}
	return 0
}

func (t *GJSON) GetInt(k string) int {
	return int(t.GetInt64(k))
}
func (t *GJSON) GetInt64(k string) int64 {
	v := t.data.Get(k)
	if v == nil {
		return 0
	}
	switch v.Type() {
	case fastjson.TypeFalse, fastjson.TypeNull, fastjson.TypeObject, fastjson.TypeArray:
		return 0
	case fastjson.TypeNumber:
		//r := t.data.GetInt64(k)
		val := v.String()
		if strings.Contains(val, ".") {
			r, err := v.Float64()
			if err != nil {
				fmt.Println(err.Error())
				return 0
			}
			return int64(r)

		} else {
			r, err := v.Int64()
			if err != nil {
				fmt.Println(err.Error())
				return 0
			}
			return r
		}
		return 0
	case fastjson.TypeTrue:
		return 1
	case fastjson.TypeString:
		val := v.String()
		val = strings.Replace(val, ",", "", -1)
		val = string([]byte(val)[1 : len(val)-1])
		res := fastfloat.ParseInt64BestEffort(val)
		return res
	}
	return 0
}
func (t *GJSON) GetObject(k string) (result *GJSON) {
	result = new(GJSON)
	result.data = t.data.Get(k)
	return
}
func (t *GJSON) GetArray(k string) (result *GJSON) {
	result = new(GJSON)
	result.data = t.data.Get(k)
	return
}

func (t *GJSON) Item(index int) (result *GJSON) {
	v, err := t.data.Array()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result = new(GJSON)
	result.data = v[index]
	return
}

//Load  load string
func (t *GJSON) Load(s string) (err error) {
	v, err := fastjson.Parse(s)
	if err != nil {
		return
	}
	t.data = v
	return
}

func (t *GJSON) ItemCount() (num int) {
	t.initGJSON(2)
	if t.data.Type() == fastjson.TypeArray {
		slic, err := t.data.Array()
		if err != nil {
			return 0
		}
		num = len(slic)
	}
	return
}

//AddItem  添加一个子对象
func (t *GJSON) AddItem() (result *GJSON) {
	t.initGJSON(2)
	if t.data.Type() == fastjson.TypeArray {
		slic, err := t.data.Array()
		if err != nil {
			return nil
		}
		var temp GJSON
		temp.initGJSON(1)
		result = &temp
		t.data.SetArrayItem(len(slic), temp.data)
	}
	return
}

func (t *GJSON) ToString() string {
	return t.data.String()
}

func (t *GJSON) initGJSON(flag int) {
	if t.data == nil {
		switch flag {
		case 1:
			t.data = t.arena.NewObject()
		case 2:
			t.data = t.arena.NewArray()
		default:
			t.data = t.arena.NewObject()
		}
	}
}

//////////////////////////////////////////////////Set*/////////////////////////
func (t *GJSON) SetBool(k string, v bool) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		if v == true {
			t.data.Set(k, t.arena.NewTrue())
		} else {
			t.data.Set(k, t.arena.NewFalse())
		}
	}
}

func (t *GJSON) SetFloat64(k string, v float64) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewNumberFloat64(v))
	}
}

func (t *GJSON) SetInt(k string, v int) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewNumberInt(v))
	}
}

func (t *GJSON) SetString(k, v string) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewString(v))
	}
}

func (t *GJSON) SetBytes(k string, v []byte) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewStringBytes(v))
	}
}

func (t *GJSON) SetObject(k string, v GJSON) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, v.data)
	}
}

func (t *GJSON) SetArray(k string, v GJSON) {
	t.initGJSON(1)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, v.data)
	}
}
