//go:build wireinject

//According to the Bug-resistant build constraints
//proposal, //+build will be replaced by //go:build.
//A transition period from //+build to //go:build syntax
//will last from Go version 1.16 through version 1.18.
//In the 1.16 version of Go, you can use either the old syntax or both syntaxes at the same time.
// the build tag makes sure the stub is not build in the final build

/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */
package app

import (
	"bailu/app/api"
	"bailu/app/core"
	"bailu/app/domain/repo"
	"bailu/app/router"
	"bailu/app/service"
	"bailu/app/service/cron"
	"bailu/pkg/casbin"
	"bailu/pkg/jwt"
	"bailu/pkg/log"
	"bailu/pkg/sms"
	"bailu/pkg/store"
	"bailu/utils/captcha"
	"bailu/utils/dict"
	"bailu/utils/upload"
	"github.com/google/wire"
)

func BuildInjector(www string) (*Injector, func(), error) {
	wire.Build(
		log.InitLogger,
		//validator.New,
		sms.New,
		core.InitGorm,
		upload.NewOSS,
		store.NewStore,
		repo.RepoSet,
		dict.DictUtilSet,
		captcha.NewDefaultRedisStore,
		jwt.JwtProviderSet,
		casbin.CasbinAdapterSet,
		InitCasbin,
		//service
		service.ServiceSet,
		cron.CrontabSet,
		api.APISet,
		router.RouterSet,
		core.InitRouter,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
