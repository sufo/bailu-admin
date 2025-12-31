/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 清理日志
 */

package jobs

import (
	"bailu/app/config"
	"bailu/utils/time"
	"context"
	"os"
	"path"
	"strings"
)

type CleanLogJob struct {
}

func (c *CleanLogJob) Invoke(ctx context.Context, args map[string]any) (result string, err error) {
	//默认清理一个月之前的日志
	daysAgo, exist := args["daysAgo"]
	if !exist {
		daysAgo = config.Conf.Zap.CleanDaysAgo
	}
	files, err := os.ReadDir(config.Conf.Zap.Director)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		arr := strings.Split(file.Name(), ".")
		if time.IsNDaysAgo(daysAgo.(int), arr[0]) {
			if err = os.Remove(path.Join(config.Conf.Zap.Director, file.Name())); err != nil {
				return "", err
			}
		}
	}
	return "处理成功", nil
}

func (c *CleanLogJob) Name() string {
	return "清理日志"
}
