/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package base

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/app/domain/repo/util"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

type Repository[T entity.IModel] struct {
	DB *gorm.DB
}

// 这里是否考虑使用session
// Get gorm.DB from context
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	tx, ok := appctx.FromTrans(ctx)
	if ok && !appctx.FromNoTrans(ctx) {
		if appctx.FromTransLock(ctx) {
			tx = tx.Clauses(clause.Locking{Strength: "UPDATE"})
		}
		return tx
	}
	//return defDB
	return defDB.WithContext(ctx)
}

// Get gorm.DB.Model from context
func GetDBWithModel[T any](ctx context.Context, defDB *gorm.DB) *gorm.DB {
	var t T
	return GetDB(ctx, defDB).Model(&t)
}

func (r *Repository[T]) WithQueryBuilder(ctx context.Context, builer *QueryBuilder) *gorm.DB {
	var db *gorm.DB
	if builer.Table != "" {
		db = r.GetDB(ctx).Table(builer.Table)
	} else {
		db = GetDBWithModel[T](ctx, r.DB)
	}
	//预加载
	if len(builer.Preloads) > 0 {
		for _, preload := range builer.Preloads {
			if len(preload.Args) == 0 {
				db = db.Preload(preload.Query)
			} else {
				db = db.Preload(preload.Query, preload.Args...)
			}
		}
	}

	if len(builer.Selects) > 0 {
		db = db.Select(builer.Selects)
	}
	if len(builer.Distincts) > 0 {
		db = db.Distinct(builer.Distincts)
	}
	if len(builer.Omits) > 0 {
		db = db.Omit(builer.Omits...)
	}
	//join
	if builer.Joins != nil && len(builer.Joins) > 0 {
		for _, join := range builer.Joins {
			db = db.Joins(join)
		}
	}

	//判断map是否为空
	if builer.Wheres != nil && len(builer.Wheres) > 0 {
		db.Where(builer.Wheres["query"], (builer.Wheres["args"].([]interface{}))...)
	}

	if builer.Order != nil {
		db = db.Order(strings.Join(builer.Order, ","))
	}

	if builer.Group != "" {
		db.Group(builer.Group)
	}

	if builer.Paginate != nil {
		db = db.Offset(builer.Paginate.Offset).Limit(builer.Paginate.Limit)
	}

	return db
}

func (r *Repository[T]) FindByBuilder(ctx context.Context, builder *QueryBuilder) (any, error) {
	//增加默认排序
	if builder.Order == nil {
		var alias = ""
		if builder.Table != "" {
			nameArr := strings.Split(builder.Table, " ")
			if len(nameArr) > 0 { //说明存在别名
				alias = nameArr[len(nameArr)-1] //取别名
			} else { //说明只有表名
				alias = nameArr[0] //取表名
			}
			alias = fmt.Sprint(alias, ".")
		}

		//builder.WithOrder(fmt.Sprint(alias, "sort desc"), fmt.Sprint(alias, "updated_at desc"))
		var t T
		v := reflect.TypeOf(t)
		//判断是否存在sort，存在则默认使用sort排序
		if _, exist := v.FieldByName("Sort"); exist {
			builder.WithOrder(fmt.Sprint(alias, "sort asc"))
		}
	}

	if builder.Paginate == nil {
		return r.FindAllByBuilder(ctx, builder)
	} else {
		return r.ListByBuilder(ctx, builder)
	}
}

func (r *Repository[T]) FindAllByBuilder(ctx context.Context, builder *QueryBuilder) (t []*T, err error) {
	err = r.WithQueryBuilder(ctx, builder).Find(&t).Error
	return
}
func (r *Repository[T]) ListByBuilder(ctx context.Context, builder *QueryBuilder) (*resp.PageResult[T], error) {
	result := new(resp.PageResult[T])
	//分页校验
	if builder.Paginate == nil {
		log.L.Warn("page.StartPage() must be invoke in controller")
		return nil, fmt.Errorf("page.StartPage() must be invoke in controller")
	}
	var pageSize = builder.Paginate.Limit

	err := r.WithQueryBuilder(ctx, builder).Find(&result.List).Offset(-1).Count(&result.ItemCount).Error
	if err == nil {
		if (int(result.ItemCount) % pageSize) > 0 {
			result.PageCount = result.ItemCount/int64(pageSize) + 1
		} else {
			result.PageCount = result.ItemCount / int64(pageSize)
		}
		result.PageIndex = builder.Paginate.PageIndex
		result.PageSize = pageSize
	}
	return result, err
}

