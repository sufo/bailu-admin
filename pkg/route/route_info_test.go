/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"testing"
)

func TestSignature_Generate(t *testing.T) {
	// 创建 Gin 实例
	r := gin.New()

	// 创建文档解析器
	parser, err := NewDocParser()
	if err != nil {
		panic(err)
	}

	// 注册路由和处理函数
	r.GET("/users", GetUserList)
	r.POST("/users", CreateUser)

	// 解析并打印文档
	printHandlerDoc(parser, "/users", "GET", GetUserList)
	printHandlerDoc(parser, "/users", "POST", CreateUser)
}

// 打印处理函数文档
func printHandlerDoc(parser *DocParser, path, method string, handler gin.HandlerFunc) {
	doc, err := parser.ParseHandlerDoc(handler)
	if err != nil {
		fmt.Printf("Error parsing doc for %s %s: %v\n", method, path, err)
		return
	}

	doc.Path = path
	doc.Method = method

	fmt.Printf("API Documentation for %s %s\n", method, path)
	fmt.Println(strings.Repeat("=", 80))

	if doc.Deprecated {
		fmt.Println("⚠️ DEPRECATED")
	}

	if doc.Summary != "" {
		fmt.Printf("Summary: %s\n", doc.Summary)
	}

	if doc.Description != "" {
		fmt.Printf("\nDescription:\n%s\n", doc.Description)
	}

	if len(doc.Parameters) > 0 {
		fmt.Printf("\nParameters:\n")
		for _, param := range doc.Parameters {
			required := ""
			if param.Required {
				required = " (required)"
			}
			fmt.Printf("  - %s (%s)%s: %s\n",
				param.Name, param.Type, required, param.Description)
		}
	}

	if len(doc.Returns) > 0 {
		fmt.Printf("\nReturns:\n")
		for _, ret := range doc.Returns {
			fmt.Printf("  - %s\n", ret)
		}
	}

	if len(doc.Tags) > 0 {
		fmt.Printf("\nAdditional Tags:\n")
		for key, value := range doc.Tags {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Println(strings.Repeat("-", 80))
}

// GetUserList 获取用户列表
// @summary 获取系统中的所有用户
// 支持分页和过滤
// @param page int 页码 (默认1)
// @param size int 每页数量 (默认20)
// @param status string 用户状态过滤
// @return 用户列表及分页信息
// @version 1.0
func GetUserList(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get user list"})
}

// CreateUser 创建新用户
// @summary 创建新用户账户
// 创建新用户并返回用户信息
// @param name string 用户名 required
// @param email string 邮箱地址 required
// @param role string 用户角色
// @return 新创建的用户信息
// @version 1.0
func CreateUser(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Create user"})
}
