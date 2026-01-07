/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc
 */

package version

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/sufo/bailu-admin/global/consts"
)

func StartCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Get version info",
		Action: func(c *cli.Context) error {
			fmt.Println(consts.Version)
			return nil
		},
	}
}
