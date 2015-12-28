package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/denkhaus/rancher-meta-template/config"
	"github.com/denkhaus/rancher-meta-template/logging"
	"github.com/juju/errors"
)

//////////////////////////////////////////////////////////////////////////////
func main() {
	app := cli.NewApp()
	app.Name = "rancher-meta-template"
	app.Version = fmt.Sprintf("%s-%s", AppVersion, Revision)
	app.Usage = "A rancher-metadata config writer."

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "run config generation",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "config, c", Value: "/etc/rancher-meta-template/config.toml", Usage: "Configuration File path"},
				cli.IntFlag{Name: "repeat, r", Value: 60, Usage: "Repeat config creation every x seconds", EnvVar: "RANCHER_META_REPEAT"},
				cli.StringFlag{Name: "host, H", Value: "http://rancher-metadata", Usage: "Rancher metadata host", EnvVar: "RANCHER_META_HOST"},
				cli.StringFlag{Name: "template, t", Usage: "template path", EnvVar: "RANCHER_META_TEMPLATE_PATH"},
				cli.StringFlag{Name: "prefix, p", Value: "/latest", Usage: "api prefix", EnvVar: "RANCHER_META_PREFIX"},
				cli.StringFlag{Name: "destination, d", Usage: "the destination path", EnvVar: "RANCHER_META_DEST_PATH"},
				cli.StringFlag{Name: "user, u", Value: "nouser", Usage: "run as user", EnvVar: "RANCHER_META_USER"},
				cli.StringFlag{Name: "group, g", Value: "nogroup", Usage: "run as group", EnvVar: "RANCHER_META_GROUP"},
				cli.StringFlag{Name: "loglevel, l", Value: "warning", Usage: "the loglevel", EnvVar: "RANCHER_META_LOGLEVEL"},
			},
			Action: func(ctx *cli.Context) {
				logging.Info("startup")

				var cnf *config.Config
				cnf, _ = config.NewFromFile(ctx.String("config"))
				if cnf == nil {
					logging.Warning("config file not found, get config from cli context")
					c, err := config.NewFromCtx(ctx)
					if err != nil {
						logging.Error(errors.Annotate(err, "new config from context"))
						return
					}
					cnf = c

				} else {
					cnf.OverrideFromCtx(ctx)
				}

				cnf.Print()
				if err := cnf.Validate(); err != nil {
					logging.Error(errors.Annotate(err, "validate config"))
					return
				}

				if err := logging.SetLogLevel(cnf.LogLevel); err != nil {
					logging.Error(errors.Annotate(err, "set log level"))
					return
				}

				if err := processTemplates(cnf); err != nil {
					logging.Error(errors.Annotate(err, "process templates"))
				}
			},
		},
	}

	app.RunAndExitOnError()
}
