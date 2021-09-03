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
		dodocker run --ti [command]`,
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
		&cli.StringFlag{
			Name:  "cwd",
			Usage: "设置子进程的当前工作目录",
		},
		&cli.StringFlag{
			Name:  "v",
			Usage: "volume, host:container",
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
		Run(tty, cmdArray, resConf, ctx.String("cwd"), ctx.String("v"))
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

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "创建一个新的容器镜像",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "cwd",
			Usage: "设置子进程的当前工作目录",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.Args().Len() < 1 {
			return fmt.Errorf("缺少容器的镜像名称")
		}
		imageName := ctx.Args().Get(0)
		commitContainer(imageName, ctx.String("cwd"))
		return nil
	},
}
