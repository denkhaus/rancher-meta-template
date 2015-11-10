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
				cli.StringFlag{"config, c", "/etc/rancher-meta-template/config.toml", "Configuration File path", ""},
				cli.IntFlag{"repeat, r", 60, "Repeat config creation every x seconds", "RANCHER_META_REPEAT"},
				cli.StringFlag{"host, H", "http://rancher-metadata", "Rancher metadata host", "RANCHER_META_HOST"},
				cli.StringFlag{"template, t", "", "template path", "RANCHER_META_TEMPLATE_PATH"},
				cli.StringFlag{"prefix, p", "/latest", "api prefix", "RANCHER_META_PREFIX"},
				cli.StringFlag{"destination, d", "", "the destination path", "RANCHER_META_DEST_PATH"},
				cli.StringFlag{"user, u", "nouser", "run as user", "RANCHER_META_USER"},
				cli.StringFlag{"group, g", "nogroup", "run as group", "RANCHER_META_GROUP"},
				cli.StringFlag{"loglevel, l", "warning", "the loglevel", "RANCHER_META_LOGLEVEL"},
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
					cnf.Group = ctx.String("loglevel")
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
