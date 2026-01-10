/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package app

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/sufo/bailu-admin/app/config"
	"time"
)

func InitCasbin(adapter persist.Adapter) (*casbin.SyncedEnforcer, func(), error) {
	cfg := config.Conf.Casbin
	if cfg.Model == "" {
		return new(casbin.SyncedEnforcer), func() {}, nil
	}
	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	if err != nil {
		return nil, nil, err
	}
	e.EnableLog(cfg.Debug)
	err = e.InitWithModelAndAdapter(e.GetModel(), adapter)
	if err != nil {
		return nil, nil, err
	}
	e.EnableEnforce(cfg.Enable)
	clearFunc := func() {}
	if cfg.AutoLoad {
		e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInterval) * time.Second)
		clearFunc = func() {
			e.StopAutoLoadPolicy()
		}
	}
	return e, clearFunc, nil
}
