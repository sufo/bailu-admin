/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package admin

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/sufo/bailu-admin/app"
	"github.com/sufo/bailu-admin/global"
)

func StartCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "Start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App configuration file(.json,.yaml,.toml)",
				Required: false,
			},
			//&cli.StringFlag{
			//	Name:     "domain",
			//	Aliases:  []string{"m"},
			//	Usage:    "Casbin domain configuration(.conf)",
			//	Required: true,
			//},
			&cli.StringFlag{
				Name:  "www",
				Usage: "Static site directory",
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "Initialize menu's data configuration(.yaml)",
			},
		},
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetWWWDir(c.String("www")),
				app.SetMenuFile(c.String("menu")),
				app.SetVersion(global.Version))
		},
	}
}
