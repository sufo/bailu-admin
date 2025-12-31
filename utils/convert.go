/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 命名转换、 类型转换
 */

package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// Ordinary way using bit shifting
func Uint64ToBytes2(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}

	return r
}

// Use `encoding/binary` package
func Uint64ToBytes(val uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)
	return b
}

// 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}

// Remove struct name from validation error messages
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

var errNegativeNotAllowed = errors.New("unable to cast negative value")

// 基本数据类型转换
func ToT[T any](value any) (T, error) {
	var key string
	if value == nil {
		return any(key).(T), nil
	}
	var t T
	var v any
	var err error
	switch any(t).(type) {
	case bool:
		v, err = ToBool(value)
	case float64:
		v, err = ToFloat[float64](value)
	case float32:
		v, err = ToFloat[float64](value)
	case int:
		v, err = ToInt[int](value)
	case uint:
		v, err = ToUint[uint](value)
	case int8:
		v, err = ToInt[int8](value)
	case uint8:
		v, err = ToUint[uint8](value)
	case int16:
		v, err = ToInt[int16](value)
	case uint16:
		v, err = ToUint[uint16](value)
	case int32:
		v, err = ToInt[int32](value)
	case uint32:
		v, err = ToUint[uint32](value)
	case int64:
		v, err = ToInt[int64](value)
	case uint64:
		v, err = ToUint[uint64](value)
	case string:
		key = Strval(value)
	//case []byte:
	//	key = string(value.([]byte))
	default:
		err = fmt.Errorf("unable to cast %#v of type %T to bool", value, value)
	}
	if err != nil {
		return t, err
	}
	t = any(v).(T)
	return t, nil
}

// ToBoolE casts any type to a bool type.
func ToBool(a any) (bool, error) {
	a = indirect(a)
	switch b := a.(type) {
	case bool:
		return b, nil
	case nil:
		return false, nil
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8, float64, float32, uintptr, complex64, complex128:
		return !reflect.ValueOf(a).IsZero(), nil
	case string:
		return strconv.ParseBool(a.(string))
	case time.Duration:
		return b != 0, nil
	case json.Number:
		v, err := b.Float64()
		return v != 0, err
	default:
		return false, fmt.Errorf("unable to cast %#v of type %T to bool", a, a)
	}
}

// toInt returns the int value of v if v or v's underlying type is an int.
// Note that this will return false for int64 etc. types.
func _toInt(v any) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
}
func ToInt[T int8 | int16 | int | int32 | int64](i any) (T, error) {
	i = indirect(i)

	intv, ok := _toInt(i)
	if ok {
		return T(intv), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		return T(s), nil
	case int32:
		return T(s), nil
	case int16:
		return T(s), nil
	case int8:
		return T(s), nil
	case uint:
		return T(s), nil
	case uint64:
		return T(s), nil
	case uint32:
		return T(s), nil
	case uint16:
		return T(s), nil
	case uint8:
		return T(s), nil
	case float64:
		return T(s), nil
	case float32:
		return T(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if err == nil {
			return T(v), nil
		}
		return 0, fmt.Errorf("unable to cast %#v of type %T to int64", i, i)
	case json.Number:
		v, err := s.Int64()
		return T(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to int", i, i)
	}
}

func ToUint[T uint8 | uint16 | uint | uint32 | uint64](i any) (T, error) {
	i = indirect(i)
	intv, ok := _toInt(i)
	if ok {
		if intv < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(intv), nil
	}

	switch s := i.(type) {
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case int64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(s), nil
	case int32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(s), nil
	case int16:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(s), nil
	case int8:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(s), nil
	case uint:
		return T(s), nil
	case uint64:
		return any(s).(T), nil
	case uint32:
		return T(s), nil
	case uint16:
		return T(s), nil
	case uint8:
		return T(s), nil
	case float32:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(s), nil
	case float64:
		if s < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(s), nil
	case string:
		v, err := strconv.ParseInt(trimZeroDecimal(s), 0, 0)
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(v), err
	case json.Number:
		v, err := s.Int64()
		if v < 0 {
			return 0, errNegativeNotAllowed
		}
		return T(v), err
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to uint64", i, i)
	}
}

func ToFloat[T float32 | float64](i any) (T, error) {
	i = indirect(i)
	intv, ok := _toInt(i)
	if ok {
		return T(intv), nil
	}
	switch s := i.(type) {
	case float64:
		return T(s), nil
	case float32:
		return T(s), nil
	case int64:
		return T(s), nil
	case int32:
		return T(s), nil
	case int16:
		return T(s), nil
	case int8:
		return T(s), nil
	case uint:
		return T(s), nil
	case uint64:
		return T(s), nil
	case uint32:
		return T(s), nil
	case uint16:
		return T(s), nil
	case uint8:
		return T(s), nil
	case string:
		var t T
		switch any(t).(type) {
		case float64:
			if r, err := strconv.ParseFloat(s, 64); err != nil {
				return 0, err
			} else {
				return any(r).(T), nil
			}
		case float32:
			if r, err := strconv.ParseFloat(s, 64); err != nil {
				return 0, err
			} else {
				return any(r).(T), nil
			}
		default:
			return t, nil
		}
	case json.Number:
		if r, err := s.Float64(); err != nil {
			return 0, err
		} else {
			return T(r), err
		}
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("unable to cast %#v of type %T to float64", i, i)
	}
}

// Copied from html/template/content.go.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a any) any {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Pointer {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// trimZeroDecimal trims the zero decimal.
// E.g. 12.00 to 12 while 12.01 still to be 12.01.
func trimZeroDecimal(s string) string {
	var foundZero bool
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				return s[:i-1]
			}
		case '0':
			foundZero = true
		default:
			return s
		}
	}
	return s
}

