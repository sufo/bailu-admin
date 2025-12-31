/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package service

import (
	"github.com/google/wire"
	"bailu/app/service/message"
	"bailu/app/service/sys"
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
