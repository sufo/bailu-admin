/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"gorm.io/gorm"
)

func NewNoticeRepo(db *gorm.DB) *NoticeRepo {
	return &NoticeRepo{base.Repository[entity.Notice]{db}}
}

// 通知公告
type NoticeRepo struct {
	base.Repository[entity.Notice]
}

func NewNoticeSendRepo(db *gorm.DB) *NoticeSendRepo {
	return &NoticeSendRepo{base.Repository[entity.NoticeSend]{db}}
}

type NoticeSendRepo struct {
	base.Repository[entity.NoticeSend]
}

// 消息配置
func NewMsgUserConfigRepo(db *gorm.DB) *MsgUserConfigRepo {
	return &MsgUserConfigRepo{base.Repository[entity.MsgUserConfig]{db}}
}

// 消息个人配置
type MsgUserConfigRepo struct {
	base.Repository[entity.MsgUserConfig]
}
