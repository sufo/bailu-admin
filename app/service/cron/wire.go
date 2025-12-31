/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package cron

import (
	"bailu/app/service/cron/jobs"
	"github.com/google/wire"
)

var CrontabSet = wire.NewSet(
	jobs.NewScheduleNotice,
	Inject2JobSet,
	NewCronTask,
	TaskSet,
)
