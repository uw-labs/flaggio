package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/victorkt/flaggio/cmd/admin"
	"github.com/victorkt/flaggio/cmd/api"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		admin.Command(),
		api.Command(),
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
