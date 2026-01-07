/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 通知公告
 */

package sys

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/wire"
	"bailu/app/config"
	"bailu/app/core/appctx"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	"bailu/global"
	"bailu/pkg/log"
	"bailu/pkg/mq"
	"bailu/pkg/store"
	"bailu/utils"
	"strings"
	"time"
)

const (
	SCOPE_ALL    = "all"
	SCOPE_ROLE   = "role"
	SCOPE_DEPART = "depart"
	SCOPE_USER   = "user"
)

const (
	CHANNEL_MAIL = "mail"
	CHANNEL_SMS  = "sms"
	CHANNEL_WEB  = "web"
	CHANNEL_APP  = "app"
)

var NoticeSrvSet = wire.NewSet(wire.Struct(new(NoticeService), "*"))

type NoticeService struct {
	NoticeRepo     *repo.NoticeRepo
	NoticeSendRepo *repo.NoticeSendRepo
	UserRepo       *repo.UserRepo
	TranRepo       *repo.Trans
	////OnlineSrv      *service.OnlineService
	DB store.IStore
}

// 记录定时器，方便取消
var timers = make(map[uint64]*time.Timer, 0)

// 查询我发送的通知
func (a *NoticeService) List(ctx context.Context, params dto.NoticeParams) (*resp.PageResult[entity.Notice], error) {
	builder := base.NewQueryBuilder().WithPagination(ctx)
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if user.IsSuper() { //查所有
		builder.WithWhere("title like ?", fmt.Sprint("%", params.Title, "%"))
	} else {
		builder.WithWhere("title like ? and senderId=?", params.Title, user.ID) //查自己发送的
	}
	if params.SendStatus != "" {
		builder.WithWhere("send_status = ?", params.SendStatus)
	}
	if params.IfScheduled != "" {
		scheduled, err := utils.ToInt[int32](params.IfScheduled)
		if err != nil {
			return nil, err
		}
		query := global.Ternary(scheduled == 1, "scheduled_time IS NOT NULL", "scheduled_time IS NULL")
		builder.WithWhere(query)
	}
	if params.ReadStatus != nil {
		builder.WithJoin(fmt.Sprintf("left join %s as ns on ns.msgId=%s.id", entity.NoticeSendTN, entity.NoticeTN)).
			WithWhere("ns.read_flag=?", *params.ReadStatus)
	}

	if result, err := a.NoticeRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		return result.(*resp.PageResult[entity.Notice]), err
	}
}

// 查询我发送的通知
func (a *NoticeService) MySendedList(ctx context.Context, title string, status *int) (*resp.PageResult[entity.Notice], error) {
	builder := base.NewQueryBuilder().WithPagination(ctx)
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if user.IsSuper() { //查所有
		builder.WithWhere("title like ?", fmt.Sprint("%", title, "%"))
	} else {
		builder.WithWhere("title like ? and senderId=?", user.ID) //查自己发送的
	}
	if status != nil {
		builder.WithJoin(fmt.Sprintf("left join %s as ns on ns.msgId=%s.id", entity.NoticeSendTN, entity.NoticeTN)).
			WithWhere("ns.read_flag=?", *status)
	}
	if result, err := a.NoticeRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		return result.(*resp.PageResult[entity.Notice]), err
	}
}

func (a *NoticeService) Create(ctx context.Context, notice *entity.Notice) error {
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	notice.SenderId = user.ID
	notice.Sender = user.Username //保存这个字段，避免连表查询
	return a.NoticeRepo.Create(ctx, notice)
}

// 修改(必须撤销之后才可以修改)
func (a *NoticeService) Update(ctx context.Context, notice *entity.Notice) (*entity.Notice, error) {
	err := a.NoticeRepo.Update(ctx, notice)
	return notice, err
}

