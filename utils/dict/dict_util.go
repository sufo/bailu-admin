/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 字典tool
 */

package dict

import (
	"encoding/json"
	"github.com/google/wire"
	"github.com/sufo/bailu-admin/app/domain/entity"
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/pkg/store"
	"time"
)

var DictUtilSet = wire.NewSet(wire.Struct(new(DictUtil), "*"))

type DictUtil struct {
	Store store.IStore
}

func (d *DictUtil) SetDictCache(key string, items []entity.DictItem, expiration time.Duration) error {
	return (d.Store).Set(cacheKey(key), items, expiration)
}

func (d *DictUtil) GetDictCache(key string) ([]entity.DictItem, error) {
	res, err := (d.Store).Get(cacheKey(key))
	if err != nil {
		return nil, err
	} else {
		var items = make([]entity.DictItem, 0)
		if err := json.Unmarshal([]byte(res), &items); err != nil {
			return nil, err
		} else {
			return items, nil
		}
	}
}

func (d *DictUtil) RemoveDictCache(key string) error {
	return (d.Store).Del(cacheKey(key))
}

func (d *DictUtil) Clear(pattern string) error {
	return (d.Store).Clear(pattern)
}

// 缓存key
func cacheKey(key string) string {
	return consts.DICT_CACHE_KEY + key
}
