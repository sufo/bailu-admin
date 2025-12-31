/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Addable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string
}

type Comparable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string | bool
}

type Uint64_string interface {
	uint64 | string
}

// UpperFirst upper first char
func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	rs := []rune(s)
	f := rs[0]

	if 'a' <= f && f <= 'z' {
		return string(unicode.ToUpper(f)) + string(rs[1:])
	}
	return s
}

func JsonStr2Map(jsonStr string) (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &tempMap)
	if err != nil {
		return nil, nil
	}
	return tempMap, nil
}

func JsonStr2FlatMap[T any](jsonStr string) (map[string]T, error) {
	d, err := JsonStr2Map(jsonStr)
	if err != nil {
		return nil, err
	}
	var flatMap map[string]T
	FlatMap[T]("", d, flatMap)
	return flatMap, nil
}

// flatten map
func FlatMap[T any](prefix string, src map[string]interface{}, dest map[string]T) {
	if len(prefix) > 0 {
		prefix += "."
	}
	for k, v := range src {
		switch child := v.(type) {
		case map[string]interface{}:
			FlatMap(prefix+k, child, dest)
		case []interface{}:
			for i := 0; i < len(child); i++ {
				dest[prefix+k+"."+strconv.Itoa(i)] = child[i].(T)
			}
		default:
			dest[prefix+k] = v.(T)
		}
	}
}

// 遍历取struct单个属性组成新数组
// func StructsField2Arr[T any](source []interface{}, key string) ([]T, error) {
func StructsField2Arr[T any](source interface{}, key string) ([]T, error) {
	v := reflect.ValueOf(source)
	if v.Kind() != reflect.Slice {
		// 如果输入不是 slice，就返回错误
		// todo: handle error
		return nil, errors.New("source param must be a slice")
	}
	// 获取 slice 的长度
	n := v.Len()
	var arr []T
	for i := 0; i < n; i++ {
		if v.Index(i).IsValid() {
			var val T
			var err error
			var item = v.Index(i).Interface()
			if realVal, ok := item.(map[string]interface{}); ok {
				var success bool
				val, success = realVal[key].(T)
				if !success {
					err = errors.New("field type convert fail")
				}
			} else {
				val, err = getStructFieldVal[T](item, key)
			}
			if err != nil {
				return nil, err
			}
			arr = append(arr, val)
		} else {
			return arr, errors.New("source param resolve failed")
		}
	}
	return arr, nil
}

func getStructFieldVal[T any](obj interface{}, field string) (t T, err error) {
	/* 获取反射对象信息 */
	point2Struct := reflect.ValueOf(obj)
	structVal := point2Struct.Elem()
	structFieldValue := structVal.FieldByName(field)
	if !structFieldValue.IsValid() {
		return t, fmt.Errorf("No such field: %s in obj", field)
	}
	//structFieldType := structFieldValue.Type()
	return structFieldValue.Interface().(T), nil
}

func getStructFieldVal2[T any](v interface{}, property string) (t T, err error) {
	m, _ := json.Marshal(v)
	var x map[string]interface{}
	err = json.Unmarshal(m, &x)
	if err != nil {
		return t, err
	}
	return x[property].(T), nil
}

// ContainsInSlice 判断字符串是否在 slice 中
func ContainsInSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func Includes[T Comparable](items []T, item T) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 符号后面第一个字符转大写
// @s 原始字符串
// @symbol 符号
func ToUpperForFirstCharAtSymbolBehind(s, symbol string) string {
	if s == "" {
		return s
	}
	s = strings.ReplaceAll(s, symbol, "/")
	b := []byte(s)
	size := len(b)
	for i := 0; i < size; i++ {
		if i == 0 {
			b[0] = byte(unicode.ToUpper(rune(b[0])))
		}
		if b[i] == 47 { // 47 ascii "/"
			if i+1 < size {
				b[i+1] = byte(unicode.ToUpper(rune(b[i+1])))
			}
		}
	}
	name := strings.ReplaceAll(string(b), "/", "")
	return name
}

// 去重
func RemoveDuplicates(selects []string) {
	// 对Select进行去重
	selectMap := make(map[string]int, len(selects))
	for _, e := range selects {
		if _, ok := selectMap[e]; !ok {
			selectMap[e] = 1
		}
	}
}

// 生成0-max之间随机数
func RandNumber(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return r.Intn(max)
}

// 生成长度为length的随机字符串
func RandString(length int64) string {
	sources := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sourceLength := len(sources)
	var i int64 = 0
	for ; i < length; i++ {
		result = append(result, sources[r.Intn(sourceLength)])
	}

	return string(result)
}
