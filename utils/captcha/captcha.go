/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package captcha

import (
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/pkg/store"
	"time"
)

var _ base64Captcha.Store = (*CaptchaStore)(nil)

func NewDefaultRedisStore(store store.IStore) *CaptchaStore {
	conf := config.Conf.Captcha
	if store == nil {
		return nil
	}
	return &CaptchaStore{
		Expiration: time.Second * time.Duration(conf.Expire),
		PreKey:     conf.Prefix,
		Store:      store,
	}
}

type CaptchaStore struct {
	Expiration time.Duration
	PreKey     string
	Store      store.IStore
}

func (rs *CaptchaStore) Set(id string, value string) error {
	err := rs.Store.Set(rs.PreKey+id, value, rs.Expiration)
	return err
}

func (rs *CaptchaStore) Get(key string, clear bool) string {
	val, err := rs.Store.Get(key)
	if err != nil {
		log.L.Error("RedisStoreGetError!", zap.Error(err))
		return ""
	}
	if clear {
		err := rs.Store.Del(key)
		if err != nil {
			log.L.Error("RedisStoreClearError!", zap.Error(err))
			return ""
		}
	}
	return val
}

func (rs *CaptchaStore) Verify(id, answer string, clear bool) bool {
	key := rs.PreKey + id
	v := rs.Get(key, clear)
	return v == answer
}
