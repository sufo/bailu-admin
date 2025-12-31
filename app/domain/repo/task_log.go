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

func NewTaskLogRepo(db *gorm.DB) *TaskLogRepo {
	r := base.Repository[entity.TaskLog]{db}
	return &TaskLogRepo{r}
}

type TaskLogRepo struct {
	base.Repository[entity.TaskLog]
}

func (t *TaskLogRepo) table() *gorm.DB {
	return t.DB.Table(entity.TaskLogTN)
}

// taskLog
func (t *TaskLogRepo) CreateTaskLog(taskLog *entity.TaskLog) (uint64, error) {
	err := t.table().Create(taskLog).Error
	return taskLog.ID, err
}

func (t *TaskLogRepo) UpdateTaskLog(id uint64, vals map[string]any) error {
	return t.table().Where("id=?", id).Updates(vals).Error
}
