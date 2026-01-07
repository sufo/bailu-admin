/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package base

import (
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/utils"
	"github.com/sufo/bailu-admin/utils/page"
	"github.com/sufo/bailu-admin/utils/types"
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

type Paginate struct {
	Limit     int
	Offset    int
	PageIndex int //记录pageIndex
}

type Preload struct {
	Query string
	Args  []any
}

type QueryBuilder struct {
	Table string //名称优先于Model名称
	//Model     *IQueryModel //接收结果实体,
	Wheres    map[string]interface{}
	Omits     []string
	Order     []string
	Paginate  *Paginate
	Distincts []interface{}
	Selects   []string
	Joins     []string
	DataScope string
	Preloads  []Preload //预加载
	Group     string    //多个用","号隔开
}

//type IQueryModel interface {
//	//table() string //table name
//}

// Deprecated
func IsZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, float64, int32, int16, int64, float32:
		if v == 0 {
			return true
		}
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	case types.JSONTime:
		return v.IsZero()
	default:
		return false
	}
	return false
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{Wheres: make(map[string]interface{})}
}

func (q *QueryBuilder) WithTable(table string) *QueryBuilder {
	q.Table = table
	return q
}

//func (q *QueryBuilder) WithModel(ModelPtr *IQueryModel) *QueryBuilder {
//	q.Model = ModelPtr
//	return q
//}

func (q *QueryBuilder) getTableAlias() string {
	if q.Table != "" {
		arr := strings.Split(q.Table, " ")
		return arr[len(arr)-1]
	} else {
		return q.Table
	}
}

// sql build where
//
//	cond, vals, err := whereBuild(map[string]interface{}{
//		"name": "jinzhu",
//		"age in": []int{20, 19, 18},
//	})
//	if err != nil {
//		t.Fatal(err)
//	}
//	db.Where(cond, vals...).Find(&users)
//
// where 查询条件
// tableAlias 表名或表别名（解决连表查询同名字段问题，没有连接查询的话可传""）
// Deprecated
func (q *QueryBuilder) WithWhereMapAndAlias(where map[string]interface{}, tableAlias string) *QueryBuilder {
	var whereSQL string
	if query, ok := q.Wheres["query"]; ok {
		whereSQL = query.(string)
	}
	var vals []interface{}
	if args, ok := q.Wheres["args"]; ok {
		vals = args.([]interface{})
	}

	var pre = tableAlias
	if pre != "" {
		pre = tableAlias + "."
	}

	var paginaton = make(map[string]*int, 2)
	for k, v := range where {
		//说明条件没值
		if IsZero(v) {
			continue
		}

		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			log.L.Errorf("Error in query condition: %s. ", k)
			break
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}

		switch len(ks) {
		case 1:
			//处理分页 (这里为了处理条件里面含有分页参数)
			if k == "pageIndex" || k == "pageSize" {
				switch v := v.(type) {
				case *int:
					paginaton[k] = v
				case int: //应该只允许指针的，值类型就无法判断有没有传值
					paginaton[k] = &v
				}
			} else {

				switch v := v.(type) {
				case NullType:
					if v == IsNotNull {
						whereSQL += fmt.Sprint(pre, k, " IS NOT NULL")
					} else {
						whereSQL += fmt.Sprint(pre, k, " IS NULL")
					}
				default:
					whereSQL += fmt.Sprint(pre, k, "=?")
					vals = append(vals, v)
				}
			}
		case 2:
			//k = ks[0]
			k = fmt.Sprint(pre, ks[0])
			switch strings.ToLower(ks[1]) {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
			case "in":
				//whereSQL += fmt.Sprint(k, " in (?)")
				whereSQL += fmt.Sprint(k, " in ?")
				vals = append(vals, v)
			case "like":
				whereSQL += fmt.Sprint(k, " like %?%")
				vals = append(vals, v)
			case "*like":
				whereSQL += fmt.Sprint(k, " like %?")
				vals = append(vals, v)
			case "like*":
				whereSQL += fmt.Sprint(k, " like ?%")
				vals = append(vals, v)
			case "between":
				whereSQL += fmt.Sprint(k, " between ? and ?")
				if vArr, ok := v.([]string); ok {
					vals = append(vals, vArr)
				} else {
					panic(respErr.BadRequestErrorWithMsg(k + "must be an string array"))
				}
			}
		}
	}

	if whereSQL != "" {
		q.Wheres["query"] = whereSQL
		q.Wheres["args"] = vals
	}
	//处理条件里面的分页参数
	//q.WithPaginate(paginaton["pageIndex"], paginaton["pageSize"])
	return q
}

