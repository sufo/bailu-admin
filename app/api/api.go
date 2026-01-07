/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package api

import (
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/api/admin/content"
	"github.com/sufo/bailu-admin/app/api/admin/mine"
	"github.com/sufo/bailu-admin/app/api/admin/monitor"
	"github.com/sufo/bailu-admin/app/api/admin/system"
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
