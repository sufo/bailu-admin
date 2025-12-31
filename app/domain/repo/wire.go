/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package repo

import "github.com/google/wire"

var RepoSet = wire.NewSet(
	//UserSet,
	//TransSet,
	//DeptSet,
	//PostSet,
	//RoleSet,
	//MenuSet,
	NewDictItemRepo,
	NewDictRepo,
	NewUserRepo,
	TransSet,
	NewDeptRepo,
	NewPostRepo,
	NewRoleRepo,
	NewMenuApiRepo,
	NewMenuRepo,
	NewOperationRecoderRepo,
	NewNoticeRepo,
	NewNoticeSendRepo,
	NewMsgUserConfigRepo,
	NewTaskLogRepo,
	NewTaskRepo,
	NewLoginLogRepo,
	NewSysConfigRepo,
	//file
	NewCategoryRepo,
	NewFileRepo,
)
