/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package repo

import (
	"bailu/app/domain/entity"
	"bailu/app/domain/repo/base"
	"gorm.io/gorm"
)

func NewFileRepo(db *gorm.DB) *FileRepo {
	r := base.Repository[entity.FileInfo]{db}
	return &FileRepo{r}
}

type FileRepo struct {
	base.Repository[entity.FileInfo]
}

func NewCategoryRepo(db *gorm.DB) *FileCategoryRepo {
	r := base.Repository[entity.FileCategory]{db}
	return &FileCategoryRepo{r}
}

type FileCategoryRepo struct {
	base.Repository[entity.FileCategory]
}
