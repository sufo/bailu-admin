/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 使用这个就表示必须会分页了
 */
//Deprecated
package page

import (
	"github.com/sufo/bailu-admin/app/core/appctx"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	PAGINATION_KEY = "gin-pagination"

	DEFAULT_INDEX_KEY = "pageIndex"
	DEFAULT_SIZE_KEY  = "pageSize"
	DEFAULT_INDEX     = 1
	DEFAULT_SIZE      = 10
	DEFAULT_MAX_SIZE  = 100
)

type Pagination struct {
	Limit     int
	Offset    int
	PageIndex int
	PageSize  int
}

func Default(c *gin.Context) {
	New(
		c,
		DEFAULT_INDEX_KEY,
		DEFAULT_SIZE_KEY,
		DEFAULT_INDEX,
		DEFAULT_SIZE,
		DEFAULT_MAX_SIZE,
	)
}

func New(c *gin.Context, indexKey string, sizeKey string, defaultIndex int, defaultSize int, maxSize int) {
	//处理pageIndex
	indexStr := c.DefaultQuery(indexKey, strconv.Itoa(defaultIndex))
	pIndex, err := strconv.Atoi(indexStr)
	if err != nil {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be an integer", indexKey)))
	}
	if pIndex <= 0 {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be positive", indexKey)))
	}

	//处理pageSize
	sizeStr := c.DefaultQuery(sizeKey, strconv.Itoa(defaultSize))
	pSize, err := strconv.Atoi(sizeStr)
	if err != nil {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s number must be an integer", sizeKey)))
	}
	if pSize <= 0 {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be positive", sizeKey)))
	}
	//处理max
	// Validate for min and max page size
	if pSize > maxSize {
		panic(respErr.BadRequestErrorWithMsg(fmt.Sprintf("%s must be between 1 and %d", sizeKey, maxSize)))
	}

	offset := (pIndex - 1) * pSize
	p := &Pagination{
		Limit:     pSize,
		Offset:    offset,
		PageIndex: pIndex,
		PageSize:  pSize,
	}

	//page存储在request context中
	ctx := appctx.NewPageCtx(c.Request.Context(), p)
	c.Request = c.Request.WithContext(ctx)
	//c.Set(PAGINATION_KEY, p)
}

/**
 * 设置请求分页数据
 */
func StartPage(c *gin.Context) {
	Default(c)
}

func ClearPage(c *gin.Context) {
	c.Set(PAGINATION_KEY, nil)
}
