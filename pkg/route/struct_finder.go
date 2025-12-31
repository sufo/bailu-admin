/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package route

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// StructFinder 用于存储结构体查找的结果
type StructFinder struct {
	targetStruct string
	foundFile    string
}

// Visit 实现ast.Visitor接口
func (f *StructFinder) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	// 检查是否为类型声明
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		// 检查是否为结构体定义
		if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
			// 如果结构体名称匹配
			if typeSpec.Name.Name == f.targetStruct {
				// 获取文件对象
				file := node.Pos().IsValid()
				if file {
					f.foundFile = f.targetStruct //可随便赋值表示找到了
				}
			}
		}
	}
	return f
}

func findStructInFile(filePath string, structName string) (bool, error) {
	// 创建文件集
	fset := token.NewFileSet()

	// 解析Go源文件
	file, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return false, fmt.Errorf("解析文件失败: %v", err)
	}

	// 创建结构体查找器
	finder := &StructFinder{
		targetStruct: structName,
	}

	// 遍历AST
	ast.Walk(finder, file)

	return finder.foundFile != "", nil
}

//该方法也可以使用
//func findStructInFile(filePath string, structName string) (bool, error) {
//
//	// 创建结构体查找器
//	finder := &StructFinder{
//		targetStruct: structName,
//	}
//	// 创建文件集
//	fset := token.NewFileSet()
//
//	// 解析Go源文件
//	node, err := parser.ParseFile(fset, filePath, nil, 0)
//	if err != nil {
//		return false, fmt.Errorf("解析文件失败: %v", err)
//	}
//
//	// Inspect the file's declarations to find struct types.
//	for _, decl := range node.Decls {
//		genDecl, ok := decl.(*ast.GenDecl)
//		if !ok || genDecl.Tok != token.TYPE {
//			continue
//		}
//
//		// Find struct declarations.
//		for _, spec := range genDecl.Specs {
//			typeSpec, ok := spec.(*ast.TypeSpec)
//			if !ok {
//				continue
//			}
//
//			// Check if the name matches the struct we're looking for.
//			if typeSpec.Name.Name == structName {
//				finder.foundFile = filePath
//				return finder.foundFile != "", nil // Stop after finding the file.
//			}
//		}
//	}
//	return finder.foundFile != "", nil
//}

func FindStruct(dir string, structName string) (string, error) {
	var foundFile string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理Go源文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			found, err := findStructInFile(path, structName)
			if err != nil {
				return err
			}
			if found {
				foundFile = path
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("查找过程出错: %v", err)
	}

	if foundFile == "" {
		return "", fmt.Errorf("未找到结构体 %s", structName)
	}

	return foundFile, nil
}