// where map
// Deprecated
func (q *QueryBuilder) WithWhereMap(where map[string]interface{}) *QueryBuilder {
	return q.WithWhereMapAndAlias(where, q.getTableAlias())
}

func (q *QueryBuilder) WithWhereStruct(where any) *QueryBuilder {
	return q.WithWhereStructAndAlias(where, q.getTableAlias())
}

// where struct
// tableAlias (table name or table alise)
func (q *QueryBuilder) WithWhereStructAndAlias(where any, tableAlias string) *QueryBuilder {
	//判断是否strcut
	rv := reflect.ValueOf(where)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct { // Non-structural return error
		log.L.Errorf("ToMap only accepts struct or struct pointer; got %T", rv)
		panic(respErr.BadRequestErrorWithError)
	}

	//带上前面的条件
	var whereSQL string
	if query, ok := q.Wheres["query"]; ok {
		whereSQL = query.(string)
	}
	var vals = make([]interface{}, 0)
	if args, ok := q.Wheres["args"]; ok {
		vals = args.([]interface{})
	}

	//字段前缀（表名或别名+.）
	var pre = tableAlias
	if pre != "" {
		pre = tableAlias + "."
	}

	//生成where
	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {

		fi := t.Field(i)
		queryTag := fi.Tag.Get("query")
		if queryTag == "-" || queryTag == "" {
			continue
		}
		cons := strings.Split(queryTag, ",")
		var qName = cons[0]
		v := rv.Field(i).Interface()

		if !rv.Field(i).IsZero() {
			//拼接AND，这里暂时支持AND
			if whereSQL != "" {
				whereSQL += " AND "
			}

			if len(cons) == 1 {
				switch vt := v.(type) {
				case NullType:
					if vt == IsNotNull {
						whereSQL += fmt.Sprint(pre, qName, " IS NOT NULL")
					} else {
						whereSQL += fmt.Sprint(pre, qName, " IS NULL")
					}
				default:
					whereSQL += fmt.Sprint(pre, qName, "=?")
					vals = append(vals, v)
				}
			} else {
				//比较运算符基本只会是单个字符，但是between会出现两个字符用空格隔开，所以这里要处理下
				//example =、>、< 、 between xxx、......
				c := strings.Split(cons[1], " ")
				qName = fmt.Sprint(pre, qName) //加上表名或表别名
				switch strings.ToLower(c[0]) {
				case "=": //为加快速度
					whereSQL += fmt.Sprint(qName, "=?")
					vals = append(vals, v)
				//case "like":
				//	whereSQL += fmt.Sprint(qName, " like %?%")
				//	vals = append(vals, v)
				//case "*like":
				//	whereSQL += fmt.Sprint(qName, " like %?")
				//	vals = append(vals, v)
				//case "like*":
				//	whereSQL += fmt.Sprint(qName, " like ?%")
				//	vals = append(vals, v)
				case "like":
					whereSQL += fmt.Sprint(qName, " like ?")
					vals = append(vals, fmt.Sprint("%", v, "%"))
				case "*like":
					whereSQL += fmt.Sprint(qName, " like ?")
					vals = append(vals, fmt.Sprint("%", v))
				case "like*":
					whereSQL += fmt.Sprint(qName, " like ?")
					vals = append(vals, fmt.Sprint(v, "%"))
				case "in":
					whereSQL += fmt.Sprint(qName, " in ?")
					vals = append(vals, v)
				case "between": // example: beginDate,between endDate
					var endVal = v
					if len(c) == 0 {
						log.L.Errorf("between not valid for 'query' tag,example 'between endDate'")
					} else {
						endVal = rv.FieldByName(c[1]).String()
						if endVal == "" {
							log.L.Errorf("'%s' not found", c[1])
							panic(respErr.BadRequestError)
						}
					}
					whereSQL += fmt.Sprint(qName, " between ? and ?")
					vals = append(vals, v, endVal)
				default:
					whereSQL += fmt.Sprint(qName, fmt.Sprintf(" %s ?", c[0]))
					vals = append(vals, v)
				}
			}
		}
	}
	if whereSQL != "" {
		q.Wheres["query"] = whereSQL
		q.Wheres["args"] = vals
	}
	return q
}