func (r *Repository[T]) FindRowsByBuilder(ctx context.Context, builder *QueryBuilder) (*sql.Rows, error) {
	return r.WithQueryBuilder(ctx, builder).Rows()
}
func (r *Repository[T]) ListAnyByBuilder(ctx context.Context, builder *QueryBuilder) (*resp.PageResult[any], error) {
	result := new(resp.PageResult[any])
	var pageSize = builder.Paginate.Limit
	err := r.WithQueryBuilder(ctx, builder).Scan(&result.List).Offset(-1).Count(&result.ItemCount).Error
	if err == nil {
		if (int(result.ItemCount) % pageSize) > 0 {
			result.PageCount = result.ItemCount/int64(pageSize) + 1
		} else {
			result.PageCount = result.ItemCount / int64(pageSize)
		}
		result.PageIndex = builder.Paginate.PageIndex
		result.PageSize = pageSize
	}
	return result, err
}

// 获取指定model
func (r *Repository[T]) FindModelByBuilder(ctx context.Context, builder *QueryBuilder, model any) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		log.L.Errorf("%v must be Pointer", model)
		//return fmt.Errorf("%v must be Pointer", model)
		return nil
	}
	if model == nil {
		log.L.Errorf("model mustn't be nil")
		return nil
	}
	return r.WithQueryBuilder(ctx, builder).Scan(model).Error
}

func (r *Repository[T]) FindById(ctx context.Context, id uint64) (t *T, err error) {
	err = GetDBWithModel[T](ctx, r.DB).Where("id=?", id).Find(&t).Error
	return
}
func (r *Repository[T]) FindByIds(ctx context.Context, ids []uint64) (t []*T, err error) {
	err = GetDBWithModel[T](ctx, r.DB).Where("id in ?", ids).Find(&t).Error
	return
}

// 硬删除
func (r *Repository[T]) Delete(ctx context.Context, ids interface{}) error {
	db := GetDBWithModel[T](ctx, r.DB)
	var t T
	return db.Unscoped().Delete(&t, ids).Error
}

// 软删除
func (r *Repository[T]) SoftDel(ctx context.Context, ids interface{}) error {
	var t T
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Delete(&t, ids).Error
}

// 更新
func (r *Repository[T]) Update(ctx context.Context, t *T) error {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Where("id=?", (*t).GetID()).Updates(t).Error
}

// 新增或更新
func (r *Repository[T]) Save(ctx context.Context, t *T) error {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Save(t).Error
}

// 支持批量
func (r *Repository[T]) Create(ctx context.Context, t any) error {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Create(t).Error
}

func (r *Repository[T]) CreateInBatch(ctx context.Context, t []T) error {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.CreateInBatches(t, len(t)).Error
}

// 返回值
func (r *Repository[T]) Insert(ctx context.Context, t *entity.IModel) (entity.IModel, error) {
	db := GetDBWithModel[T](ctx, r.DB)
	err := db.Create(t).Error
	return *t, err
}

func (r *Repository[T]) First(ctx context.Context) (t T, err error) {
	err = GetDBWithModel[T](ctx, r.DB).First(&t).Error
	return
}

func (r *Repository[T]) Last(ctx context.Context) (t T, err error) {
	err = GetDBWithModel[T](ctx, r.DB).Last(&t).Error
	return
}

func (r *Repository[T]) Take(ctx context.Context) (t T, err error) {
	err = GetDBWithModel[T](ctx, r.DB).Take(&t).Error
	return
}

// whether it exists
func (r *Repository[T]) IsExist(ctx context.Context, query interface{}, args ...interface{}) (bool, error) {
	var exists = false
	db := GetDBWithModel[T](ctx, r.DB).
		Select("count(*) > 0")
	if query != nil {
		db.Where(query, args...)
	}
	err := db.Find(&exists).Error
	return exists, err
}

// page
func (r *Repository[T]) Paginate(ctx context.Context, pageIndex, pageSize int) (*resp.PageResult[T], error) {
	db := GetDBWithModel[T](ctx, r.DB)
	return util.Paginate[T](db, pageIndex, pageSize, "")
}

