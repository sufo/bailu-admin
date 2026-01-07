/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 基于swagger.json解析
 */

package route

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sufo/bailu-admin/docs"
	"strings"
)

// RouteInfo 路由信息结构体
type routeInfo struct {
	Method     string
	Path       string
	Summary    string //处理函数简短摘要 @Summary
	Tag        string //@Tag
	Deprecated bool   // 是否已废弃
}

type Option struct {
	Key      string           `json:"key"`
	Label    string           `json:"label"`
	Children []map[string]any `json:"children"`
}

// 获取api tree
// Parse Swagger.json
func RouteTree(engine *gin.Engine, excludeRoutePaths ...string) ([]map[string]any, error) {

	if docs.SwaggerJson == "" {
		return nil, fmt.Errorf("Swagger is not generated. Please execute swag init on the command line to generate it.")
	}

	var swaggerMap map[string]any
	err := json.Unmarshal([]byte(docs.SwaggerJson), &swaggerMap)
	if err != nil {
		return nil, err
	}

	var category = make(map[string]*Option)
	pathsMap := swaggerMap["paths"].(map[string]any)
	for url, value := range pathsMap {
		apiMap, ok := value.(map[string]any)
		if !ok {
			continue
		}
		//处理包含{xxx}的动态路由,将{xxx} -> :xxx
		path := strings.ReplaceAll(strings.ReplaceAll(url, "{", ":"), "}", "")
		//是否跳过
		if isSkip(path, excludeRoutePaths...) {
			continue
		}

		for method, val := range apiMap {
			reqMap := val.(map[string]any)
			summary := reqMap["summary"].(string)
			tags := reqMap["tags"].([]any)
			var tag string
			if len(tags) == 0 {
				tag = summary
			} else {
				tag = tags[0].(string)
			}

			item := map[string]any{
				"key":    method + "_" + path,
				"label":  summary,
				"method": method,
				"path":   path,
				"isLeaf": true,
			}

			if opt, exist := category[tag]; exist {
				opt.Children = append(opt.Children, item)
			} else {
				category[tag] = &Option{
					Key:      path,
					Label:    tag,
					Children: []map[string]any{item},
				}
			}
		}
	}

	var result = make([]map[string]any, 0)
	for _, value := range category {
		result = append(result, map[string]any{"key": value.Key, "label": value.Label, "children": value.Children})
	}
	return result, nil
}

func isSkip(path string, excludes ...string) bool {
	for _, eachItem := range excludes {
		if strings.HasPrefix(path, eachItem) {
			return true
		}
	}
	return false
}
