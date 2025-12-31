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
	"mime/multipart"
	"bailu/app/domain/entity"
	"bailu/app/domain/repo"
	repoBase "bailu/app/domain/repo/base"
	"bailu/app/domain/resp"
	"bailu/app/service/base"
	"bailu/utils/upload"
)

var FileSet = wire.NewSet(wire.Struct(new(FileOption), "*"), NewFileService)

type FileOption struct {
	FileRepo     *repo.FileRepo
	CategoryRepo *repo.FileCategoryRepo
	Oss          upload.OSS
}

type FileService struct {
	base.BaseService[entity.FileInfo]
	FileOption
}

func NewFileService(opt FileOption) *FileService {
	return &FileService{base.BaseService[entity.FileInfo]{opt.FileRepo.Repository}, opt}
}

// cid->classifyId
func (f *FileService) List(ctx context.Context, cid *uint64, tags string) (*resp.PageResult[entity.FileInfo], error) {
	builder := repoBase.NewQueryBuilder()
	if cid != nil {
		builder.WithWhere("category_id=?", *cid)
	}
	if tags != "" {
		builder.WithWhere("tag like ?", fmt.Sprint("%", tags, "%"))
	}
	builder.WithPagination(ctx).WithOrder("created_at desc")
	return f.FileRepo.ListByBuilder(ctx, builder)
}

func (f *FileService) Create(ctx context.Context, files []*multipart.FileHeader, cid uint64, tag string) error {

	var fileInfos = make([]entity.FileInfo, len(files))
	// 遍历所有图片
	for _, file := range files {
		url, name, err := f.Oss.UploadFile(file)
		if err != nil {
			//报错则删除本次保存的文件
			for _, info := range fileInfos {
				f.Oss.DeleteFile(info.Name)
			}
			return err
		}
		var fileInfo = entity.FileInfo{
			CategoryId: cid,
			OriginName: file.Filename,
			MIME:       file.Header.Get("Content-Type"),
			Size:       file.Size,
			Path:       url,
			Name:       name,
			Url:        url,
		}
		fileInfos = append(fileInfos, fileInfo)
	}

	var err error
	if len(fileInfos) > 0 {
		err = f.FileRepo.CreateInBatch(ctx, fileInfos)
	}
	if err != nil {
		//报错则删除本次所有的文件
		for _, info := range fileInfos {
			f.Oss.DeleteFile(info.Name)
		}
	}
	return err
}

func (f *FileService) UpdateFileInfo(ctx context.Context, id uint64, cid uint64, tags string) error {
	return f.FileRepo.UpdateColumns(ctx, id, map[string]any{"classify_id": cid, "tags": tags}).Error
}

// 批量删除
func (f *FileService) Delete(ctx context.Context, ids []uint64) error {
	files, err := f.FileRepo.FindByIds(ctx, ids)
	if err != nil {
		return err
	}
	//删除数据
	err = f.FileRepo.Delete(ctx, ids)
	if err == nil {
		//删除文件
		for _, file := range files {
			f.Oss.DeleteFile(file.Name)
		}
	}
	return err
}

// classify
func (f *FileService) FileClassifies() (classifies []entity.FileCategory, err error) {
	err = f.CategoryRepo.DB.Find(&classifies).Error
	return
}

func (f *FileService) UpdateCategory(ctx context.Context, id uint64, name string) error {
	return f.CategoryRepo.UpdateColumn(ctx, id, "name=?", name).Error
}

func (f *FileService) CreateCategory(ctx context.Context, name string) (entity.FileCategory, error) {
	classify := entity.FileCategory{Name: name}
	err := f.CategoryRepo.Create(ctx, classify)
	return classify, err
}
func (f *FileService) DeleteCategory(id uint64) error {
	classify := &entity.FileCategory{ID: id}
	return f.CategoryRepo.DB.Delete(classify).Error
}
