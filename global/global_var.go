/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package global

import "time"

var StartTime time.Time //记录启动时间
var Root string         //项目根目录
var Version string      //项目版本

// 三元表达式
func Ternary[T any](expr bool, a, b T) T {
	if expr {
		return a
	}
	return b
}
