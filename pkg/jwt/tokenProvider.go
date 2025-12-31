/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */
package jwt

import (
	"bailu/app/config"
	"bailu/global/consts"
	"bailu/pkg/store"
	"bailu/utils"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/google/wire"
	"strconv"
	"strings"
	"time"
)

var TokenInvalid = errors.New("invalid token")

var JwtProviderSet = wire.NewSet(wire.Struct(new(JwtProvider), "*"))

//var (
//	TokenExpired     = errors.New("Token is expired")
//	TokenNotValidYet = errors.New("Token not active yet")
//	TokenMalformed   = errors.New("That's not even a token")
//	TokenInvalid     = errors.New("invalid token")
//)

type JwtProvider struct {
	Store store.IStore
}

// subject (OnlineUserDto)
func GenerateToken(subject string) (string, error) {
	now := time.Now()
	conf := config.Conf.JWT
	expiresAt := now.Add(time.Duration(conf.Expired) * time.Second).Unix()

	token := jwt.NewWithClaims(signingMethod(conf), &jwt.StandardClaims{
		Id:        uuid.NewString(),
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		//Subject:   strconv.FormatUint(uint64(userId), 10),
		Subject: subject,
		Issuer:  "bailu", //发行者
	})
	return token.SignedString([]byte(conf.SigningKey))
}

func signingMethod(cfg config.JWT) jwt.SigningMethod {
	switch cfg.SigningMethod {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	default:
		return jwt.SigningMethodHS512
	}
}

// 解析令牌
func ParseToken(tokenStr string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(t *jwt.Token) (i interface{}, e error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, TokenInvalid
		}
		return []byte(config.Conf.JWT.SigningKey), nil
	})
	//判断错误类型，其实也可以直接返回return token.Claims.(*jwt.StandardClaims), err
	//if err != nil {
	//	if ve, ok := err.(*jwt.ValidationError); ok {
	//		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	//			return nil, TokenMalformed
	//		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
	//			// Token is expired
	//			return nil, TokenExpired
	//		} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
	//			return nil, TokenNotValidYet
	//		} else {
	//			return nil, TokenInvalid
	//		}
	//	}
	//}
	//if token != nil {
	//	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
	//		return claims, nil
	//	}
	//	return nil, TokenInvalid
	//
	//} else {
	//	return nil, TokenInvalid
	//}
	return token.Claims.(*jwt.StandardClaims), err
}

// 获取到期时间戳
// 得到10位时间戳
func GetExpireAt(tokenStr string) (int64, error) {
	claim, err := ParseToken(tokenStr)
	if err != nil {
		return 0, err
	}
	return claim.IssuedAt + config.Conf.JWT.Expired, nil
}

//func ParseOnlineUser(tokenStr string) (dto.OnlineUserDto, error) {
//	if tokenStr == "" {
//		return dto.OnlineUserDto{}, TokenInvalid
//	}
//	claims, err := ParseToken(tokenStr)
//	if err != nil {
//		return dto.OnlineUserDto{}, err
//	}
//	onlineUser := dto.OnlineUserDto{}
//	err2 := json.Unmarshal([]byte(claims.Subject), &onlineUser)
//	return onlineUser, err2
//}

// 通过token获取用户信息获取userId
//
//	func ParseUserID(token string) (int, error) {
//		user, err := ParseOnlineUser(token)
//		if err != nil {
//			return 0, err
//		}
//		return user.UserId, err
//	}

//	func ParseUserID(token string) (uint64, error) {
//		if token == "" {
//			return 0, TokenInvalid
//		}
//		claims, err := ParseToken(token)
//		if err != nil {
//			return 0, err
//		}
//		id, err := strconv.ParseUint(claims.Subject, 10, 64)
//
//		return id, err
//	}
func ParseUserID(token string) (uint64, error) {
	if token == "" {
		return 0, TokenInvalid
	}
	subject, err := ParseUserKey(token)
	if err != nil {
		return 0, err
	}
	subArr := strings.Split(subject, ":")
	id, err := strconv.ParseUint(subArr[len(subArr)-1], 10, 64)
	return id, err
}
func ParseUserKey(token string) (string, error) {
	if token == "" {
		return "", TokenInvalid
	}
	claims, err := ParseToken(token)
	if err != nil {
		return "", err
	}
	return claims.Subject, err
}

/**
 * @param token 需要检查的token
 * token续期 如果token没有被服务端储存，那么无法续期
 */
func (j *JwtProvider) CheckRenewal(token string) error {
	// 判断是否续期token,计算token的过期时间
	remaining, err := j.Store.TTL(config.Conf.JWT.OnlineKey + token)
	if err != nil {
		return err
	}
	// time.Now().Unix() 当前时间时间戳（10位）
	//判断当前时间与过期时间的时间差
	differ := int64(remaining * 1000)
	// 如果在续期检查的范围内，则续期
	if differ <= config.Conf.JWT.Detect {
		renew := remaining + time.Duration(config.Conf.JWT.Renew)
		err := j.Store.Expire(config.Conf.JWT.OnlineKey+token, renew)
		return err
	} else {
		return errors.New("续期失败！")
	}
}

/**
 * @param token 需要检查的token
 * token续期 如果token没有被服务端储存，那么无法续期
 */
//func (j *JwtProvider) CheckRenewalByUserId(userId string) exception {
//	// 判断是否续期token,计算token的过期时间
//	remaining, err := j.Store.TTL(config.Conf.JWT.OnlineKey + userId)
//	if err != nil {
//		return err
//	}
//	// time.Now().Unix() 当前时间时间戳（10位）
//	//判断当前时间与过期时间的时间差
//	differ := int64(remaining * 1000)
//	// 如果在续期检查的范围内，则续期
//	if differ <= config.Conf.JWT.Detect {
//		renew := remaining + time.Duration(config.Conf.JWT.Renew)
//		err := j.Store.Expire(config.Conf.JWT.OnlineKey+token, renew)
//		return err
//	} else {
//		return errors.New("续期失败！")
//	}
//}

func (j *JwtProvider) callStore(fn func(store.IStore) error) error {
	if store := j.Store; store != nil {
		return fn(store)
	}
	return nil
}

func (j *JwtProvider) Release() error {
	return j.callStore(func(iStore store.IStore) error {
		return iStore.Close()
	})
}

func UserKey(userAgent string, userId string) string {
	if config.Conf.Server.UseMultiDevice {
		isMobile := utils.IsMobile(userAgent)
		if isMobile {
			return fmt.Sprint(consts.DEVICE_MOBILE, ":", utils.MD5(userId))
		} else {
			return fmt.Sprint(consts.DEVICE_PC, ":", utils.MD5(userId))
		}
	} else { //只能一个帐号同时在一种设备上登录
		return fmt.Sprint(consts.DEVICE_ALL, ":", utils.MD5(userId))
	}
}

func UserKeyBy(isMobile bool, userId string) string {
	if config.Conf.Server.UseMultiDevice {
		if isMobile {
			return fmt.Sprint(consts.DEVICE_MOBILE, ":", utils.MD5(userId))
		} else {
			return fmt.Sprint(consts.DEVICE_PC, ":", utils.MD5(userId))
		}
	} else { //只能一个帐号同时在一种设备上登录
		return fmt.Sprint(consts.DEVICE_ALL, ":", utils.MD5(userId))
	}
}