// converts a slice or array to map[interface{}]struct{} with error
func SliceToMap(i interface{}) (map[interface{}]struct{}, error) {
	// judge the validation of the input
	if i == nil {
		return nil, fmt.Errorf("unable to converts %#v of type %T to map[interface{}]struct{}", i, i)
	}
	kind := reflect.TypeOf(i).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("the input %#v of type %T isn't a slice or array", i, i)
	}
	// execute the convert
	v := reflect.ValueOf(i)
	m := make(map[interface{}]struct{}, v.Len())
	for j := 0; j < v.Len(); j++ {
		m[v.Index(j).Interface()] = struct{}{}
	}
	return m, nil
}

// struct to Map[string]interface{}
func StructToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // Non-structural return error
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// Traversing structure fields
	// Specify the tagName value as the key in the map; the field value as the value in the map
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

// Converting structures to single-level maps
func NestedStructToMap(in interface{}, tag string) (map[string]interface{}, error) {

	// The current function only receives struct types
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr { // Structure Pointer
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	out := make(map[string]interface{})
	queue := make([]interface{}, 0, 1)
	queue = append(queue, in)

	for len(queue) > 0 {
		v := reflect.ValueOf(queue[0])
		if v.Kind() == reflect.Ptr { // Structure Pointer
			v = v.Elem()
		}
		queue = queue[1:]
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			vi := v.Field(i)
			if vi.Kind() == reflect.Ptr { // Embedded Pointer
				vi = vi.Elem()
				if vi.Kind() == reflect.Struct { // Structures
					queue = append(queue, vi.Interface())
				} else {
					ti := t.Field(i)
					if tagValue := ti.Tag.Get(tag); tagValue != "" {
						// Save to map
						out[tagValue] = vi.Interface()
					}
				}
				break
			}
			if vi.Kind() == reflect.Struct { // Embedded Structs
				queue = append(queue, vi.Interface())
				break
			}
			// General Fields
			ti := t.Field(i)
			if tagValue := ti.Tag.Get(tag); tagValue != "" {
				// Save to map
				out[tagValue] = vi.Interface()
			}
		}
	}
	return out, nil
}

// struct -> map
// 直接使用json.Marshal方法来强制转化struct。
func Struct2Map(in any) map[string]any {
	var target map[string]interface{}
	if marshalContent, err := json.Marshal(in); err != nil {
		fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&target); err != nil {
			fmt.Println(err)
		} else {
			for k, v := range target {
				target[k] = v
			}
		}
	}
	return target
}

// map to url params
func Map2UrlParams(m map[string]interface{}) (result string) {
	list := make([]string, 0)
	for k, v := range m {
		t1 := fmt.Sprintf("%s=%s", k, fmt.Sprint(v))
		list = append(list, t1)
	}
	result = strings.Join(list, "&")
	return
}

func Map2String(m map[string]interface{}) (result string) {
	data, err := json.Marshal(m)
	if err == nil {
		return string(data)
	} else {
		fmt.Errorf("map to string err. err=%v", err)
		return ""
	}
}

func StrArr2Arr[T any](s []string) ([]T, error) {
	var res = make([]T, 0)
	for _, v := range s {
		t, err := ToT[T](v)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

// 四舍五入
// f 待处理的数字
// n 保留小数位
func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
