/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 我的通知消息
 */

package message

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/domain/vo"
	"github.com/sufo/bailu-admin/utils"
)

var NoticeSrvSet = wire.NewSet(wire.Struct(new(NoticeService), "*"))

type NoticeService struct {
	NoticeRepo     *repo.NoticeRepo
	NoticeSendRepo *repo.NoticeSendRepo
}

// 查询通知/公告
func (m *MessageService) NoticeList(ctx context.Context, params dto.MessageParams) (*resp.PageResult[vo.Message], error) {
	builder := base.NewQueryBuilder().WithPagination(ctx).
		WithOrder("read_flag asc", "created_at desc")
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if params.Title != "" {
		builder.WithWhere("title like ?", fmt.Sprint("%", params.Title, "%"))
	}
	var pageResult = &resp.PageResult[vo.Message]{
		PageIndex: builder.Paginate.PageIndex,
		PageSize:  builder.Paginate.Limit,
		List:      make([]*vo.Message, 0),
	}
	builder.WithSelect("msg_notice_send.read_flag", "msg_notice.*", "case type WHEN 1 then 'ant-design:notification-outlined' ELSE 'icon-park-outline:announcement' end as icon").
		WithWhere("receive_id=?", user.ID).WithWhere("read_flag != ?", 2).
		WithJoin("left join msg_notice_send on msg_notice_send.msgId=msg_notice.id")

	if result, err := m.NoticeRepo.ListByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		//return result.(*resp.PageResult[entity.Notice]), err
		for _, n := range result.List {
			//item := (*entity.Notice)(unsafe.Pointer(&n))
			msg := vo.Message{
				ID:      utils.Strval(n.ID),
				Icon:    n.Icon,
				Title:   n.Title,
				Content: n.Content,
				IsRead:  n.ReadFlag == 1,
			}
			pageResult.List = append(pageResult.List, &msg)
		}
		pageResult.ItemCount = result.ItemCount
		pageResult.PageCount = result.PageCount
	}
	return pageResult, nil
}
