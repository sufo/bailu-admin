/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"gorm.io/gorm"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
)

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	r := base.Repository[entity.Task]{db}
	return &TaskRepo{r}
}

type TaskRepo struct {
	base.Repository[entity.Task]
}

// 更新
func (t *TaskRepo) UpdateEntryId(id uint64, column string, value interface{}) *gorm.DB {
	return t.DB.Table(entity.TaskTN).Where("id=?", id).UpdateColumn(column, value)
}

func (t *TaskRepo) UpdateColumn(id uint64, column string, value interface{}) *gorm.DB {
	return t.DB.Table(entity.TaskTN).Where("id=?", id).UpdateColumn(column, value)
}
