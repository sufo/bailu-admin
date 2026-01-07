/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package util

import (
	"errors"
	"gorm.io/gorm"
	"github.com/sufo/bailu-admin/app/domain/resp"
	pagination "github.com/sufo/bailu-admin/utils/page"
)

type PageQuery struct {
	PageSize  int `json:"pageSize"`  //每页条数
	PageIndex int `json:"pageIndex"` //当前页码
}

type QueryBuilder struct {
	query *gorm.DB
	PageQuery
	Sort string
}

type Option func(*QueryBuilder)

func WithPageSize(pageSize int) Option {
	return func(builder *QueryBuilder) {
		builder.PageSize = pageSize
	}
}

func WithPageIndex(PageIndex int) Option {
	return func(builder *QueryBuilder) {
		builder.PageIndex = PageIndex
	}
}

func WithPage(page PageQuery) Option {
	return func(builder *QueryBuilder) {
		builder.PageIndex = page.PageIndex
		builder.PageSize = page.PageSize
	}
}

func WithSort(sort string) Option {
	return func(builder *QueryBuilder) {
		builder.Sort = sort
	}
}

func WithQuery(query *gorm.DB) Option {
	return func(builder *QueryBuilder) {
		builder.query = query
	}
}

func PaginateByOptions[T any](opts ...Option) (*resp.PageResult[T], error) {
	builder := &QueryBuilder{}
	for _, o := range opts {
		o(builder)
	}
	if builder.query == nil {
		return nil, errors.New("WithQuery must be invoke")
	} else {
		return PaginateQuery[T](builder.query, builder.PageQuery, builder.Sort)
	}
}

func Paginate[T any](query *gorm.DB, PageIndex int, pageSize int, sort string) (*resp.PageResult[T], error) {
	if query == nil {
		return nil, errors.New("query must be not nil")
	}
	return PaginateQuery[T](query, PageQuery{
		PageIndex: PageIndex,
		PageSize:  pageSize,
	}, sort)
}

func PaginateQuery[T any](query *gorm.DB, page PageQuery, sort string) (*resp.PageResult[T], error) {
	result := new(resp.PageResult[T])
	if page.PageIndex == 0 {
		page.PageIndex = 1
	}

	if page.PageSize == 0 {
		page.PageSize = pagination.DEFAULT_SIZE
	}
	if sort == "" {
		sort = "id desc"
	}

	result.PageIndex = page.PageIndex
	result.PageSize = page.PageSize

	offset := (page.PageIndex - 1) * page.PageSize
	var rows []*T
	err := query.Limit(page.PageSize).Offset(offset).Order(sort).Find(&rows).Count(&result.ItemCount).Error
	if (int(result.ItemCount) % page.PageSize) > 0 {
		result.PageCount = result.ItemCount/int64(page.PageSize) + 1
	} else {
		result.PageCount = result.ItemCount / int64(page.PageSize)
	}

	result.List = rows
	return result, err
}

// 查询列表
func List[T any](query *gorm.DB) ([]T, error) {
	var rows []T
	err := query.Find(&rows).Error
	return rows, err
}

// 查询单个对象
func GetOne[T any](query *gorm.DB) (T, error) {
	return Take[T](query)
}

func Take[T any](query *gorm.DB) (T, error) {
	var row T
	err := query.Take(&row).Error
	return row, err
}

// 第一条记录
func First[T any](query *gorm.DB) (T, error) {
	var row T
	err := query.First(&row).Error
	return row, err
}

// 最后一条记录
func Last[T any](query *gorm.DB) (T, error) {
	var row T
	err := query.Last(&row).Error
	return row, err
}
