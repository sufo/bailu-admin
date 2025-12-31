/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 参数检验辅助
 */

package utils

import (
	"github.com/dlclark/regexp2"
	"regexp"
)

/** 密码正则(密码为6-18位数字/字符/符号的组合) */
var REGEXP_PWD = regexp2.MustCompile(`^(?![0-9]+$)(?![a-z]+$)(?![A-Z]+$)(?!([^(0-9a-zA-Z)]|[()])+$)(?!^.*[\u4E00-\u9FA5].*$)([^(0-9a-zA-Z)]|[()]|[a-z]|[A-Z]|[0-9]){6,18}$`, regexp2.RE2)

// 密码检验
func PasswordStrength(pwd string) bool {
	isMatch, _ := REGEXP_PWD.MatchString(pwd)
	return isMatch
}

// mysql 1062 //Duplicate entry 'sufo' for key 'sys_user.user_name'
var REGEXP_1062 = regexp.MustCompile(`^Duplicate entry [\'](\S+)[\'] for key [\'](.*?)[\']$`)
