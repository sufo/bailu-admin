/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package dto

type LoginParams struct {
	Username  string `json:"username" form:"username" binding:"required"`   // 用户名
	Password  string `json:"password" form:"password" binding:"required"`   // 密码(rsa加密)
	CaptchaId string `json:"captchaId" form:"captchaId" binding:"required"` // 验证码ID
	ImgCode   string `json:"imgCode" form:"imgCode" binding:"required"`     // 验证码
}

type LoginUser struct {
	Username string `json:"username" form:"username" binding:"required"` // 用户名
	Password string `json:"password" form:"password" binding:"required"` // 密码(rsa加密)
}

type SMSLoginParams struct {
	Phone   string `form:"phone" binding:"required"`        // 手机好
	SMSCode string `form:"smsCode" binding:"numeric,len=6"` // 短信验证码
}

type UserQueryParams struct {
	//PaginationParam
	Username  string `form:"username" query:"username,like"`
	Status    *int   `form:"status" query:"status"`
	DeptId    string `form:"deptId" query:"-"`
	Phone     string `form:"phone" query:"phone"`
	BeginDate string `form:"beginDate" query:"createAt,between endDate" binding:"omitempty,len=8"` //omitempty表示可选，存在则继续向后校验 YYYYMMDD
	EndDate   string `form:"endDate" binding:"omitempty,len=8"`
}

type ResetParams struct {
	Password string `json:"password" form:"password" binding:"required"` // 密码(rsa加密)
	DialCode string `json:"dialCode" form:"dialCode"`
	Phone    string `form:"phone" binding:"required"`        //e164
	SMSCode  string `form:"smsCode" binding:"numeric,len=6"` // 短信验证码
}

type RegisterParams struct {
	//json对应前端传递字段名
	Username string `json:"username" form:"username" binding:"required"` // 用户名
	ResetParams
}

type SetPwdParams struct {
	Id       uint64 `json:"id" form:"id" binding:"required,numeric,gt=1"` //1为管理员id
	Password string `json:"password" form:"password" binding:"required"`
}

type ChangePwdParams struct {
	Id          uint64 `json:"id" form:"id" binding:"required,numeric"`
	Password    string `json:"password" form:"password" binding:"required"`
	NewPassword string `json:"newPassword" form:"newPassword" binding:"required"`
}

type UserDto struct {
	ID       uint64  `json:"id,string"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	NickName string  `json:"nickName,omitempty"`
	Email    string  `json:"email,omitempty"`
	DialCode string  `json:"dialCode,omitempty"`
	Phone    string  `json:"phone,omitempty"`
	Sex      *uint8  `json:"sex,omitempty"`
	Avatar   string  `json:"avatar,omitempty"`
	DeptId   *uint64 `json:"deptId"`
	Status   uint8   `json:"status"`
	Remark   string  `json:"remark,omitempty"`
	//创建时使用
	RoleIds []uint64 `json:"roleIds,omitempty" gorm:"-"`
	PostIds []uint64 `json:"postIds,omitempty" gorm:"-"`
}
