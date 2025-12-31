/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc trace生成
 */

package utils

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	version string
	incrNum uint64
	pid     = os.Getegid()
)

func NewTraceId() string {
	return fmt.Sprintf("trace-id-%d-%s-%d",
		os.Getppid(),
		time.Now().Format("2006.01.02.15.04.05.999"),
		atomic.AddUint64(&incrNum, 1))
}
