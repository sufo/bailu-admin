/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package service

import (
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/service/message"
	"github.com/sufo/bailu-admin/app/service/sys"
)

var ServiceSet = wire.NewSet(
	sys.OnlineSet,
	sys.DictItemSet,
	sys.DictSet,
	sys.SMSSet,
	sys.MenuSet,
	sys.DeptSet,
	sys.RoleSet,
	sys.PostSet,
	sys.UserSet,
	sys.AuthSet,
	sys.OperationSet,
	sys.NoticeSrvSet,
	message.MessageSrvSet,
	message.NoticeSrvSet,
	sys.NewLoginLogService,
	//message.NoticeSrvSet, //import cycle not allowed
	sys.SysConfigSet,
	sys.FileSet,
)
