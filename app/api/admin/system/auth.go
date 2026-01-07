/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc login
 */

package system

import (
	"github.com/sufo/bailu-admin/app/api/admin"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/domain/dto"
	"github.com/sufo/bailu-admin/app/domain/resp"
	"github.com/sufo/bailu-admin/app/domain/vo"
	"github.com/sufo/bailu-admin/app/service/sys"
	"github.com/sufo/bailu-admin/global/consts"
	respErr "github.com/sufo/bailu-admin/pkg/exception"
	"github.com/sufo/bailu-admin/pkg/i18n"
	"github.com/sufo/bailu-admin/pkg/jwt"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/pkg/rsa"
	"github.com/sufo/bailu-admin/pkg/translate"
	captcha2 "github.com/sufo/bailu-admin/utils/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jinzhu/copier"
	captcha "github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

var AuthSet = wire.NewSet(wire.Struct(new(AuthApi), "*"))

type AuthApi struct {
	CaptchaStore *captcha2.CaptchaStore
	LoginSrv     *sys.AuthService
	OnlineSrv    *sys.OnlineService
	SMSService   *sys.SMSService
	LoginLogSrv  *sys.LoginLogService
}

//func (a *AuthApi) Login(c *gin.Context) {
//	var loginParams dto.LoginParams
//	if err := c.ShouldBindJSON(&loginParams); err != nil {
//		panic(err)
//	}
//	if !a.CaptchaStore.Verify(loginParams.CaptchaId, loginParams.ImgCode, true) {
//		resp.FailWithMsg(c, i18n.DefTr("tip.captchaValid"))
//		return
//	}
//	//处理密码
//	decryptStr, err := rsa.PrivateDecrypt(loginParams.Password, config.Conf.RSA.PrivateKey)
//	if err != nil {
//		resp.FailWithMsg(c, i18n.DefTr("admin.AccountOrPwdErr"))
//		return
//	}
//	loginUser := dto.LoginUser{loginParams.UserName, decryptStr}
//	if user, err := a.LoginSrv.Login(c, loginUser); err != nil {
//		panic(respErr.BadRequestError)
//	} else {
//		token, err := jwt.GenerateToken(strconv.FormatUint(user.ID, 10))
//		if err != nil {
//			panic(respErr.InternalServerErrorWithMsg)
//		}
//		//判断是否为多点登录
//		if !config.Conf.Server.UseMultipoint {
//			//不是，则踢出之前登录的用户
//			if err := a.OnlineSrv.KickOut(strconv.FormatUint(user.ID, 10)); err != nil {
//				log.L.Error(err)
//			}
//		}
//		saveErr := a.OnlineSrv.Save(user, c.Request, token)
//		if saveErr != nil {
//			resp.FailWithError(c, saveErr)
//			return
//		}
//		resp.OKWithData(c, map[string]any{
//			"token":    token,
//			"userInfo": user,
//			"expire":   config.Conf.JWT.Expired,
//		})
//	}
//}

func (a *AuthApi) _login(c *gin.Context, hasCapatch bool) {
	var loginUser dto.LoginUser
	//是否有验证码
	if hasCapatch {
		var loginParams dto.LoginParams
		if err := c.ShouldBindJSON(&loginParams); err != nil {
			panic(err)
		}
		if !a.CaptchaStore.Verify(loginParams.CaptchaId, loginParams.ImgCode, true) {
			msg := i18n.DefTr("tip.captchaValid")
			//记录日志
			a.LoginLogSrv.Create(c, loginUser.Username, consts.ERROR, msg)

			resp.FailWithMsg(c, msg)
			return
		}
		loginUser.Username = loginParams.Username
		loginUser.Password = loginParams.Password
	} else {
		if err := c.ShouldBindJSON(&loginUser); err != nil {
			panic(err)
		}
	}
	//处理密码
	decryptStr, err := rsa.PrivateDecrypt(loginUser.Password, config.Conf.RSA.PrivateKey)
	if err != nil {
		msg := i18n.DefTr("admin.AccountOrPwdErr")
		//记录日志
		a.LoginLogSrv.Create(c, loginUser.Username, consts.ERROR, msg)

		resp.FailWithMsg(c, msg)
		return
	}
	loginUser.Password = decryptStr
	if user, err := a.LoginSrv.Login(c, loginUser); err != nil {
		a.LoginLogSrv.Create(c, loginUser.Username, consts.ERROR, err.Error())
		panic(err)
	} else {
		//token, err := jwt.GenerateToken(strconv.FormatUint(user.ID, 10))
		uKey := jwt.UserKey(c.Request.UserAgent(), utils.ToString(user.ID))
		token, err := jwt.GenerateToken(uKey)
		if err != nil {
			a.LoginLogSrv.Create(c, loginUser.Username, consts.ERROR, err.Error())
			panic(respErr.InternalServerError)
		}

		_, saveErr := a.OnlineSrv.Save(user, c.Request, uKey)
		if saveErr != nil {
			resp.FailWithError(c, saveErr)
			return
		}

		//用户登录日志
		a.LoginLogSrv.Create(c, loginUser.Username, consts.SUCCESS, "登录成功")

		//处理userVo roles
		var userVo = &vo.User{}
		var roles = make([]vo.KV, 0)
		copier.Copy(userVo, user)
		for _, r := range user.Roles {
			op := vo.KV{r.RoleKey, r.Name}
			roles = append(roles, op)
		}
		userVo.Roles = roles

		//post
		var posts = make([]vo.KV, 0)
		for _, p := range user.Posts {
			op := vo.KV{p.PostCode, p.Name}
			posts = append(posts, op)
		}
		userVo.Posts = posts
		//avatar
		if user.Avatar != "" && !strings.HasPrefix(user.Avatar, "http") {
			userVo.Avatar = admin.FileUrl(c.Request, user.Avatar)
		}

		resp.OKWithData(c, map[string]any{
			"token":    token,
			"userInfo": userVo,
			"expire":   config.Conf.JWT.Expired,
		})
	}
}

// @title 管理员登录
// @Summary 管理员登录
// @Description 管理员登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginParams true "请求体"
// //@Success 200 {object} resp.Response[any]{data=object{token=string,userInfo=vo.User,expire=integer}}
// @Router /api/login [post]
func (a *AuthApi) Login(c *gin.Context) {
	a._login(c, true)
}

// @title UnLock 管理员屏幕解锁
// @Summary 管理端屏幕解锁
// @Description 管理端屏幕解锁
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginParams true "请求体"
// @Failure      1  {object}  resp.Response[any]
// //@Success 200 {object} resp.Response[any]{data=object{token=string,userInfo=vo.User,expire=integer}}
// @Router /api/unlock [post]
// @Security Bearer
func (a *AuthApi) UnLock(c *gin.Context) {
	a._login(c, false)
}

/**
 * 生成验证码
 */
// @title 生成验证码
// @Summary 用于生成登录验证码
// @Description 用于生成登录验证码
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
////@Success 200 {object} resp.Response[any]{data=object{captchaId=string,picPath=string}}
// @Router /api/captcha [get]
func (a *AuthApi) GetCaptcha(c *gin.Context) {
	var driver captcha.Driver
	conf := config.Conf.Captcha
	//create base64 encoding captcha
	switch conf.CaptchaType {
	case "audio":
		driver = &captcha.DriverAudio{}
	case "string":
		driver = (&captcha.DriverString{
			Height: conf.Height,
			Width:  conf.Width,
			Length: conf.Length,
		}).ConvertFonts()
	case "math":
		driver = (&captcha.DriverMath{
			Height: conf.Height,
			Width:  conf.Width,
		}).ConvertFonts()
	case "chinese":
		driver = (&captcha.DriverChinese{
			Height: conf.Height,
			Width:  conf.Width,
		}).ConvertFonts()
	default:
		driver = &captcha.DriverDigit{
			Height: conf.Height,
			Width:  conf.Width,
			Length: conf.Length,
		}
	}
	var store captcha.Store
	if conf.Store == "memory" {
		//store = captcha.DefaultMemStore
		store = captcha.NewMemoryStore(10204, time.Duration(conf.Expire))
	} else {
		store = a.CaptchaStore
	}
	res := captcha.NewCaptcha(driver, store)
	if id, b64s, _, err := res.Generate(); err != nil {
		log.L.Error(i18n.DefTr("admin.captchaError"), zap.Error(err))
		resp.FailWithMsg(c, i18n.DefTr("admin.captchaError"))
		return
	} else {
		resp.Result(c, resp.Data(map[string]interface{}{
			"captchaId": id,
			"picPath":   b64s,
			//"length":    config.Conf.Captcha.Length,
		}), resp.Msg(i18n.DefTr("tip.captchaOk")))
	}
}

// @title 发送短信验证码
// @Summary 发送短信验证码
// @Description 发送短信验证码
// @Tags Auth
// @Accept json
// @Produce json
// @Param phone query string true "手机号"
// @response 200 {object} resp.Response[any]
// @Router /api/sms [post]
func (a *AuthApi) SendSMS(c *gin.Context) {
	phone := c.Query("phone")
	countryCode := admin.Query(c, "dialCode", "86")
	if valid := translate.Validate.Var("+"+countryCode+phone, "required,e164"); valid != nil {
		//panic(respErr.BadRequestErrorWithError(valid))
		panic(valid) //Validate错误直接抛出交给中间件处理，包要用respErr包装，包装后不会自动翻译
	}
	//非测试
	if config.Conf.Server.Mode != "debug" {
		if err := a.SMSService.SendSMS(c, phone, countryCode); err != nil {
			resp.FailWithError(c, err)
			return
		}
	}
	resp.OkWithMsg(c, i18n.DefTr("tip.smsSendOk"))
}

// TODO
// @title 短信验证码登录
// @Summary 手机好和短信验证码登录
// @Description 手机好和短信验证码登录
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body dto.SMSLoginParams true "请求体"
// //@Success 200 {object} resp.Response[any]{data=object{captchaId=string,picPath=string}}
// @Router /api/sms-login [post]
func (a *AuthApi) SMSLogin(c *gin.Context) {
	var loginUser dto.SMSLoginParams
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		panic(respErr.BadRequestError)
	}

	resp.OKWithData(c, map[string]any{
		"token":    "token",
		"userInfo": "user",
		"expired":  config.Conf.JWT.Expired,
	})
}
