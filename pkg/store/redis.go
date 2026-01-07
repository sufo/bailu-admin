/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 如果需要排序则选择使用ZSET
 */

package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/pkg/log"
	"strings"
	"time"
)

// redis
var RedisClient *redis.Client

func NewRedisClient(db int) *redis.Client {
	conf := config.Conf.Store.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		client = nil
		log.L.Error("redis connect ping failed, err:", zap.Any("err", err))
		//return nil
		panic(err)
	}
	return client
}

var _ IStore = (*RedisStore)(nil)
var ctx = context.Background()

type RedisStore struct {
	Redis *redis.Client
}

/*
 * 指定缓存失效时间
 * duration 失效时长（秒）
 */
func (r *RedisStore) expire(key string, duration int64) bool {
	flag, err := r.Redis.Expire(ctx, key, time.Duration(duration)*time.Second).Result()
	if err != nil {
		panic(err)
	}
	return flag
}

func (r *RedisStore) Set(key string, val interface{}, duration time.Duration) error {
	return r.Redis.Set(ctx, key, val, duration).Err()
}

func (r *RedisStore) Get(key string) (string, error) {
	res, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	} else {
		return res, nil
	}
}

// pattern key模糊匹配
// filter value过滤
// return key value exception
func (r *RedisStore) Find(pattern, filter string) (KV, error) {
	var cursor uint64
	iter := r.Redis.Scan(ctx, cursor, pattern+"*", 1000).Iterator()
	for iter.Next(ctx) {
		val := iter.Val()
		res, err := r.Get(val)
		if err != nil {
			return KV{}, err
		}
		if strings.Contains(res, filter) {
			return KV{val, res}, nil
		}
		if cursor == 0 { // no more keys
			break
		}
	}
	return KV{}, errors.New("not found")
}

func (r *RedisStore) Scan(pattern string) ([]KV, error) {
	var cursor uint64
	var kvs = make([]KV, 0)
	iter := r.Redis.Scan(ctx, cursor, pattern+"*", 1000).Iterator()
	for iter.Next(ctx) {
		val := iter.Val()
		res, err := r.Get(val)
		if err != nil {
			return nil, err
		}
		kv := KV{val, res}
		kvs = append(kvs, kv)
		if cursor == 0 { // no more keys
			break
		}
	}
	return kvs, nil
}

func (r *RedisStore) DelByFilter(pattern, filter string) (string, error) {
	var cursor uint64
	iter := r.Redis.Scan(ctx, cursor, pattern+"*", 1000).Iterator()
	for iter.Next(ctx) {
		res, err := r.Get(iter.Val())
		if err != nil {
			return "", err
		}
		if strings.Contains(res, filter) {
			return res, nil
		}
		if cursor == 0 { // no more keys
			break
		}
	}
	return "", errors.New("not found")
}

// 判断key是否存在
func (r *RedisStore) Check(keys string) (bool, error) {
	n, err := r.Redis.Exists(ctx, keys).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// 判断key是否存在
func (r *RedisStore) Exists(keys ...string) (bool, error) {
	n, err := r.Redis.Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	//return int(n) > len(keys)-1, err
	return int(n) == len(keys), nil
}

// del
func (r *RedisStore) MultiDel(keys ...string) (bool, error) {
	n, err := r.Redis.Del(ctx, keys...).Result()
	if err != nil {
		return false, err
	}
	//return int(n) > len(keys)-1, err
	return int(n) == len(keys), nil
}

// del
func (r *RedisStore) Del(key string) error {
	n, err := r.Redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if int(n) > 0 {
		return nil
	} else {
		return errors.New("delete failed")
	}
}

// @pattern key前缀
func (r *RedisStore) Clear(pattern string) error {
	var cursor uint64
	//如果匹配的数据量大不要用Keys，Keys命令具有阻塞特性，若匹配的key数量过多，则会占用很多服务器资源，导致Redis性能下降。
	iter := r.Redis.Scan(ctx, cursor, pattern+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		// 使用 DEL 命令删除 key
		if err := r.Redis.Del(ctx, key).Err(); err != nil {
			log.L.Errorf("Failed to delete key %s: %v", key, err)
			return err
		} else {
			if config.Conf.IsDebug() {
				fmt.Printf("Deleted key: %s\n", key)
			}
		}
	}
	return iter.Err()
}

/**
 * 获取剩余过期时长,单位:秒(s)
 */
func (r *RedisStore) TTL(key string) (time.Duration, error) {
	ttl, err := r.Redis.TTL(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	return ttl, nil
}

// 过期时间戳  （精确到秒）
func (r *RedisStore) GetExpireAt(key string) (time.Duration, error) {
	ttl, err := r.TTL(key)
	if err != nil {
		return 0, err
	}
	return time.Duration(time.Now().Unix()) + ttl, nil
}

// 设置过期时间
func (r *RedisStore) Expire(key string, dur time.Duration) error {
	res, err := r.Redis.Expire(ctx, key, dur).Result()
	if err != nil {
		return err
	}
	if !res {
		return errors.New("设置过期失败")
	}
	return nil
}

// Close ...
func (r *RedisStore) Close() error {
	return r.Redis.Close()
}

func (r *RedisStore) WrapFunc(fn func(client *redis.Client) (interface{}, error)) (interface{}, error) {
	return fn(r.Redis)
}
