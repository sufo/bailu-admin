/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 文件管理
 */

package content

import (
	"bailu/app/api/admin"
	"bailu/app/domain/entity"
	"bailu/app/domain/resp"
	"bailu/app/service/sys"
	respErr "bailu/pkg/exception"
	"bailu/pkg/i18n"
	"bailu/pkg/log"
	"bailu/utils"
	"bailu/utils/page"
	"github.com/gin-gonic/gin"
	"strings"
)

type FileApi struct {
	FileService *sys.FileService
}

func NewFileApi(FileService *sys.FileService) *FileApi {
	return &FileApi{FileService}
}

// @title Index 文件列表
// @Summary 文件列表接口
// @Description 可按文件分类和标签查询文件列表接口
// @Tags File
// @Accept json
// @Produce json
// @Param cid query string false "文件分类ID"
// @Param tag query string false "文件tag"
// @Param pageIndex query int true "页码"
// @Param PageSize query int true "每页条数"
// @response default {object} resp.Response[resp.PageResult[entity.FileInfo]]
// @Success 200 {object} resp.Response[any]
// @Router /api/files [get]
// @Security Bearer
func (f *FileApi) Index(c *gin.Context) {
	_cid := c.DefaultQuery("cid", "")
	tag := c.DefaultQuery("tag", "")
	page.StartPage(c)

	var cid *uint64
	if _cid != "" {
		id, err := utils.ToUint[uint64](_cid)
		if err != nil {
			resp.FailWithError(c, err)
			return
		}
		cid = &id
	}
	result, err := f.FileService.List(c.Request.Context(), cid, tag)
	if err != nil {
		panic(respErr.InternalServerErrorWithMsg(err.Error()))
	}
	//处理url
	for _, f := range result.List {
		if !strings.HasPrefix(f.Url, "http") {
			path := f.Url
			if path != "" && !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			f.Url = admin.ReqSchema(c.Request) + "://" + c.Request.Host + path
		}
	}
	resp.OKWithData(c, result)
}

// @title 新增文件
// @Summary 批量上传文件
// @Description 文件管理-批量上传文件
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "文件"
// @Param tag formData string false "文件tag"
// @Param cid formData int false "文件分类ID"
// @Success 200 {object} resp.Response[any]
// @Router /api/file [post]
// @Security Bearer
func (f *FileApi) Create(c *gin.Context) {
	_cid := c.PostForm("cid")
	tag := c.PostForm("tag")

	//分类字段校验
	var cid uint64
	if _cid != "" {
		id, err := utils.ToUint[uint64](_cid)
		if err == nil {
			resp.FailWithError(c, err)
			return
		}
		cid = id
	}

	form, err := c.MultipartForm()
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	// 获取所有文件
	files := form.File["files"]
	if len(files) == 0 {
		resp.FailWithMsg(c, "No files uploaded")
		return
	}
	//格式检验

	err = f.FileService.Create(c.Request.Context(), files, cid, tag)

	if err != nil {
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title Destroy 批量删除文件接口
// @Summary 批量删除文件接口
// @Description 批量删除文件接口
// @Tags File
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param ids path []string true "字典编码集合"
// @Success 200 {object} resp.Response[any]
// @Router /api/file/{ids} [delete]
// @Security Bearer
func (f *FileApi) Destroy(c *gin.Context) {
	ids := admin.ParseParamArray[uint64](c, "ids")
	if err := f.FileService.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}

// @title 文件分类查询
// @Summary 文件分类查询
// @Description 可按分类名称模糊查询
// @Tags File
// @Accept json
// @Produce json
// @Param name query string false "文件分类名称"
// @response default {object} resp.Response[any]{data=array{value=string,label=string}}
// @Success 200 {object} resp.Response[any]{data=array{value=string,label=string}}
// @Router /api/file/category [get]
// @Security Bearer
func (f *FileApi) Category(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	var query string
	var args any = nil
	if name != "" {
		query = "name like ?"
		args = "%" + name + "%"
	}
	category, err := f.FileService.CategoryRepo.FindBy(c.Request.Context(), query, args)
	if err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.OKWithData(c, category)
}

// @title 新增文件分类
// @Summary 新增文件分类
// @Description 文件管理-新增文件分类
// @Tags File
// @Accept json
// @Produce json
// @Param body body entity.FileCategory true "文件分类信息"
// @Success 200 {object} resp.Response[any]{data=array{label=string,value=string}}
// @Router /api/file/category/add [post]
// @Security Bearer
func (f *FileApi) CategoryCreate(c *gin.Context) {
	var category entity.FileCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		panic(respErr.BadRequestError)
	}
	ctx := c.Request.Context()
	//检查分类名称
	if !f.FileService.CheckUnique(ctx, "name=?", category.Name) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", category.Name))
		return
	}

	if err := f.FileService.CategoryRepo.Create(c.Request.Context(), &category); err != nil {
		resp.FailWithError(c, err)
	}

	result := map[string]string{"label": category.Name, "value": string(category.ID)}
	resp.OKWithData(c, &result)
}

// @title 修改文件分类
// @Summary 修改文件分类接口
// @Description 修改文件分类接口
// @Tags File
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.FileCategory true "文件分类信息"
// @Success 200 {object} resp.Response[entity.FileCategory]
// @Router /api/file/category [put]
func (f *FileApi) CategoryEdit(c *gin.Context) {
	var category entity.FileCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		panic(respErr.BadRequestError)
	}

	ctx := c.Request.Context()
	//检查分类名称
	if !f.FileService.CheckUnique(ctx, "name=? and id !=?", category.Name, category.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", category.Name))
		return
	}

	if err := f.FileService.CategoryRepo.Update(ctx, &category); err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.OKWithData(c, category)
}

// @title 创建或修改文件分类
// @Summary 创建或修改文件分类接口
// @Description 创建或修改文件分类接口
// @Tags File
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body entity.FileCategory true "文件分类信息"
// @Success 200 {object} resp.Response[entity.FileCategory]
// @Router /api/file/category [post]
func (f *FileApi) CategorySave(c *gin.Context) {
	var category entity.FileCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		panic(respErr.BadRequestError)
	}

	ctx := c.Request.Context()
	//检查分类名称
	if !f.FileService.CheckUnique(ctx, "name=? and id !=?", category.Name, category.ID) {
		resp.FailWithMsg(c, i18n.DefTr("admin.existed", category.Name))
		return
	}

	if err := f.FileService.CategoryRepo.Save(ctx, &category); err != nil {
		resp.FailWithError(c, err)
		return
	}
	resp.OKWithData(c, category)
}

// @title Destroy 批量删除文件分类
// @Summary 批量删除文件分类接口
// @Description 批量删除文件分类接口
// @Tags File
// @Accept json
// @Produce json
// @Param ids path []string true "字典编码集合"
// @Success 200 {object} resp.Response[any]
// @Router /api/file/category/{ids} [delete]
// @Security Bearer
func (f *FileApi) CategoryDestroy(c *gin.Context) {
	ids := admin.ParseParamArray[uint64](c, "ids")
	if err := f.FileService.CategoryRepo.Delete(c.Request.Context(), ids); err != nil {
		log.L.Error(err)
		resp.FailWithError(c, err)
	} else {
		resp.Ok(c)
	}
}
