package main

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/codegangsta/cli"
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
				printInfo("startup")

				cnf, err := readConfig(ctx.String("config"))
				if err != nil {
					printError(errors.Annotate(err, "read config"))
					return
				}

				if cnf == nil {
					templatePath := ctx.String("template")
					destinationPath := ctx.String("destination")

					if templatePath == "" {
						printError(errors.New("no template path provided"))
						return
					}

					if destinationPath == "" {
						printError(errors.New("no destination path provided"))
						return
					}

					cnf = new(Config)
					cnf.Repeat = ctx.Int("repeat")
					cnf.Host = ctx.String("host")
					cnf.Prefix = ctx.String("prefix")
					cnf.User = ctx.String("user")
					cnf.Group = ctx.String("group")
					cnf.LogLevel = ctx.String("loglevel")
					cnf.Sets = make([]TemplateSet, 0)

					cnf.Sets = append(cnf.Sets, TemplateSet{
						TemplatePath:    templatePath,
						DestinationPath: destinationPath,
					})

				} else {
					if cnf.Repeat == 0 || ctx.IsSet("repeat") {
						cnf.Repeat = ctx.Int("repeat")
					}
					if cnf.Host == "" || ctx.IsSet("host") {
						cnf.Host = ctx.String("host")
					}
					if cnf.Prefix == "" || ctx.IsSet("prefix") {
						cnf.Host = ctx.String("prefix")
					}
					if cnf.User == "" || ctx.IsSet("user") {
						cnf.User = ctx.String("user")
					}
					if cnf.Group == "" || ctx.IsSet("group") {
						cnf.Group = ctx.String("group")
					}
					if cnf.LogLevel == "" || ctx.IsSet("loglevel") {
						cnf.LogLevel = ctx.String("loglevel")
					}
				}

				if !govalidator.IsRequestURL(cnf.Host) {
					printError(errors.New("provide a valid host url"))
					return
				}

				cnf.Print()
				if err := cnf.Check(); err != nil {
					printError(errors.Annotate(err, "check config"))
					return
				}

				if err := setLogLevel(cnf.LogLevel); err != nil {
					printError(errors.Annotate(err, "set log level"))
					return
				}

				if err := processTemplates(cnf); err != nil {
					printError(errors.Annotate(err, "process templates"))
				}
			},
		},
	}

	app.RunAndExitOnError()
}
