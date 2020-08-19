package gjson

import (
	"fmt"
	"strings"
	"time"

	"github.com/valyala/fastjson/fastfloat"

	"github.com/valyala/fastjson"
)

type GjsonType int

const (
	_ = iota
	TypeObject
	TypeArray
)

//GJSON struct
type GJSON struct {
	data  *fastjson.Value
	arena fastjson.Arena
}

func NewGJSON(t GjsonType) *GJSON {
	result := new(GJSON)
	result.initGJSON(t)
	return result
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
func (t *GJSON) GetTime(k string) (result time.Time) {
	v := t.data.Get(k)
	if v == nil || v.Type() != fastjson.TypeString {
		return
	}
	str := v.String()
	tm, err := time.Parse("2006-01-02", str)
	if err != nil {
		tm, err = time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			tm, err = time.Parse("2006-01-02T15:04:05.999999999+08:00", str)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return tm
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

func (t *GJSON) Interface() interface{} {
	//map[string]interface{}
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		obj, err := t.data.Object()
		if err != nil {
			return make(map[string]interface{})
		} else {
			result := make(map[string]interface{}, obj.Len())
			obj.Visit(func(key []byte, v *fastjson.Value) {
				switch v.Type() {
				case fastjson.TypeArray, fastjson.TypeObject, fastjson.TypeNull:
					result[string(key)] = ""
				case fastjson.TypeString:
					result[string(key)] = string([]byte(v.String())[1 : len(v.String())-1])
				case fastjson.TypeTrue:
					result[string(key)] = true
				case fastjson.TypeFalse:
					result[string(key)] = false
				case fastjson.TypeNumber:
					if strings.Contains(v.String(), ".") {
						f, _ := v.Float64()
						result[string(key)] = f
					} else {
						i, _ := v.Int64()
						result[string(key)] = i
					}
				}
			})
			return result
		}
	}
	return nil
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
	t.initGJSON(TypeArray)
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
	t.initGJSON(TypeArray)
	if t.data.Type() == fastjson.TypeArray {
		slic, err := t.data.Array()
		if err != nil {
			return nil
		}
		var temp GJSON
		temp.initGJSON(TypeObject)
		result = &temp
		t.data.SetArrayItem(len(slic), temp.data)
	}
	return
}

func (t *GJSON) ToString() string {
	return t.data.String()
}

func (t *GJSON) initGJSON(flag GjsonType) {
	if t.data == nil {
		switch flag {
		case TypeObject:
			t.data = t.arena.NewObject()
		case TypeArray:
			t.data = t.arena.NewArray()
		default:
			t.data = t.arena.NewObject()
		}
	}
}

//////////////////////////////////////////////////Set*/////////////////////////
func (t *GJSON) SetBool(k string, v bool) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		if v == true {
			t.data.Set(k, t.arena.NewTrue())
		} else {
			t.data.Set(k, t.arena.NewFalse())
		}
	}
}

func (t *GJSON) SetFloat64(k string, v float64) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewNumberFloat64(v))
	}
}

func (t *GJSON) SetInt(k string, v int) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewNumberInt(v))
	}
}

func (t *GJSON) SetString(k, v string) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewString(v))
	}
}

func (t *GJSON) SetBytes(k string, v []byte) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, t.arena.NewStringBytes(v))
	}
}

func (t *GJSON) SetObject(k string, v GJSON) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, v.data)
	}
}

func (t *GJSON) SetArray(k string, v GJSON) {
	t.initGJSON(TypeObject)
	if t.data.Type() == fastjson.TypeObject {
		t.data.Set(k, v.data)
	}
}
