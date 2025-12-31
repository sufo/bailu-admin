/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 操作日志记录
 */

package sys

import (
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	"context"
	"github.com/google/wire"
	"github.com/mssola/user_agent"
)

var OperationSet = wire.NewSet(wire.Struct(new(OperationService), "*"))

type OperationService struct {
	OperRepo *repo.OperationRecoderRepo
}

func (o *OperationService) Create(ctx context.Context, oper *entity.OperationRecord) error {
	return o.OperRepo.Create(ctx, oper)
}

// 查询所有符合条件的，无dataScope限制
func (o *OperationService) List(ctx context.Context, params dto.OperParams) (*resp.PageResult[entity.OperationRecord], error) {
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(params).
		WithPagination(ctx)
	status := params.Status
	if status != "" {
		if status == "0" {
			builder.WithWhere("resp_code = 0")
		} else {
			builder.WithWhere("resp_code <> 0")
		}
	}
	if result, err := o.OperRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {

		//处理os
		pageRecord := result.(*resp.PageResult[entity.OperationRecord])
		for _, record := range pageRecord.List {
			ua := user_agent.New(record.Agent)
			record.OS = ua.OS()
			name, _ := ua.Browser()
			record.Browser = name
		}
		return pageRecord, nil
	}
}

func (o *OperationService) Delete(ctx context.Context, ids []uint64) error {
	return o.OperRepo.Delete(ctx, ids)
}
