package main

import (
	"net/http"
	"os"

	"github.com/ernoaapa/linuxkit-server/pkg/version"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "linuxkit-server"
	app.Usage = "Server for building Linuxkit distributions"
	app.UsageText = `linuxkit-server [arguments...]

	 # By default listen port 5000
	 linuxkit-server`
	app.Description = `Server which builds Linuxkit distributions.`
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output in logs",
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
		router := mux.NewRouter()
		log.Println("Start listen on :8000")
		return http.ListenAndServe(":8000", router)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