// 删除
func (a *NoticeService) Delete(ctx context.Context, ids []uint64) error {
	return a.TranRepo.Exec(ctx, func(ctx context.Context) error {
		//删除send表
		err := a.NoticeSendRepo.Where(ctx, "msg_id in ?", ids).Unscoped().Delete(&entity.NoticeSend{}).Error
		if err != nil {
			return err
		}
		return a.NoticeRepo.Delete(ctx, ids)
	})
}

// 用户删除自己通知公告
func (a *NoticeService) UserDel(ctx context.Context, msgIds []uint64) error {
	return a.NoticeSendRepo.Where(ctx, "msg_id in ?", msgIds).
		Update("read_time", time.Now()).Update("read_flag", 2).Error
}

// // 发布通知
//
//	func (n *NoticeService) releaseNotice(ctx context.Context, notice *entity.Notice) error {
//		var err error
//		user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
//		notice.SenderId = user.ID
//		notice.SendStatus = "1"
//		notice.SendTime = time.Now() //这里设置不太准确，因为消息会先进入队列
//
//		//插入send
//		var noticeSends = make([]entity.NoticeSend, 0)
//
//		//根据发送发送方式采取不同发送策略
//		//选择近2个月登录过系统的用户，如果是web、app通知的方式，则只推送给在线的用户
//		send_channel := notice.NotifyChannel
//		if send_channel == CHANNEL_MAIL || send_channel == CHANNEL_SMS {
//			uIds := make([]uint64, 0)
//			//只需要id
//			builder := (&base.QueryBuilder{}).WithSelect("id")
//			switch notice.SendScope {
//			case SCOPE_ALL:
//				builder.WithWhere("status=? and DATE_SUB(CURDATE(),INTERVAL 60 DAY) <= date(last_login_time)", 1)
//				err := n.UserRepo.Where(ctx, "status=? and DATE_SUB(CURDATE(),INTERVAL 60 DAY) <= date(last_login_time)", 1).Find(&uIds).Error
//				if err != nil {
//					return err
//				}
//			case SCOPE_ROLE:
//				receiversStr := strings.Split(notice.Receivers, ",")
//				ids, _ := utils.StrArr2Arr[uint64](receiversStr)
//				builder.WithJoin("left join sys_user_role on r on ur.user_id= sys_user.id").
//					WithWhere("ur.role_id in ?", ids).
//					WithWhere("status=? and DATE_SUB(CURDATE(),INTERVAL 60 DAY) <= date(last_login_time)", 1)
//			case SCOPE_DEPART:
//				receiversStr := strings.Split(notice.Receivers, ",")
//				ids, _ := utils.StrArr2Arr[uint64](receiversStr)
//				builder.WithWhere("dept_id in ?", ids)
//			case SCOPE_USER:
//				receiversStr := strings.Split(notice.Receivers, ",")
//				unames, _ := utils.StrArr2Arr[string](receiversStr)
//				builder.WithWhere("username in ?", unames)
//			}
//			err := n.UserRepo.WithQueryBuilder(ctx, builder).Find(&uIds).Error
//			if err != nil {
//				return err
//			}
//
//			//notice send数据
//			for _, id := range uIds {
//				noticeSends = append(noticeSends, entity.NoticeSend{MsgId: notice.ID, ReceiveId: id})
//			}
//		} else {
//			users := make([]*entity.OnlineUserDto, 0)
//			//onlineUsers, err := n.OnlineSrv.GetAll(ctx)
//			onlineUsers, err := n.getOnlineUsers()
//			if err != nil {
//				return err
//			}
//
//			switch notice.SendScope {
//			case SCOPE_ALL:
//				users = onlineUsers
//			case SCOPE_ROLE:
//				for _, online := range onlineUsers {
//					for _, role := range online.Roles {
//						if strings.Contains(notice.Receivers, utils.Strval(role.ID)) {
//							users = append(users, online)
//							break
//						}
//					}
//				}
//			case SCOPE_DEPART:
//				for _, online := range onlineUsers {
//					if strings.Contains(notice.Receivers, utils.Strval(online.DeptId)) {
//						users = append(users, online)
//					}
//				}
//			case SCOPE_USER:
//				for _, online := range onlineUsers {
//					if strings.Contains(notice.Receivers, utils.Strval(online.ID)) {
//						users = append(users, online)
//					}
//				}
//			}
//
//			//notice send数据
//			for _, user := range users {
//				noticeSends = append(noticeSends, entity.NoticeSend{MsgId: notice.ID, ReceiveId: user.ID})
//			}
//
//		}
//
//		err = n.TranRepo.Exec(ctx, func(ctx context.Context) error {
//			err := n.NoticeRepo.Update(ctx, notice)
//			if err != nil {
//				return err
//			}
//			//批量插入notice_send数据
//			return n.NoticeSendRepo.Create(ctx, noticeSends)
//		})
//		if err != nil {
//			return err
//		}
//		//发送消息
//		for _, notice := range noticeSends {
//			_, err := mq.Publisher.SendMsg(notice)
//			if err != nil {
//				log.L.Error(err)
//			}
//		}
//		return nil
//	}
//

