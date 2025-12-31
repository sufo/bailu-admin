/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 自定义json输出时间格式
 * 在 gorm 中使用覆盖 MarshalJSON 的方式来自定义时间字段的格式时，只重写 MarshalJSON 是不够的，
 * 只写这个方法会在写数据库的时候会提示 delete_at 字段不存在，还需要加上对 database/sql 中的 Value 和 Scan 方法的实现。
 */

package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const CSTLayout = "2006-01-02 15:04:05"

//type JSONTime time.Time  //别名方式

// 内嵌方式（推荐）
// JSONTime format json time field by myself
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	//formatted := fmt.Sprintf("\"%s\"", t.Format(CSTLayout))
	formatted := fmt.Sprintf("%q", t.Format(CSTLayout))
	return []byte(formatted), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse(`"2006-01-02 15:04:05"`, string(data)) //layout使用CSTLayout变量会报错
	*t = JSONTime{tt}
	return err
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
