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
	"github.com/sufo/bailu-admin/app/api"
	"github.com/sufo/bailu-admin/app/core"
	"github.com/sufo/bailu-admin/app/domain/repo"
	"github.com/sufo/bailu-admin/app/router"
	"github.com/sufo/bailu-admin/app/service"
	"github.com/sufo/bailu-admin/app/service/cron"
	"github.com/sufo/bailu-admin/pkg/casbin"
	"github.com/sufo/bailu-admin/pkg/jwt"
	"github.com/sufo/bailu-admin/pkg/log"
	"github.com/sufo/bailu-admin/pkg/sms"
	"github.com/sufo/bailu-admin/pkg/store"
	"github.com/sufo/bailu-admin/utils/captcha"
	"github.com/sufo/bailu-admin/utils/dict"
	"github.com/sufo/bailu-admin/utils/upload"
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
