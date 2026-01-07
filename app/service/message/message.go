/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 所有类型消息处理
 */

package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/vo"
	"github.com/sufo/bailu-admin/app/service/sys"
	"github.com/sufo/bailu-admin/utils"
	"time"
)

var MessageSrvSet = wire.NewSet(wire.Struct(new(MessageService), "*"))

type MessageService struct {
	NoticeRepo     *repo.NoticeRepo
	NoticeSendRepo *repo.NoticeSendRepo
	Tran           *repo.Trans
	Dict           *sys.DictItemService
}

const (
	NOTICE = "notice"
	EVENT  = "event"
	CHAT   = "chat"
)

type UnionMessage interface {
	entity.Notice
}

// 查询消息
// // @Param msgType notice event chat
//func (m *MessageService) List(ctx context.Context, msgType string, title string) (*resp.PageResult[vo.Message], error) {
//	builder := base.NewQueryBuilder().WithPagination(ctx).
//		WithOrder("read_flag asc", "created_at desc")
//	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
//	if title != "" {
//		builder.WithWhere("title like ?", fmt.Sprint("%", title, "%"))
//	}
//	var pageResult = &resp.PageResult[vo.Message]{
//		PageIndex: builder.Paginate.PageIndex,
//		PageSize:  builder.Paginate.Limit,
//		List:      make([]*vo.Message, 0),
//	}
//	if NOTICE == msgType {
//		builder.WithSelect("msg_notice_send.read_flag", "msg_notice.*", "case type WHEN 1 then 'ant-design:notification-outlined' ELSE 'icon-park-outline:announcement' end as icon").
//			WithWhere("receive_id=?", user.ID).WithWhere("read_flag != ?", 2).
//			WithJoin("left join msg_notice_send on msg_notice_send.msgId=msg_notice.id")
//
//		if result, err := m.NoticeRepo.ListByBuilder(ctx, builder); err != nil {
//			return nil, err
//		} else {
//			//return result.(*resp.PageResult[entity.Notice]), err
//			for _, n := range result.List {
//				//item := (*entity.Notice)(unsafe.Pointer(&n))
//				msg := vo.Message{
//					ID:      n.ID,
//					Icon:    n.Icon,
//					Title:   n.Title,
//					Content: n.Content,
//					IsRead:  n.ReadFlag == 1,
//				}
//				pageResult.List = append(pageResult.List, &msg)
//			}
//			pageResult.ItemCount = result.ItemCount
//			pageResult.PageCount = result.PageCount
//		}
//
//	} else if EVENT == msgType {
//		//if result, err := m.RemindRepo.
//		//	FindByBuilder(ctx, builder); err != nil {
//		//	return nil, err
//		//} else {
//		//	return result.(*resp.PageResult[entity.EventRemind]), err
//		//}
//		//这个要聚合查询,只查需要通知的, 不分页
//		builder.Paginate = nil
//		builder.WithTable("msg_event_remind e").WithSelect("ers.source_id", "ers.action", "ers.source_id", "ers.source_id", "ers.last_read_time", "ers.userId", "count(id) as total").
//			WithJoin("left join msg_event_remind_status ers on ers.action=e.action and ers.source_id=e.source_id and ers.source_type=e.source_type").
//			WithWhere("e.receive_id=?", user.ID).
//			WithWhere("e.created_at > ers.last_read_time").
//			WithGroup("action,source_id,source_type,receive_id")
//		if result, err := m.RemindRepo.ListAnyByBuilder(ctx, builder); err != nil {
//			return nil, err
//		} else {
//			//TODO这个最好加上字典缓存，然后查询缓存
//			opts, _ := m.Dict.FindOptions(ctx, "")
//			//return result.(*resp.PageResult[entity.EventRemind]), nil
//			for _, item := range result.List {
//				e := (*entity.EventRemindStatus)(unsafe.Pointer(&item))
//				action, exist := lo.Find[*vo.Option[string]](opts, func(item *vo.Option[string]) bool { return item.Value == e.Action })
//				if !exist {
//					return nil, errors.New("action is not found")
//				}
//				title := fmt.Sprintf("有%d%s了你的%s", e.Total, action, e.SourceType)
//
//				msg := vo.Message{
//					ID:      e.ID,
//					Icon:    "mdi:event-clock",
//					Title:   title,
//					Content: title,
//					IsRead:  false,
//					Date:    e.LastReadTime,
//				}
//				pageResult.List = append(pageResult.List, &msg)
//			}
//			pageResult.ItemCount = int64(len(pageResult.List))
//			pageResult.PageCount = -1
//		}
//
//	} else if CHAT == msgType {
//		builder.WithWhere("user1_id=? or user2_id=?", user.ID, user.ID)
//		if result, err := m.ChatSessionRepo.
//			FindByBuilder(ctx, builder); err != nil {
//			return nil, err
//		} else {
//			//return result.(*resp.PageResult[entity.ChatSession]), err
//			data := result.(*resp.PageResult[entity.ChatSession])
//			for _, c := range data.List {
//				lastReadTime := global.Ternary[time.Time](c.User1Id == user.ID, c.LastRead1Time, c.LastRead2Time)
//				msg := vo.Message{
//					ID:      0,
//					Icon:    "",
//					Avatar:  c.Avatar,
//					Title:   "",
//					Content: c.Last_message,
//					IsRead:  lastReadTime.Before(time.Now()),
//					Date:    lastReadTime,
//				}
//				pageResult.List = append(pageResult.List, &msg)
//			}
//			pageResult.ItemCount = data.ItemCount
//			pageResult.PageCount = data.PageCount
//		}
//	} else {
//		return nil, nil
//	}
//	return pageResult, nil
//}

