package main

import (
	"os"

	"github.com/ernoaapa/linuxkit-server/pkg/api"
	"github.com/ernoaapa/linuxkit-server/pkg/version"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "linuxkit-server"
	app.Usage = "Server for building Linuxkit distributions"
	app.UsageText = `linuxkit-server [arguments...]

	 # By default listen port 8000
	 linuxkit-server
	 
	 # Use port number 8080
	 linuxkit-server --port 8080`
	app.Description = `Server which builds Linuxkit distributions.`
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
		},
		cli.IntFlag{
			Name:   "port",
			EnvVar: "PORT",
			Usage:  "HTTP port number",
			Value:  8000,
		},
	}
	app.Version = version.VERSION
	app.Before = func(context *cli.Context) error {
		if context.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Action = func(clicontext *cli.Context) error {
		return api.New(clicontext.Int("port")).Serve()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