// 发布通知
// 简化上面的逻辑，不针对在线用户推送
func (n *NoticeService) releaseNotice(ctx context.Context, notice *entity.Notice) error {
	var err error
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	notice.SenderId = user.ID
	notice.SendStatus = "1"
	notice.SendTime = time.Now() //这里设置不太准确，因为消息会先进入队列

	//插入send
	var noticeSends = make([]entity.NoticeSend, 0)

	//根据发送发送方式采取不同发送策略
	//选择近2个月登录过系统的用户，如果是web、app通知的方式，则只推送给在线的用户
	uIds := make([]uint64, 0)
	//只需要id
	builder := (&base.QueryBuilder{}).WithSelect("id")
	switch notice.SendScope {
	case SCOPE_ALL:
		builder.WithWhere("status=? and DATE_SUB(CURDATE(),INTERVAL 60 DAY) <= date(last_login_time)", 1)
		err := n.UserRepo.Where(ctx, "status=? and DATE_SUB(CURDATE(),INTERVAL 60 DAY) <= date(last_login_time)", 1).Find(&uIds).Error
		if err != nil {
			return err
		}
	case SCOPE_ROLE:
		receiversStr := strings.Split(notice.Receivers, ",")
		ids, _ := utils.StrArr2Arr[uint64](receiversStr)
		builder.WithJoin("left join sys_user_role on r on ur.user_id= sys_user.id").
			WithWhere("ur.role_id in ?", ids).
			WithWhere("status=? and DATE_SUB(CURDATE(),INTERVAL 60 DAY) <= date(last_login_time)", 1)
	case SCOPE_DEPART:
		receiversStr := strings.Split(notice.Receivers, ",")
		ids, _ := utils.StrArr2Arr[uint64](receiversStr)
		builder.WithWhere("dept_id in ?", ids)
	case SCOPE_USER:
		receiversStr := strings.Split(notice.Receivers, ",")
		unames, _ := utils.StrArr2Arr[string](receiversStr)
		builder.WithWhere("username in ?", unames)
	}
	err = n.UserRepo.WithQueryBuilder(ctx, builder).Find(&uIds).Error
	if err != nil {
		return err
	}

	//notice send数据
	for _, id := range uIds {
		noticeSends = append(noticeSends, entity.NoticeSend{MsgId: notice.ID, ReceiveId: id})
	}

	err = n.TranRepo.Exec(ctx, func(ctx context.Context) error {
		err := n.NoticeRepo.Update(ctx, notice)
		if err != nil {
			return err
		}
		//批量插入notice_send数据
		return n.NoticeSendRepo.Create(ctx, noticeSends)
	})
	if err != nil {
		return err
	}

	//发送消息
	// 如果是web端推送消息，这里不用筛选在线用户去发送，直接使用近2个月登录过系统的用户就行
	// 如果用户不在线，下次进入系统会拉取自己的未读消息
	for _, notice := range noticeSends {
		_, err := mq.Publisher.SendMsg(notice)
		if err != nil {
			log.L.Error(err)
		}
	}
	return nil
}