// query
func (q *QueryBuilder) WithWhere(query string, args ...interface{}) *QueryBuilder {
	var whereSQL string
	if query, ok := q.Wheres["query"]; ok {
		whereSQL = query.(string)
	}
	var vals = make([]interface{}, 0)
	if args, ok := q.Wheres["args"]; ok {
		vals = args.([]interface{})
	}
	if whereSQL != "" {
		whereSQL += " AND "
	}
	whereSQL += query
	//vals = append(vals, args)
	vals = append(vals, args...)
	if whereSQL != "" {
		q.Wheres["query"] = whereSQL
		q.Wheres["args"] = vals
	}
	return q
}

// where map[string]any or struct
// tableAlias 表名或别名
func (q *QueryBuilder) WithWhereAndAlias(where any, tableAlias string) *QueryBuilder {
	if where != nil {
		reflectObj := reflect.TypeOf(where)
		if reflectObj.Elem().Kind() == reflect.Struct {
			q.WithWhereStructAndAlias(where, tableAlias)
		} else if reflectObj.Kind() == reflect.Map {
			q.WithWhereMapAndAlias(where.(map[string]any), tableAlias)
		} else {
			panic(respErr.BadRequestErrorWithMsg(""))
		}
	}
	return q
}

func (q *QueryBuilder) WithOmit(columns ...string) *QueryBuilder {
	q.Omits = append(q.Omits, columns...)
	return q
}

func (q *QueryBuilder) WithOrder(order ...string) *QueryBuilder {
	q.Order = append(q.Order, order...)
	return q
}

// 分页策略（分也与不分页只用一个api的情况）
// 1. 有pageNum，有pageSize，正常执行。
// 2. 有pageNum，有pageSize，且pageSize=0，正常执行。
// 3. 无pageNum，有pageSize，且pageSize=0，正常执行。
// 4. 有pageNum，无pageSize，pageSize取默认值，正常执行。
// 5. 无pageNum，有pageSize，且pageSize不是0，报错：pageNum不能空。
// 7. 无pageNum，无pageSize，报错：pageNum不能空。
// **若要返回全部记录，则必须传pageSize=0。**
func (q *QueryBuilder) WithPaginateParams(pageIndex *int, pageSize *int) *QueryBuilder {

	//分页策略 (这里没有去防范恶意攻击的情况)
	var pIndex = 1
	var pSize = page.DEFAULT_SIZE
	if pageSize == nil {
		if pageIndex == nil {
			panic(respErr.BadRequestErrorWithMsg("pageIndex不能空"))
		} else {
			pIndex = *pageIndex
			if pIndex == 0 {
				panic(respErr.BadRequestErrorWithMsg("pageIndex不能0"))
			}
			pSize = page.DEFAULT_SIZE
		}
	} else if *pageSize == 0 {
		//表示查所有
		return q
	} else {
		pSize = *pageSize
		if pageIndex == nil {
			panic(respErr.BadRequestErrorWithMsg("pageIndex不能空"))
		} else {
			pIndex = *pageIndex
			if pIndex == 0 {
				panic(respErr.BadRequestErrorWithMsg("pageIndex不能0"))
			}
		}
	}

	offset := (pIndex - 1) * pSize
	q.Paginate = &Paginate{
		Limit:     pSize,
		Offset:    offset,
		PageIndex: pIndex,
	}
	return q
}

// 需配合startPage使用
func (q *QueryBuilder) WithPagination(c context.Context) *QueryBuilder {
	if p, exist := appctx.GetPageCtx[page.Pagination](c); exist {
		q.Paginate = &Paginate{
			Limit:     p.Limit,
			Offset:    p.Offset,
			PageIndex: p.PageIndex,
		}
	}
	return q
}

func (q *QueryBuilder) WithDistinct(args ...interface{}) *QueryBuilder {
	q.Distincts = append(q.Distincts, args...)
	return q
}

//	func (q *QueryBuilder) WithSelect(idName string, selects ...string) *QueryBuilder {
//		selects = append(q.Selects, selects...)
//		if len(selects) > 0 {
//			if len(idName) > 0 {
//				selects = append(selects, idName)
//			}
//			// 对Select进行去重
//			selectMap := make(map[string]int, len(selects))
//			for _, e := range selects {
//				if _, ok := selectMap[e]; !ok {
//					selectMap[e] = 1
//				}
//			}
//
//			newSelects := make([]string, 0, len(selects))
//			for k := range selectMap {
//				if len(k) > 0 {
//					newSelects = append(newSelects, k)
//				}
//			}
//			selects = newSelects
//		}
//		q.Selects = selects
//		return q
//	}

