/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package sys

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"bailu/app/domain/dto"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	"bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	"bailu/app/domain/vo"
	base2 "bailu/app/service/base"
	respErr "bailu/pkg/exception"
	"bailu/utils"
)

var DictItemSet = wire.NewSet(wire.Struct(new(DictItemOption), "*"), NewDictItemService)

type DictItemOption struct {
	DictItemRepo *repo.DictItemRepo
}

type DictItemService struct {
	base2.BaseService[entity.DictItem]
	DictItemOption
}

func NewDictItemService(opt DictItemOption) *DictItemService {
	return &DictItemService{base2.BaseService[entity.DictItem]{opt.DictItemRepo.Repository}, opt}
}

func (d *DictItemService) List(ctx context.Context, code, label, status string) (*resp.PageResult[entity.DictItem], error) {
	builder := base.NewQueryBuilder().WithWhere("code=?", code).WithPagination(ctx)
	if label != "" {
		builder.WithWhere("label like ?", fmt.Sprint("%", label, "%"))
	}
	if status != "" {
		_sta, err := utils.ToUint[uint8](status)
		if err != nil {
			panic(respErr.BadRequestErrorWithMsg("status must be number"))
		}
		builder.WithWhere("label = ?", _sta)
	}
	return d.DictItemRepo.ListByBuilder(ctx, builder)
}

func (d *DictItemService) Create(ctx context.Context, dict *entity.DictItem) error {
	return d.DictItemRepo.Create(ctx, dict)
}

func (d *DictItemService) Update(ctx context.Context, dictItem *entity.DictItem) (*entity.DictItem, error) {
	err := d.DictItemRepo.Update(ctx, dictItem)
	return dictItem, err
}

func (d *DictItemService) Delete(ctx context.Context, ids []uint64) error {
	return d.DictItemRepo.Where(ctx, "id in ?", ids).Unscoped().Delete(&entity.DictItem{}).Error
}

func (d *DictItemService) FindByCode(ctx context.Context, code string) ([]entity.DictItem, error) {
	dicts := make([]entity.DictItem, 0)
	err := d.DictItemRepo.Where(ctx, "code = ?", code).Find(&dicts).Error
	return dicts, err
}

func (p *DictItemService) FindOptions(ctx context.Context, code string) ([]*vo.Option[string], error) {
	builder := base.NewQueryBuilder()
	builder.WithSelect("value").WithSelect("label").
		WithSelect("is_default").
		WithWhere("code=?", code).
		WithWhere("status=?", 1). //启用
		WithOrder("sort asc")
	options := make([]*vo.Option[string], 0)
	err := p.DictItemRepo.FindModelByBuilder(ctx, builder, &options)
	return options, err
}

func (d *DictItemService) ChangeStatus(ctx context.Context, param dto.StatusParam) error {
	return d.DictItemRepo.Where(ctx, "id=?", param.ID).UpdateColumn("status", param.Status).Error
}
