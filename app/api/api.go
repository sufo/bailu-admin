/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package api

import (
	"github.com/google/wire"
	"bailu/app/api/admin/content"
	"bailu/app/api/admin/mine"
	"bailu/app/api/admin/monitor"
	"bailu/app/api/admin/system"
)

var APISet = wire.NewSet(
	system.NewUploadApi,
	system.DictItemSet,
	system.DictSet,
	system.AuthSet,
	system.UserSet,
	system.ProfileSet,
	system.RoleSet,
	system.MenuSet,
	system.DeptSet,
	system.PostSet,
	system.NoticeSet,
	monitor.NewOperApi,
	monitor.NewLoginLogApi,
	monitor.NewServerInfo,
	monitor.NewOnlineUserApi,
	monitor.NewTaskApi,
	monitor.NewEvent,
	system.NewSysConfigApi,
	mine.MessageSet,
	content.NewFileApi,
)