// 这里需要注意，如果select给字段取了别名，那后续如果还使用到该字段就要当心
func (q *QueryBuilder) WithSelect(selects ...string) *QueryBuilder {
	selects = append(q.Selects, selects...)
	if len(selects) > 0 {
		// 对Select进行去重
		selectMap := make(map[string]int, len(selects))
		for _, e := range selects {
			if _, ok := selectMap[e]; !ok {
				selectMap[e] = 1
			}
		}

		newSelects := make([]string, 0, len(selects))
		for k := range selectMap {
			if len(k) > 0 {
				newSelects = append(newSelects, k)
			}
		}
		selects = newSelects
	}
	q.Selects = selects
	return q
}

func (q *QueryBuilder) WithJoin(join ...string) *QueryBuilder {
	q.Joins = append(q.Joins, join...)
	return q
}

func (q *QueryBuilder) WithPreload(query string, args ...any) *QueryBuilder {
	q.Preloads = append(q.Preloads, Preload{query, args})
	return q
}

func (q *QueryBuilder) WithDataScope(ctx context.Context, deptAlias string, userAlias string) *QueryBuilder {
	if deptAlias == "" {
		deptAlias = entity.DeptTN
	}
	if userAlias == "" {
		userAlias = entity.UserTN
	}
	user := appctx.GetAuthUser[entity.OnlineUserDto](ctx)
	if user == nil {
		log.L.Warn("current user not found, context must be Request context")
		q.DataScope = fmt.Sprintf("%s.dept_id = 0", deptAlias) //什么也查不到
		return q
	}
	// 如果是超级管理员，则不过滤数据
	if user.IsSuper() {
		return q
	}
	var condition []string
	var sqlBuilder strings.Builder
	for _, role := range user.Roles {
		dataScope := role.DataScope
		if consts.DATA_SCOPE_CUSTOM != dataScope && utils.ContainsInSlice(condition, dataScope) {
			continue
		}
		if consts.DATA_SCOPE_ALL == dataScope {
			sqlBuilder = strings.Builder{}
			break
		} else if consts.DATA_SCOPE_CUSTOM == dataScope {
			sql := fmt.Sprintf(" OR sys_dept.dept_id IN ( SELECT dept_id FROM sys_role_dept WHERE role_id = %d ) ", role.GetID())
			sqlBuilder.WriteString(sql)
		} else if consts.DATA_SCOPE_DEPT == dataScope {
			sql := fmt.Sprintf(" OR %s.dept_id = %d ", deptAlias, user.DeptId)
			sqlBuilder.WriteString(sql)
		} else if consts.DATA_SCOPE_DEPT_AND_CHILD == dataScope {
			sql := fmt.Sprintf(" OR %s.dept_id IN ( SELECT dept_id FROM sys_dept WHERE dept_id = %d or find_in_set( %d , ancestors ) )", deptAlias, user.DeptId, user.DeptId)
			sqlBuilder.WriteString(sql)
		} else if consts.DATA_SCOPE_SELF == dataScope {
			sql := fmt.Sprintf(" OR %s.id = %d ", userAlias, user.ID)
			sqlBuilder.WriteString(sql)
		}
		condition = append(condition, dataScope)
	}
	// 多角色情况下，所有角色都不包含传递过来的权限字符, 则不查询任何数据
	if condition == nil || len(condition) == 0 {
		sqlBuilder.WriteString(fmt.Sprintf("%s.dept_id = 0", deptAlias))
	}

	sqlString := sqlBuilder.String()
	//if sqlString != "" {
	//	sqlString = " AND (" + sqlString[4:] + ")"
	//	return db.Exec(sqlString)
	//}
	if sqlString != "" {
		q.DataScope = sqlString[4:]
	}
	return q
}

func (q *QueryBuilder) WithGroup(group string) *QueryBuilder {
	q.Group = group
	return q
}

func (q *QueryBuilder) ClearWhere(ctx context.Context, deptAlias string, userAlias string) *QueryBuilder {
	q.Wheres = make(map[string]interface{})
	return q
}
