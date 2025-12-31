/**
* Create by sufo
* @Email ouamour@Gmail.com
* @Desc 公共常量
 */

package consts

const (
	Version        = "1.0.0"
	ConfigEnv      = "BAILU_CONFIG"
	ConfigFile     = "config/config.yml"
	ConfigType     = "yaml"
	ConfigDefault  = "config/config.yml"
	REQUEST_ID_KEY = "X-Request-Id"

	CLAIM_KEY     = "claim"
	REQUEST_USER  = "Req-User"
	SUPER_ROLE_ID = 1
	SUPER_USER_ID = 1

	REQ_TOKEN = "REQ_TOKEN"

	//response code
	ERROR   = 1
	SUCCESS = 0

	/**
	 * 全部数据权限
	 */
	DATA_SCOPE_ALL = "1"

	/**
	 * 自定数据权限
	 */
	DATA_SCOPE_CUSTOM = "2"

	/**
	 * 部门数据权限
	 */
	DATA_SCOPE_DEPT = "3"

	/**
	 * 部门及以下数据权限
	 */
	DATA_SCOPE_DEPT_AND_CHILD = "4"

	/**
	 * 仅本人数据权限
	 */
	DATA_SCOPE_SELF = "5"
)

const (
	MODE_DEBUG   = "debug"
	MODE_TEST    = "test"
	MODE_RELEASE = "release"
)

const (
	TYPE_DIR    = "M"
	TYPE_MENU   = "C"
	TYPE_BUTTON = "F"
)

const SUPER_PERMISSION = "*:*:*"
const SUPER_KEY = "super"

type CaptchaType int

const (
	Audio CaptchaType = iota
	String
	Math
	Chines
	Digit
)

const DICT_CACHE_KEY = "sys_dict:"

const IMG_DIR = "imgs"

// PC MOBILE
const (
	DEVICE_ALL    = "ALL" //只能一个设备登录
	DEVICE_PC     = "PC"
	DEVICE_MOBILE = "MOBILE"
)

// USER:指定用户，ALL:全体用户
const (
	ANC_RECEIVER_ALL  = "ALL"
	ANC_RECEIVER_USER = "USER"
)
