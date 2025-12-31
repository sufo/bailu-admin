/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

// 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @description: 批量创建文件夹
// @param: dirs ...string
// @return: err exception
func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			//log.L.Debug("create directory" + v)
			fmt.Println("create directory %s", v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				//log.L.Error("create directory"+v, zap.Any(" exception:", err))
				return err
			}
		}
	}
	return err
}

func ReadFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//defer file.Close()
	return file, err
}

func LoadYaml2Struct[T any](path string) (T, error) {
	absPath, _ := filepath.Abs(path)
	file, err := os.Open(absPath)
	if err != nil {
		var t T
		return t, err
	}
	defer file.Close()

	var data T
	d := yaml.NewDecoder(file)
	err = d.Decode(&data)
	return data, err
}

func LoadYaml(path string) (map[string]any, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data := make(map[string]any)
	err = yaml.Unmarshal(file, &data)
	return data, nil
}
func LoadFlatYaml[T any](path string) (map[string]T, error) {
	//Unmarshal to map
	data, err := LoadYaml(path)
	if err != nil {
		//panic("fail to load yaml , error " + err.Error())
		return nil, err
	}
	//flat
	//var dest map[string]T //assignment to entry in nil map
	var dest = make(map[string]T)
	FlatMap[T]("", data, dest)
	return dest, nil
}

//////////////////////////////////////////////////
////////////////// file size /////////////////////
//////////////////////////////////////////////////

var suffixes = [5]string{"B", "KB", "MB", "GB", "TB"}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// fs := float64( file.Size() )
// fmt.Println(file.Name(), HumanFileSize(fs), file.ModTime())
func HumanFileSize(size float64) string {
	//fmt.Println(size)
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	//fmt.Println(int(math.Floor(base)))
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}
