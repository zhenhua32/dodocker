package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const usage = `dodocker is a simple container runtime implementation.`

func main() {
	app := cli.NewApp()
	app.Name = "dodocker"
	app.Usage = usage

	app.Commands = []*cli.Command{
		&initCommand,
		&runCommand,
	}

	app.Before = func(c *cli.Context) error {
		logrus.SetFormatter(&logrus.TextFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal("发生致命错误", err)
	}
}