//	func (r *Repository[T]) List(ctx context.Context, pageIndex *int, pageSize *int, sort string) (*resp.PageResult[T], error) {
//		result := new(resp.PageResult[T])
//		var rows []*T
//		db := GetDBWithModel[T](ctx, r.DB)
//		var pIndex = 1
//		var pSize = consts.PAGE_SIZE
//		if pageIndex != nil || pageSize != nil {
//			if pageIndex != nil {
//				pIndex = *pageIndex
//			}
//			if pageSize != nil {
//				pSize = *pageSize
//			}
//
//			//说明传了分页参数，并且都为0,则不分页
//			if pIndex == 0 && pSize == 0 {
//				if sort != "" {
//					db = db.Order(sort)
//				}
//				err := db.Find(&result.List).Count(&result.ItemCount).Error
//				return result, err
//			} else {
//				//不都为0的情况
//				if pIndex == 0 {
//					pIndex = 1
//				}
//				if pSize == 0 {
//					pSize = consts.PAGE_SIZE
//				}
//			}
//		}
//		offset := (pIndex - 1) * pSize
//		db = db.Limit(pSize).Offset(offset)
//		if sort != "" {
//			db = db.Order(sort)
//		}
//		err := db.Order(sort).Find(&rows).Count(&result.ItemCount).Error
//		if (int(result.ItemCount) % pSize) > 0 {
//			result.PageCount = result.ItemCount/int64(pSize) + 1
//		} else {
//			result.PageCount = result.ItemCount / int64(pSize)
//		}
//		return result, err
//	}

//func (r *Repository[T]) List(ctx context.Context, pageIndex *int, pageSize *int, sort string) (*resp.PageResult[T], error) {
//	result := new(resp.PageResult[T])
//	var rows []*T
//	db := GetDBWithModel[T](ctx, r.DB)
//
//	var pIndex = 1
//	var pSize = consts.PAGE_SIZE
//	if pageSize == nil {
//		if pageIndex == nil {
//			panic(respErr.BadRequestErrorWithMsg("pageIndex不能空"))
//		} else {
//			pIndex = *pageIndex
//			if pIndex == 0 {
//				panic(respErr.BadRequestErrorWithMsg("pageIndex不能0"))
//			}
//			pSize = consts.PAGE_SIZE
//		}
//	} else if *pageSize == 0 {
//		//表示查所有
//		if sort != "" {
//			db = db.Order(sort)
//		}
//		err := db.Find(&result.List).Count(&result.ItemCount).Error
//		return result, err
//	} else {
//		pSize = *pageSize
//		if pageIndex == nil {
//			panic(respErr.BadRequestErrorWithMsg("pageIndex不能空"))
//		} else {
//			pIndex = *pageIndex
//			if pIndex == 0 {
//				panic(respErr.BadRequestErrorWithMsg("pageIndex不能0"))
//			}
//		}
//	}
//
//	offset := (pIndex - 1) * pSize
//	db = db.Limit(pSize).Offset(offset)
//	if sort != "" {
//		db = db.Order(sort)
//	}
//	err := db.Order(sort).Find(&rows).Count(&result.ItemCount).Error
//	if (int(result.ItemCount) % pSize) > 0 {
//		result.PageCount = result.ItemCount/int64(pSize) + 1
//	} else {
//		result.PageCount = result.ItemCount / int64(pSize)
//	}
//	return result, err
//}

func (r *Repository[T]) List(ctx context.Context, sort string) (*resp.PageResult[T], error) {
	result := new(resp.PageResult[T])
	var rows []*T
	db := GetDBWithModel[T](ctx, r.DB)

	var pIndex = 1
	var pSize = page.DEFAULT_SIZE
	p, exist := appctx.GetPageCtx[page.Pagination](ctx)
	if exist {
		pIndex = p.PageIndex
		pSize = p.PageSize
	} else {
		log.L.Warn("分页没有调用page.startPage()")
	}
	offset := (pIndex - 1) * pSize
	db = db.Limit(pSize).Offset(offset)
	if sort != "" {
		db = db.Order(sort)
	}
	err := db.Order(sort).Find(&rows).Count(&result.ItemCount).Error
	if (int(result.ItemCount) % pSize) > 0 {
		result.PageCount = result.ItemCount/int64(pSize) + 1
	} else {
		result.PageCount = result.ItemCount / int64(pSize)
	}
	return result, err
}