// 查询事件提醒、私聊，这是是聚合查询，不分页
func (m *MessageService) UnreadList(ctx context.Context, params dto.MessageParams) ([]*vo.Message, error) {
	builder := base.NewQueryBuilder()
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	msgs := make([]*vo.Message, 0)

	if NOTICE == params.Type {
		if params.Title != "" {
			builder.WithWhere("title like ?", fmt.Sprint("%", params.Title, "%"))
		}
		//未读
		builder.WithWhere("msg_notice_send.read_flag = ?", 0)
		builder.WithSelect("msg_notice_send.read_flag", "msg_notice.*", "case type WHEN 1 then 'ant-design:notification-outlined' ELSE 'icon-park-outline:announcement' end as icon").
			WithWhere("receive_id=?", user.ID).
			WithJoin("left join msg_notice_send on msg_notice_send.msg_id=msg_notice.id").
			WithOrder("msg_notice_send.created_at desc")
		if list, err := m.NoticeRepo.FindAllByBuilder(ctx, builder); err != nil {
			return nil, err
		} else {
			for _, n := range list {
				//item := (*entity.Notice)(unsafe.Pointer(&n))
				msg := vo.Message{
					ID:      utils.Strval(n.ID),
					Icon:    n.Icon,
					Title:   n.Title,
					Content: n.Content,
					IsRead:  n.ReadFlag == 1,
				}
				msgs = append(msgs, &msg)
			}
		}

	} else if EVENT == params.Type {
		//TODO

	} else if CHAT == params.Type {
		//TODO
	} else {
		return nil, nil
	}
	return msgs, nil
}

// 未读消息数量
// // @Param msgType notice event chat
func (m *MessageService) UnreadCount(ctx context.Context, msgType string) (map[string]int64, error) {

	var nCount int64 = 0
	var eCount int64 = 0
	var cCount int64 = 0
	var err error = nil
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if NOTICE == msgType { //通知/公告 未读

		err = m.NoticeSendRepo.Where(ctx, "receive_id=? and read_flag", user.ID, 0).Count(&nCount).Error

	} else if EVENT == msgType { //提醒未读
		//TODO
	} else if CHAT == msgType { // 私信未读
		//TODO

	} else { //查询所有未读
		err = m.NoticeSendRepo.Where(ctx, "receive_id=? and read_flag", user.ID, 0).Count(&nCount).Error
		if err != nil {
			return nil, err
		}

	}

	result := map[string]int64{
		"notice": nCount,
		"event":  eCount,
		"chat":   cCount,
		"total":  nCount + eCount + cCount,
	}
	return result, nil
}

// 全部设置已读
func (m *MessageService) ReadAll(ctx context.Context, msgType string) error {
	var err error = nil
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if NOTICE == msgType { //通知/公告 未读
		err = m.NoticeSendRepo.Where(ctx, "receive_id=? and read_flag", user.ID, 0).
			Update("read_flag", 1).Error
	} else if EVENT == msgType { //提醒未读
		//TODO
	} else if CHAT == msgType { // 私信未读
		//TODO
	} else { //查询所有未读
		err = m.Tran.Exec(ctx, func(ctx context.Context) error {
			//通知/公告
			if err = m.NoticeSendRepo.Where(ctx, "receive_id=? and read_flag", user.ID, 0).
				Update("read_flag", 1).Error; err != nil {
				return err
			}
			return err
		})

	}
	return err
}

// 单条设置已读(只针对通知/公告)
//
//	func (m *MessageService) ReadNotice(ctx context.Context, id uint64) error {
//		return m.NoticeSendRepo.Where(ctx, "id=? and read_flag", id, 0).
//			UpdateColumns(map[string]any{"read_flag": 1, "read_time": time.Now()}).Error
//	}
//
// 单条/单组设置已读
func (m *MessageService) Read(ctx context.Context, msgType string, id string) error {
	//user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if NOTICE == msgType {
		return m.NoticeSendRepo.Where(ctx, "id=? and read_flag", id, 0).
			UpdateColumns(map[string]any{"read_flag": 1, "read_time": time.Now()}).Error
	} else if EVENT == msgType {
		//TODO
		return nil
	} else if CHAT == msgType {
		//TODO
		return nil
	} else {
		return errors.New("msgType must be 'notice','event' or 'chat'")
	}
}

func (m *MessageService) Delete(ctx context.Context, ids []uint64, msgType string) error {
	var err error = nil
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if NOTICE == msgType {
		err = m.NoticeSendRepo.Where(ctx, "id in ? and receive_id=?", ids, user.ID).
			UpdateColumns(map[string]any{"read_flag": 2, "read_time": time.Now()}).Error
	} else if EVENT == msgType {
		//TODO
	} else if CHAT == msgType { //私信
		return m.Tran.Exec(ctx, func(ctx context.Context) error {

			//TODO
			return nil

		})
	}
	return err
}

// 清空
func (m *MessageService) Clear(ctx context.Context, msgType string) error {
	var err error = nil
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if NOTICE == msgType {
		err = m.NoticeSendRepo.Where(ctx, "receive_id=?", user.ID).
			UpdateColumns(map[string]any{"read_flag": 2, "read_time": time.Now()}).Error
	} else if EVENT == msgType {
		//TODO
	} else if CHAT == msgType { //私信
		return m.Tran.Exec(ctx, func(ctx context.Context) error {
			//TODO
			return nil

		})
	}
	return err
}
