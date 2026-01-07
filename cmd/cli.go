/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package cmd

import (
	"github.com/sufo/bailu-admin/app/config"
	"github.com/sufo/bailu-admin/app/core/appctx"
	"github.com/sufo/bailu-admin/cmd/admin"
	"github.com/sufo/bailu-admin/global"
	"github.com/sufo/bailu-admin/global/consts"
	"github.com/sufo/bailu-admin/utils"
	"context"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
)

////go:generate go env -w GO111MODULE=on
////go:generate go env -w GOPROXY=https://goproxy.cn,direct
////go:generate go mod tidy
////go:generate go mod download

// https://patorjk.com/software/taag/#p=testall&f=Graffiti&t=bailu
const name = `
   ▄▄▄▄███▄▄▄▄    ▄█  ███▄▄▄▄       ███        ▄████████ 
 ▄██▀▀▀███▀▀▀██▄ ███  ███▀▀▀██▄ ▀█████████▄   ███    ███ 
 ███   ███   ███ ███▌ ███   ███    ▀███▀▀██   ███    █▀  
 ███   ███   ███ ███▌ ███   ███     ███   ▀   ███        
 ███   ███   ███ ███▌ ███   ███     ███     ▀███████████ 
 ███   ███   ███ ███  ███   ███     ███              ███ 
 ███   ███   ███ ███  ███   ███     ███        ▄█    ███ 
  ▀█   ███   █▀  █▀    ▀█   █▀     ▄████▀    ▄████████▀  
`

func NewApp() *cli.App {
	ctx := appctx.NewTagCtx(context.Background(), "__main__")

	app := cli.NewApp()
	//app.Name = name //"bailu-admin"
	//version放到config文件里面
	//version := config.Conf.Version
	app.Version = global.Ternary(global.Version == "", consts.Version, global.Version)
	app.Usage = "bailu admin base on GIN + GORM + Casbin + JWT + WIRE"

	app.Commands = []*cli.Command{
		admin.StartCmd(ctx),
	}
	return app
}

func Run() {
	//skip 0: 当前栈帧信息； 1:当前调用栈上一帧信息
	global.Root = utils.GetCurrentAbPath(3)

	//加载配置文件, 因为这里要要到配置，所以把InitViper提前到这里
	_ = config.Conf.Default()
	//core.InitViper()
	var app = NewApp()
	err := app.Run(os.Args)
	if err != nil {
		logger, _ := zap.NewDevelopment()
		defer logger.Sync()
		logger.Error("err", zap.Error(err))
	}
}