// 更新某个字段
//
//	func (r *Repository[T]) UpdateColumn(ctx context.Context, id uint64, column string, value interface{}) *gorm.DB {
//		db := GetDBWithModel[T](ctx, r.DB)
//		user := appctx.GetAuthUser(ctx)
//		return db.Where("id=?", id).UpdateColumns(map[string]interface{}{column: value, "update_by": user.UserName})
//	}
func (r *Repository[T]) UpdateColumn(ctx context.Context, id uint64, column string, value interface{}) *gorm.DB {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Where("id=?", id).UpdateColumn(column, value)
}

// 多个字段更新
func (r *Repository[T]) UpdateColumns(ctx context.Context, id uint64, values interface{}) *gorm.DB {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Where("id=?", id).Updates(values)
}

func (r *Repository[T]) Where(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	db := GetDBWithModel[T](ctx, r.DB)
	return db.Where(query, args...)
}

func (r *Repository[T]) Find(ctx context.Context, conds ...any) (t []*T, err error) {
	db := GetDBWithModel[T](ctx, r.DB)
	err = db.Find(&t, conds).Error
	return
}

// 通过条件查询
func (r *Repository[T]) FindBy(ctx context.Context, query interface{}, args ...interface{}) (t []*T, err error) {
	db := GetDBWithModel[T](ctx, r.DB)
	if query != nil && query != "" {
		err = db.Where(query, args).Find(&t).Error
	} else {
		err = db.Find(&t).Error
	}
	return
}
func (r *Repository[T]) FirstBy(ctx context.Context, query interface{}, args ...interface{}) (t *T, err error) {
	db := GetDBWithModel[T](ctx, r.DB)
	err = db.Where(query, args).Find(&t).Error
	return
}

func (r *Repository[T]) GetDB(ctx context.Context) *gorm.DB {
	return GetDB(ctx, r.DB)
}
func (r *Repository[T]) WithModel(ctx context.Context) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB)
}
func (r *Repository[T]) Select(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Select(query, args...)
}
func (r *Repository[T]) Order(ctx context.Context, value interface{}) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Order(value)
}

func (r *Repository[T]) Model(ctx context.Context, value interface{}) *gorm.DB {
	return GetDB(ctx, r.DB).Model(value)
}

func (r *Repository[T]) Or(ctx context.Context, value interface{}) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Or(value)
}

func (r *Repository[T]) Limit(ctx context.Context, limit int) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Limit(limit)
}

func (r *Repository[T]) Distinct(ctx context.Context, args ...interface{}) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Distinct(args...)
}

func (r *Repository[T]) Omit(ctx context.Context, columns ...string) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Omit(columns...)
}

func (r *Repository[T]) Not(ctx context.Context, query interface{}, args ...interface{}) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Not(query, args...)
}

// 注意会绑定T
func (r *Repository[T]) Unscoped(ctx context.Context) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Unscoped()
}

func (r *Repository[T]) Scopes(ctx context.Context, funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return GetDBWithModel[T](ctx, r.DB).Scopes(funcs...)
}
func (r *Repository[T]) Table(ctx context.Context, name string, args ...interface{}) *gorm.DB {
	return GetDB(ctx, r.DB).Table(name, args...)
}

//	cond, vals, err := whereBuild(map[string]interface{}{
//		"name": "jinzhu",
//		"age in": []int{20, 19, 18},
//	})
//
// db.Where(cond, vals...).Find(&users)
func (r *Repository[T]) WithWhere(ctx context.Context, where map[string]interface{}) *gorm.DB {
	var whereSQL string
	var vals []interface{}
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			log.L.Errorf("Error in query condition: %s. ", k)
			break
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return GetDBWithModel[T](ctx, r.DB).Where(whereSQL, vals)
}

func (r *Repository[T]) getTableName() (string, error) {
	stmt := &gorm.Statement{DB: r.DB}
	var t T
	if err := stmt.Parse(&t); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

// 谨慎使用
func (r *Repository[T]) Truncate(ctx context.Context) error {
	tName, err := r.getTableName()
	if err != nil {
		return err
	}
	return GetDB(ctx, r.DB).Raw(fmt.Sprintf("TRUNCATE TABLE %s;", tName)).Error
}
