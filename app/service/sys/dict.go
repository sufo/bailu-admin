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
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"github.com/sufo/bailu-admin/app/domain/resp"
	base2 "github.com/sufo/bailu-admin/app/service/base"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/utils/dict"
)

var DictSet = wire.NewSet(wire.Struct(new(DictOption), "*"), NewDictService)

type DictOption struct {
	DictRepo     *repo.DictRepo
	DictItemRepo *repo.DictItemRepo
	Trans        *repo.Trans
	DictUtil     *dict.DictUtil
}

type DictService struct {
	base2.BaseService[entity.Dict]
	DictOption
}

func NewDictService(opt DictOption) *DictService {
	return &DictService{base2.BaseService[entity.Dict]{opt.DictRepo.Repository}, opt}
}

func (d *DictService) List(ctx context.Context, search string) (*resp.PageResult[entity.Dict], error) {
	builder := base.NewQueryBuilder()
	if search != "" {
		likeSearch := fmt.Sprint("%", search, "%")
		builder.WithWhere("name like ? or code like ?", likeSearch, likeSearch)
	}
	builder.WithPagination(ctx).WithOrder("id")
	return d.DictRepo.ListByBuilder(ctx, builder)
}

// 程序启动调用，加载字典数据到内存
func (d *DictService) LoadingDictCache(ctx context.Context) error {
	items := make([]entity.DictItem, 0)
	if err := d.DictItemRepo.Where(ctx, "status=?", 0).Find(&items).Error; err != nil {
		return err
	}
	//按dictCode分组存放
	if len(items) > 0 {
		var group = make(map[string][]entity.DictItem)
		for _, item := range items {
			code := item.Code
			v := make([]entity.DictItem, 0)
			exists := false
			if v, exists = group[code]; !exists {
				v = make([]entity.DictItem, 0)
			}
			group[item.Code] = v
		}

		for k, v := range group {
			d.DictUtil.SetDictCache(k, v, 0)
		}
	}
	return nil
}

func (d *DictService) CreateDict(ctx context.Context, dict *entity.Dict) error {
	err := d.DictRepo.Create(ctx, dict)
	fmt.Printf("%v", dict)
	return err
}

func (d *DictService) UpdateDict(ctx context.Context, dict *entity.Dict) (*entity.Dict, error) {
	err := d.Trans.Exec(ctx, func(ctx context.Context) error {
		code := dict.Code
		oldDict, err := d.DictRepo.FindById(ctx, dict.ID)
		if oldDict.Code != code { //dict_code发生更改才更新所有DictItem
			err = d.DictItemRepo.Where(ctx, "code=?", oldDict.Code).
				UpdateColumn("code", dict.Code).Error
			if err != nil {
				return err
			}
		}
		err = d.DictRepo.Update(ctx, dict)
		//处理缓存
		if err == nil && config.Conf.Store.EnableCache {
			if dictItems, err2 := d.DictItemRepo.FindByCode(ctx, code); err2 == nil {
				_ = d.DictUtil.SetDictCache(code, dictItems, 0)
			}
		}
		return err
	})
	return dict, err
}

func (d *DictService) Delete(ctx context.Context, dictCodes []string) error {
	count, err := d.DictItemRepo.CountByDictCodes(ctx, dictCodes)
	if err != nil {
		return err
	}
	//说明该字典下面存在字典项，这里采取的策略是不可直接删除
	if count > 0 {
		panic(respErr.WrapLogicResp("1", "%v 已分配，无法删除！", dictCodes))
	}
	//硬删除
	err = d.DictRepo.Where(ctx, "code in ?", dictCodes).Unscoped().Delete(&entity.Dict{}).Error
	//缓存
	if err == nil && config.Conf.Store.EnableCache {
		for _, code := range dictCodes {
			_ = d.DictUtil.RemoveDictCache(code)
		}
	}
	return err
}
