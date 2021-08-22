package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zhenhua32/dodocker/cgroups/subsystems"
	"github.com/zhenhua32/dodocker/container"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `创建容器, 使用 namespace 和 cgroups
		dodocker run -ti [command]`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "ti",
			Usage: "启用 tty",
		},
		&cli.StringFlag{
			Name:  "m",
			Usage: "memoery limit",
		},
		&cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		&cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() < 1 {
			return fmt.Errorf("缺少容器命令")
		}
		var cmdArray = ctx.Args().Slice()
		tty := ctx.Bool("ti")
		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}
		Run(tty, cmdArray, resConf)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "初始化容器进程, 在容器中运行用户的进程. 不要在外部调用它",
	Action: func(ctx *cli.Context) error {
		logrus.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}
