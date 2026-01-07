/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 数据存储（redis/buntdb） 主要用于存储登录用户
 */

package store

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/google/wire"
	"time"
)

type IStore interface {
	Get(key string) (string, error)
	Set(key string, val interface{}, duration time.Duration) error
	/**
	 * 检查key是否存在
	 */
	Check(key string) (bool, error)
	/**
	 * 获取过期时间戳
	 */
	//GetExpireAt(key string) (time.Duration, exception)

	/**
	 * 根据filter查询value
	 */
	Find(key, filter string) (KV, error)

	Scan(key string) ([]KV, error)
	/**
	 * 获取剩余过期时长
	 */
	TTL(key string) (time.Duration, error)
	/**
	 * 设置过期时间
	 */
	Expire(key string, dur time.Duration) error

	Del(key string) error
	/**
	 * 根据filter查询内容，返回内容对应的key
	 */
	//DelByFilter(filter string) (string, exception)

	//根据key前缀删除
	Clear(pattern string) error

	Close() error
}

//func NewStore() IStore {
//	cfg := config.Conf.JWT
//	var store interface{}
//	switch cfg.Store {
//	case "redis":
//		store = redis.NewRedis()
//	default:
//		s, err := buntdb.NewStore(cfg.FilePath)
//		if err != nil {
//			store = nil
//		}
//		store = s
//	}
//	return store.(IStore)
//}

func NewStore() IStore {
	cfg := config.Conf.Store
	//if cfg.Enable {
	var store interface{}
	if config.Conf.Server.UseRedis && cfg.StoreType == "redis" {
		redisClient := RedisClient
		if redisClient == nil {
			redisClient = NewRedisClient(cfg.Redis.DB)
		}
		store = &RedisStore{redisClient}
		wire.Bind(new(IStore), new(*RedisStore))
	} else {
		s, err := NewBuntStore(cfg.BuntDb.FilePath)
		if err != nil {
			store = nil
		}
		store = s
		wire.Bind(new(IStore), new(*BuntStore))
	}
	//println(store)
	return store.(IStore)
}

type KV struct {
	K string
	V string
}

func (kv *KV) setK(k string) *KV {
	kv.K = k
	return kv
}
func (kv *KV) setV(v string) *KV {
	kv.K = v
	return kv
}
