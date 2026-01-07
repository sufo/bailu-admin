/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package sys

import (
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/resp"
	base2 "github.com/sufo/bailu-admin/app/service/base"
	"context"
	"github.com/google/wire"
)

var SysConfigSet = wire.NewSet(wire.Struct(new(SysConfigOption), "*"), NewSysConfigService)

type SysConfigOption struct {
	SysConfigRepo *repo.SysConfigRepo
}

type SysConfigService struct {
	base2.BaseService[entity.SysConfig]
	SysConfigOption
}

func NewSysConfigService(opt SysConfigOption) *SysConfigService {
	return &SysConfigService{base2.BaseService[entity.SysConfig]{opt.SysConfigRepo.Repository}, opt}
}

func (c *SysConfigService) Create(ctx context.Context, config *entity.SysConfig) error {
	err := c.Repo.Create(ctx, config)
	if err != nil {
		return err
	}
	return nil
}

// 查询所有符合条件的
func (c *SysConfigService) List(ctx context.Context, Task dto.ConfigParams) (*resp.PageResult[entity.SysConfig], error) {
	builder := base.NewQueryBuilder()
	builder.WithWhereStruct(Task).
		WithPagination(ctx)
	if result, err := c.SysConfigRepo.FindByBuilder(ctx, builder); err != nil {
		return nil, err
	} else {
		return result.(*resp.PageResult[entity.SysConfig]), err
	}
}

func (c *SysConfigService) Update(ctx context.Context, config *entity.SysConfig) (*entity.SysConfig, error) {
	err := c.SysConfigRepo.Update(ctx, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *SysConfigService) Delete(ctx context.Context, ids []uint64) error {
	err := c.SysConfigRepo.Delete(ctx, ids)
	if err != nil {
		return err
	}
	return nil
}

func (c *SysConfigService) Status(ctx context.Context, id uint64, status int) error {
	err := c.SysConfigRepo.Where(ctx, "id=?", id).UpdateColumn("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
