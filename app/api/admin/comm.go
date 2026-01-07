/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 公共方法
 */

package admin

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func ParseParamIDs(c *gin.Context, key string) []uint64 {
	val := c.Param(key)
	if val == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be required", key)))
	}
	idStrArr := strings.Split(val, ",")
	var ids = make([]uint64, 0)
	for _, idStr := range idStrArr {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			//return 0
			panic(respErr.BadRequestErrorWithError(err))
		}
		ids = append(ids, id)
	}
	return ids
}
func ParseParamArray[T string | int | int64 | uint64 | float32 | float64](c *gin.Context, key string) []T {
	val := c.Param(key)
	if val == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be required", key)))
	}
	arr := strings.Split(val, ",")
	var params = make([]T, 0)
	var t T
	switch s := any(t).(type) {
	case string:
		for _, item := range arr {
			params = append(params, any(item).(T))
		}
	default:
		for _, item := range arr {
			if res, err := utils.ToT[T](item); err != nil {
				panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be %T type", key, s)))
			} else {
				params = append(params, res)
			}
		}
	}
	return params
}

func ParseParamId(c *gin.Context, key string) uint64 {
	val := c.Param(key)
	if val == "" {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be required", key)))
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(respErr.BadRequestErrorWithError(err))
	}
	return id
}

// wrap c.Query, add default value
func Query[T string | int | float32](c *gin.Context, key string, defaultVal T) T {
	v := c.Query(key)
	if v == "" || v == "null" {
		return defaultVal
	}
	return any(v).(T)
}

func GetLoginUser(c *gin.Context) *entity.OnlineUserDto {
	v, exist := c.Get(consts.REQUEST_USER)
	if !exist {
		panic(respErr.ForbiddenError)
	}
	return v.(*entity.OnlineUserDto)
}

func ReqSchema(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme
}

// 获取文件url
func FileUrl(r *http.Request, path string) string {
	var uploadConf = config.Conf.Upload
	if uploadConf.Model == "local" {
		return ReqSchema(r) + "://" + r.Host + config.Conf.Local.Path + "/" + path
	} else {
		return path
	}
}

/**转换树
 * @sourceField 源字段名，
 * @targetField 目标字段名
 */
func TransformTree[T any](tree []T, sourceField []string, targetField []string) ([]map[string]any, error) {
	if len(tree) == 0 {
		//return nil, nil
		return make([]map[string]any, 0), nil
	}
	if len(sourceField) != len(targetField) {
		return nil, fmt.Errorf("sourceField`s and targetField`s length must be equal")
	}
	var target = make([]map[string]any, 0)
	for _, item := range tree {
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct && v.Kind() != reflect.Map { // Non-structural return error
			return nil, fmt.Errorf("tree only accepts struct,map or struct pointer; got %T", v)
		}

		var kv = make(map[string]any, 0)
		for index, value := range sourceField {
			var vv reflect.Value
			vi := v.FieldByName(value)
			if vi.Kind() == reflect.Invalid { //表示[value]字段找不到
				continue
			}
			vv = vi
			if vi.Kind() == reflect.Ptr { // Embedded Pointer
				vi = vi.Elem()
			}
			if vi.Kind() == reflect.Array || vi.Kind() == reflect.Slice {
				var nextTree = vi.Interface().([]T)
				if nextTree == nil || len(nextTree) == 0 {
					continue //数组为空，则剔除字段，否则会出现 xx:nil 返回到前端就是xxx:null
				}
				nestedVal, err := TransformTree(nextTree, sourceField, targetField)
				if err != nil {
					return nil, err
				}
				kv[targetField[index]] = nestedVal
			} else {
				kv[targetField[index]] = vv.Interface()
			}
		}
		target = append(target, kv)
	}
	return target, nil
}
