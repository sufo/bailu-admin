/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 定时调度那些定时发送的通知公告
 */

package jobs

import (
	"github.com/sufo/bailu-admin/app/service/sys"
	"context"
)

type ScheduleNotice struct {
	NoticeSrv *sys.NoticeService
}

func NewScheduleNotice(NoticeSrv *sys.NoticeService) *ScheduleNotice {
	return &ScheduleNotice{NoticeSrv}
}

func (t *ScheduleNotice) Invoke(ctx context.Context, args map[string]any) (result string, err error) {
	t.NoticeSrv.ScheduleNotice(ctx)
	return "定时通知调度成功", nil
}

func (c *ScheduleNotice) Name() string {
	return "定时通知"
}
