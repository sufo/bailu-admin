/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
	"context"
	"gorm.io/gorm"
)

func NewDictRepo(db *gorm.DB) *DictRepo {
	r := base.Repository[entity.Dict]{db}
	return &DictRepo{r}
}

type DictRepo struct {
	base.Repository[entity.Dict]
}

// Dict item
func NewDictItemRepo(db *gorm.DB) *DictItemRepo {
	r := base.Repository[entity.DictItem]{db}
	return &DictItemRepo{r}
}

type DictItemRepo struct {
	base.Repository[entity.DictItem]
}

func (d *DictItemRepo) CountByDictCodes(ctx context.Context, codes []string) (count int64, err error) {
	err = d.Where(ctx, "code in ?", codes).Count(&count).Error
	return
}

func (d *DictItemRepo) FindByCode(ctx context.Context, code string) ([]entity.DictItem, error) {
	dicts := make([]entity.DictItem, 0)
	err := d.Where(ctx, "code = ?", code).Find(&dicts).Error
	return dicts, err
}
