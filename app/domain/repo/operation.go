/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"gorm.io/gorm"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/base"
)

func NewOperationRecoderRepo(db *gorm.DB) *OperationRecoderRepo {
	r := base.Repository[entity.OperationRecord]{db}
	return &OperationRecoderRepo{r}
}

type OperationRecoderRepo struct {
	base.Repository[entity.OperationRecord]
}
