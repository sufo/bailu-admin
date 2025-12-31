/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 不能兼容go run
func getCurrentPath() string {
	s, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println(err.Error())
	}
	s = strings.Replace(s, "\\", "/", -1)
	s = strings.Replace(s, "\\\\", "/", -1)
	i := strings.LastIndex(s, "/")
	path := string(s[0 : i+1])
	return path
}

// 最终方案-全兼容
func GetCurrentAbPath(skip int) string {
	dir := getCurrentAbPathByExcutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller(skip)
	}
	return dir
}

// 获取系统临时目录 兼容go run
func getTmpDir() string {
	//dir := os.Getenv("TEMP")
	//if dir == "" {
	//	dir = os.Getenv("TMP")
	//}
	//res, _ := filepath.EvalSymlinks(dir)
	//return res
	res, _ := filepath.EvalSymlinks(os.TempDir())
	return res
}

// 获取当前执行文件的绝对路径
func getCurrentAbPathByExcutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller(skip int) string {
	var abPath string
	_, filename, _, ok := runtime.Caller(skip)
	if ok {
		abPath = path.Dir((filename))
	}
	return abPath
}
