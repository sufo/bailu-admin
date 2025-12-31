/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc ip2region
 */

package ip2region

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

const DbPath = "./assets/ip2region.xdb"

// 1、从 dbPath 加载 VectorIndex 缓存，把下述 vIndex 变量全局到内存里面。
var VIndex []byte

func InitIp2Region() {
	vIndex, err := xdb.LoadVectorIndexFromFile(DbPath)
	if err != nil {
		fmt.Printf("failed to load vector index from `%s`: %s\n", DbPath, err)
		return
	}
	VIndex = vIndex
}