// 发布
// 暂时只考虑超级管理员
func (n *NoticeService) ReleaseNotice(ctx context.Context, id uint64) error {

	notice, err := n.NoticeRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	//判断是否是定时发送任务
	specialSendTime := notice.ScheduledTime
	if !specialSendTime.IsZero() {
		if specialSendTime.Before(time.Now()) {
			return errors.New("定时发送时间不能早于当前时间")
		} else {
			_timer := time.AfterFunc(specialSendTime.Sub(time.Now()), func() {
				err := n.releaseNotice(ctx, notice)
				log.L.Error(err)
			})
			timers[notice.ID] = _timer //记录_timer
			return nil
		}
	} else { //非定时执行的通知
		return n.releaseNotice(ctx, notice)
	}
}

// 撤销（硬删除）
func (n *NoticeService) CancelNotice(ctx context.Context, id uint64) error {
	return n.TranRepo.Exec(ctx, func(ctx context.Context) error {
		notice, err := n.NoticeRepo.FindById(ctx, id)
		if err != nil {
			return err
		}
		notice.CancelTime = time.Now()
		notice.SendStatus = "2"
		n.NoticeRepo.Update(ctx, notice)
		err = n.NoticeSendRepo.Where(ctx, "msg_id=?", notice.ID).
			Unscoped().Delete(&entity.NoticeSend{}).Error

		if err == nil {
			//如果是定时发送的通知/公告，则需要停止定时器
			_t, exist := timers[id]
			if exist {
				_t.Stop()
				delete(timers, id)
			}
		}
		return err
	})
}

// 任务执行完成发送通知
func (n *NoticeService) TaskNotice(ctx context.Context, taskModel *entity.Task, taskResult dto.TaskResult) error {
	var msg = global.Ternary[string](taskResult.Err == nil, "成功", "失败")
	//user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	//return n.NoticeRepo.Create(ctx, entity.Notice{
	//	Type:          1,
	//	Title:         fmt.Sprintf("任务[%s]执行<strong>%</strong>", taskModel.Name, msg),
	//	Content:       taskResult.Result,
	//	SenderId:      0, //系统
	//	ReceiveId:     strconv.FormatUint(user.ID, 64),
	//	SendScope:     SCOPE_USER,
	//	SendStatus:    "1", //已发布
	//	NotifyChannel: "web",
	//})
	//TODO
	fmt.Printf("通知：%s执行%s", taskModel.Name, msg)
	return nil
}

// OnlineService里面GetAll
// 这里为避免循环导入在这里重新实现了一遍
func (n *NoticeService) getOnlineUsers() ([]*entity.OnlineUserDto, error) {
	res, err := n.DB.Scan(config.Conf.JWT.OnlineKey)
	if err != nil {
		return nil, err
	}
	users := make([]*entity.OnlineUserDto, 0)
	for _, item := range res {
		var u *entity.OnlineUserDto
		if err := json.Unmarshal([]byte(item.V), u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// 处理服务器重启，重新处理未发布的定时通知
// 可以定义一个定时任务来定时调用这个函数（推荐），或者在系统启动时调用这个函数
func (n *NoticeService) ScheduleNotice(ctx context.Context) {
	result := make([]entity.Notice, 0)
	err := n.NoticeSendRepo.Where(ctx, "scheduled_time IS NOT NULL and send_status=?", 0).
		Where("scheduled_time < NOW()").Find(&result).Error
	if err != nil {
		log.L.Error(err)
	}

	for _, item := range result {
		_, exist := timers[item.ID] //定时器里面是否存在
		if !exist {
			_timer := time.AfterFunc(item.ScheduledTime.Sub(time.Now()), func() {
				err := n.releaseNotice(ctx, &item)
				log.L.Error(err)
			})
			timers[item.ID] = _timer //记录_timer
		}
	}
}
