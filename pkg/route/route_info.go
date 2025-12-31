/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package route

import (
	"bailu/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// HandlerDoc 处理函数文档结构
type HandlerDoc struct {
	Summary     string            // 接口摘要
	Description string            // 详细描述
	Path        string            // 路由路径
	Method      string            // HTTP方法
	Parameters  []ParameterDoc    // 参数文档
	Returns     []string          // 返回说明
	Deprecated  bool              // 是否已废弃
	Tags        map[string]string // 其他标签信息
}

// ParameterDoc 参数文档结构
type ParameterDoc struct {
	Name        string // 参数名
	Type        string // 参数类型
	Description string // 参数描述
	Required    bool   // 是否必须
}

// DocParser 文档解析器
type DocParser struct {
	projectRoot string               // 项目根目录
	fileCache   map[string]*ast.File // 文件解析缓存
	fset        *token.FileSet       // 文件集
}

// NewDocParser 创建文档解析器
func NewDocParser() (*DocParser, error) {
	// 获取项目根目录
	projectRoot := global.Root
	var err error
	if projectRoot == "" {
		projectRoot, err = findProjectRoot()
		if err != nil {
			return nil, err
		}
	}

	return &DocParser{
		projectRoot: projectRoot,
		fileCache:   make(map[string]*ast.File),
		fset:        token.NewFileSet(),
	}, nil
}

// findProjectRoot 查找项目根目录（包含 go.mod 的目录）
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find project root (go.mod)")
		}
		dir = parent
	}
}

// ParseHandlerDoc 解析处理函数的文档
func (p *DocParser) ParseHandlerDoc(handler gin.HandlerFunc) (*HandlerDoc, error) {
	// 获取处理函数的包路径和函数名
	handlerValue := reflect.ValueOf(handler)
	if handlerValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("not a function")
	}

	// 获取函数名和包路径
	fullName := runtime.FuncForPC(handlerValue.Pointer()).Name()
	pkgPath, structName, funcName := splitFuncName(fullName)

	// 查找源文件
	filePath, err := p.findSourceFile(pkgPath, structName)
	if err != nil {
		return nil, err
	}

	// 解析源文件
	file, err := p.parseFile(filePath)
	if err != nil {
		return nil, err
	}

	// 查找函数声明并解析注释
	return p.parseFuncDoc(file, strings.Replace(funcName, "-fm", "", 1))
}

// splitFuncName 分割完整函数名为包路径和函数名
func splitFuncName(fullName string) (string, string, string) {
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 {
		return "", "", fullName
	}
	pkg := fullName[:lastDot]
	var strcutName = ""
	// 处理方法接收器
	if strings.Contains(pkg, ")") {
		arr := strings.Split(pkg, "(")
		pkg = arr[0]
		strcutName = strings.Split(arr[1], ")")[0]
		if strings.HasPrefix(strcutName, "*") { //如果名字以*开头
			strcutName = strcutName[1:]
		}
	}
	return pkg[:len(pkg)-1], strcutName, fullName[lastDot+1:]
}

//func (p *DocParser) findSourceFile(pkgPath, structName string) (string, error) {
//	start := strings.Index(pkgPath, "/")
//	pkgDir := filepath.Join(p.projectRoot, strings.Replace(pkgPath[start+1:], "/", string(os.PathSeparator), -1))
//	entries, err := os.ReadDir(pkgDir)
//	if err != nil {
//		return "", err
//	}
//
//	// 根据结构体名称查找对应 .go 文件
//	// 这里必须 文件名包含”结构体名除去Api“ 这样的规律才行
//	for _, entry := range entries {
//		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
//			name := strings.ToLower(strings.Replace(structName, "Api", "", -1))
//			//文件名跟结构体命名方式不一样，所以去掉_分割符，为了跟结构体去匹配。
//			entryName := strings.Split(entry.Name(), ".")[0]
//			fName := strings.Replace(entryName, "_", "", -1)
//			if strings.Contains(name, fName) || strings.Contains(fName, name) {
//				return filepath.Join(pkgDir, entry.Name()), nil
//			}
//		}
//	}
//	return "", fmt.Errorf("no source file found")
//}

// 根据structName查找go文件路径
// @return file path
func (p *DocParser) findSourceFile(pkgPath, structName string) (string, error) {
	start := strings.Index(pkgPath, "/")
	pkgDir := filepath.Join(p.projectRoot, strings.Replace(pkgPath[start+1:], "/", string(os.PathSeparator), -1))

	filePath, err := FindStruct(pkgDir, structName)
	if err != nil {
		fmt.Printf("查找文件错误: %v\n", err)
		return "", err
	}
	return filePath, nil
}

// parseFile 解析源文件
func (p *DocParser) parseFile(filePath string) (*ast.File, error) {
	if file, ok := p.fileCache[filePath]; ok {
		return file, nil
	}

	file, err := parser.ParseFile(p.fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	p.fileCache[filePath] = file
	return file, nil
}

// parseFuncDoc 解析函数文档
func (p *DocParser) parseFuncDoc(file *ast.File, funcName string) (*HandlerDoc, error) {
	doc := &HandlerDoc{
		Tags:       make(map[string]string),
		Parameters: make([]ParameterDoc, 0),
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == funcName {
				if funcDecl.Doc != nil {
					p.parseComments(funcDecl.Doc.List, doc)
				}
				return false
			}
		}
		return true
	})

	return doc, nil
}

// parseComments 解析注释内容
func (p *DocParser) parseComments(comments []*ast.Comment, doc *HandlerDoc) {
	var description []string

	for _, comment := range comments {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

		// 解析特殊标记
		switch {
		case strings.HasPrefix(text, "@summary"):
			doc.Summary = strings.TrimSpace(strings.TrimPrefix(text, "@summary"))
		case strings.HasPrefix(text, "@Summary"):
			doc.Summary = strings.TrimSpace(strings.TrimPrefix(text, "@Summary"))
		case strings.HasPrefix(text, "@param"):
			param := parseParameter(text)
			doc.Parameters = append(doc.Parameters, param)
		case strings.HasPrefix(text, "@return"):
			returnDesc := strings.TrimSpace(strings.TrimPrefix(text, "@return"))
			doc.Returns = append(doc.Returns, returnDesc)
		case strings.HasPrefix(text, "@deprecated"):
			doc.Deprecated = true
		case strings.HasPrefix(text, "@"):
			// 处理其他自定义标记
			if parts := strings.SplitN(text[1:], " ", 2); len(parts) == 2 {
				doc.Tags[parts[0]] = strings.TrimSpace(parts[1])
			}
		default:
			description = append(description, text)
		}
	}

	doc.Description = strings.Join(description, "\n")
}

// parseParameter 解析参数注释
func parseParameter(text string) ParameterDoc {
	// 移除 @param 前缀
	text = strings.TrimPrefix(text, "@param")
	parts := strings.Fields(text)

	param := ParameterDoc{
		Required: strings.Contains(text, "required"),
	}

	if len(parts) >= 1 {
		param.Name = parts[0]
	}
	if len(parts) >= 2 {
		param.Type = parts[1]
	}
	if len(parts) >= 3 {
		param.Description = strings.Join(parts[2:], " ")
	}

	return param
}
