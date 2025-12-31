/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package store

import (
	"bailu/app/config"
	"bailu/global/consts"
	"bailu/utils"
	"github.com/tidwall/buntdb"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var _ IStore = (*BuntStore)(nil)

// NewStore 创建基于buntdb的文件存储
func NewBuntStore(path string) (*BuntStore, error) {
	if path != ":memory:" {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	//创建索引
	if err := db.CreateIndex(config.Conf.JWT.OnlineKey, config.Conf.JWT.OnlineKey+"*"); err != nil {
		return nil, err
	}
	if err := db.CreateIndex(consts.DICT_CACHE_KEY, consts.DICT_CACHE_KEY+"*"); err != nil {
		return nil, err
	}
	return &BuntStore{
		db: db,
	}, nil
}

// Store buntdb存储
type BuntStore struct {
	db *buntdb.DB
}

// Set ...
func (a *BuntStore) Set(key string, val interface{}, expire time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if expire > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expire}
		}
		_, _, err := tx.Set(key, utils.Strval(val), opts)
		return err
	})
}

// Delete 删除键
func (a *BuntStore) Del(key string) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})
}

// 注意这里的pattern是bundb里面已经创建好的索引名
// 索引支持pattern匹配
func (a *BuntStore) Clear(pattern string) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		//记录要删除的key
		deleteKeys := make([]string, 0)
		//不支持遍历删除
		err := tx.Ascend(pattern, func(key, value string) bool {
			deleteKeys = append(deleteKeys, key)
			return true
		})
		if err != nil {
			return err
		}

		// 删除
		for _, key := range deleteKeys {
			if _, _err := tx.Delete(key); _err != nil {
				return _err
			}
		}
		return nil
	})
}

func (a *BuntStore) Get(key string) (string, error) {
	var val string
	err := a.db.View(func(tx *buntdb.Tx) error {
		value, err := tx.Get(key)
		if err != nil {
			return err
		}
		val = value
		return nil
	})
	return val, err
}

func (a *BuntStore) Find(index, filter string) (KV, error) {
	kv := KV{}
	err := a.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(index, func(key, value string) bool {
			if strings.Contains(value, filter) {
				kv.setK(key).setV(value)
				return false
			}
			return true
		})
	})
	return kv, err
}

func (a *BuntStore) Scan(index string) ([]KV, error) {
	var kvs = make([]KV, 0)

	err := a.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(index, func(key, value string) bool {
			kv := KV{key, value}
			kvs = append(kvs, kv)
			return true
		})
	})
	return kvs, err
}

// 根据key获取存储值的有效期
func (a *BuntStore) TTL(key string) (time.Duration, error) {
	var ttl time.Duration
	err := a.db.View(func(tx *buntdb.Tx) error {
		value, err := tx.TTL(key)
		if err != nil {
			return err
		}
		ttl = value
		return nil
	})
	return ttl, err
}

/**
 * 获取有效时长,单位:秒(s)
 */
func (r *BuntStore) GetExpireAt(key string) (time.Duration, error) {
	ttl, err := r.TTL(key)
	if err != nil {
		return 0, err
	}
	return time.Duration(time.Now().Unix()) + ttl, nil
}

// 设置过期时长
func (a *BuntStore) Expire(key string, dur time.Duration) error {
	val, err := a.Get(key)
	if err != nil {
		return err
	}
	return a.Set(key, val, dur)
}

// Check ...
func (a *BuntStore) Check(key string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = val == "1"
		return nil
	})
	return exists, err
}

// Close ...
func (a *BuntStore) Close() error {
	return a.db.Close()
}
