/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Update 2024/4/30
 * @Desc 通知公告
 */

package entity

import (
	"gorm.io/gorm"
	"time"
)

var _ IModel = (*Notice)(nil)

// 通知公告
type Notice struct {
	ID            uint64         `json:"id" gorm:"primarykey"`
	Type          int            `json:"type" gorm:"type:tinyint(4);comment: 1通知，2公告(公告只能全体)"`
	Title         string         `json:"title" gorm:"not null;size:100;comment:通知标题"`
	Content       string         `json:"content" gorm:"type:text;default:null;comment:内容"`
	Sender        string         `json:"sender" gorm:"size:30;default:null;comment:发布人"`
	SenderId      uint64         `json:"senderId" gorm:"default:null;comment:发布人id,0表示是系统"`
	Receivers     string         `json:"receivers" gorm:"comment:接收者id(id类型由Group_type决定，多个以逗号分割)"`
	SendScope     string         `json:"sendScope" gorm:"size:30;default:null;comment:发送范围类型（all所有人、user指定用户、role角色、depart部门等），见字典message_scope_type"`
	SendStatus    string         `json:"sendStatus" gorm:"size:1;default:'0';comment:发布状态（0未发布，1已发布，2已撤销）"`
	StartTime     time.Time `json:"startTime" gorm:"default:null;comment:开始时间"`
	EndTime       time.Time `json:"endTime" gorm:"default:null;comment:结束时间"`
	NotifyChannel string         `json:"notifyChannel" gorm:"default:'web';size:20;default:'web';comment:通知方式：web,app,mail,sms等"`
	ScheduledTime time.Time `json:"scheduledTime" gorm:"default:null;comment:指定发送时间。如果该通知/公告是特定时间发送的"`
	SendTime      time.Time `json:"sendTime" gorm:"default:null;comment:发布时间"`
	CancelTime    time.Time `json:"cancelTime" gorm:"default:null;comment:撤销时间"`
	Icon          string         `json:"icon" gorm:"-"`
	ReadFlag      int            `json:"readFlag" gorm:"-"`
	ReadTime      time.Time `json:"readTime" gorm:"-"`
	UpdateBy      uint64         `json:"-" gorm:"column:update_by;default:0;comment:更新者"`
	CreatedAt     time.Time `json:"createdAt"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

var NoticeTN = "msg_notice"

func (m Notice) GetID() uint64 {
	return m.ID
}

func (Notice) TableName() string {
	return NoticeTN
}

// 通知公告读取状态
// 对于这种广播类的消息，并不是为每一个接收用户在notice_send表中插入一条数据，而是只有当用户处理过该消息后（status=已读||删除）才将这一行插入notice_status表
type NoticeSend struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	MsgId     uint64 `json:"msgId" gorm:"index;comment:消息id"`
	ReceiveId uint64 `json:"receiveId" gorm:"index;comment:接收者id"`
	ReadFlag  int    `json:"readFlag" gorm:"type:tinyint(4);default:0;comment:阅读状态（0未读，1已读, 2删除）"`
	//StarFlag  uint8     `json:"starFlag" gorm:"type:tinyint(4);comment:是否标星,1是标星消息"`
	ReadTime  time.Time `json:"readTime" gorm:"default:null;comment:查看时间或者删除时间"`
	CreatedAt time.Time `json:"createdAt"`
}

var NoticeSendTN = "msg_notice_send"

func (m NoticeSend) GetID() uint64 {
	return m.ID
}

func (NoticeSend) TableName() string {
	return NoticeSendTN
}

// 通知订阅
type RemindSubscription struct {
	ID         uint64         `json:"id" gorm:"primarykey"`
	SourceType string         `json:"sourceType" gorm:"size:50; comment:目标类型（文章等）"`
	SourceId   uint64         `json:"sourceId" gorm:"comment:目标id"`
	Action     string         `json:"action" gorm:"size:50;comment:动作类型（ 1、点赞  2、评论 3、回复 4、@  5、关注' 等）"`
	UserId     uint64         `json:"userId" gorm:"comment:订阅用户"`
	CreatedAt  time.Time `json:"createdAt"`
}

var RemindSubscriptionTN = "remind_subscription"

func (n RemindSubscription) GetID() uint64 {
	return n.ID
}

func (RemindSubscription) TableName() string {
	return RemindSubscriptionTN
}

// 通知配置
//type NotifyConfig struct {
//	Channel string `json:"channel" gorm:"size:30;comment:通知渠道"`
//}

// 用户自定义消息配置配置
// 主要针对提醒类型的消息的
type MsgUserConfig struct {
	ID            uint64 `json:"id" gorm:"primarykey"`
	SwitchBit     int    `json:"switchBit" gorm:"comment:位开关，第一位总开关，第二位：点赞，后面依次为评论、回复、@、关注"`
	NotifyChannel string `json:"notifyChannel" gorm:"size:20;default:'web';comment:通知方式：web,app,mail,sms等"`
	UserId        uint64 `json:"userId" gorm:"comment:用户"`
}

var MsgUserConfigTN = "msg_user_config"

func (n MsgUserConfig) GetID() uint64 {
	return n.ID
}

func (MsgUserConfig) TableName() string {
	return MsgUserConfigTN
}

//https://cloud.tencent.com/developer/article/1684449
