//对json文件的操作

package jsonx

import (
	"encoding/json"
)

//Js json数据结构体
type Js struct {
	mata interface{}
}

//NewJSONString 通过string创建json结构
func NewJSONString(data string) (*Js, error) {
	j := new(Js)
	var f interface{}
	err := json.Unmarshal([]byte(data), &f)
	if err != nil {
		return j, err
	}
	j.mata = f
	return j, err
}

//NewJSONByte 通过[]byte创建json结构
func NewJSONByte(data []byte) (*Js, error) {
	j := new(Js)
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return j, err
	}
	j.mata = f
	return j, err
}

//NewGet 通过json结构创建json结构
func NewGet(j *Js, key string) *Js {
	value := new(Js)
	value.mata = j.get(key)
	return value
}

//NewGetIndex 通过json结构创建json结构
func NewGetIndex(j *Js, i int) *Js {
	value := new(Js)
	value.mata = j.getList(i)
	return value
}

//GetNum 获取[]数组数量
func (j *Js) GetNum() int {
	if m, ok := (j.mata).([]interface{}); ok {
		return len(m)
	}
	return 0
}

//IsValid 判断json结构是否有数据
func (j *Js) IsValid() bool {
	if nil == j.mata {
		return false
	}
	return true
}

//GetUint16 通过key获取uint16
func (j *Js) GetUint16(key string) uint16 {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(float64); ok {
			return uint16(m)
		}
	}
	return 0
}

//GetString 通过key获取string
func (j *Js) GetString(key string) string {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(string); ok {
			return m
		}
	}
	return ""
}

//GetInt32 通过key获取int32
func (j *Js) GetInt32(key string) int32 {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(float64); ok {
			return int32(m)
		}
	}
	return 0
}

//GetInt 通过key获取int
func (j *Js) GetInt(key string) int {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(float64); ok {
			return int(m)
		}
	}
	return -1
}

//GetBool 通过key获取bool
func (j *Js) GetBool(key string) bool {
	data := j.get(key)
	if data != nil {
		if m, ok := data.(bool); ok {
			return m
		}
	}
	return false
}

func (j *Js) get(key string) interface{} {
	if m, ok := (j.mata).(map[string]interface{}); ok {
		if data, oki := m[key]; oki {
			return data
		}
	}
	return nil
}
func (j *Js) getList(i int) interface{} {

	if m, ok := (j.mata).([]interface{}); ok {
		if i >= 0 && i < len(m) {
			return m[i]
		}
	}
	return nil
}
