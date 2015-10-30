package main

import (
	"fmt"
	"os"

	"bitbucket.org/denkhaus/mirsvc/util"

	"github.com/BurntSushi/toml"
	"github.com/asaskevich/govalidator"
	"github.com/codegangsta/cli"
	"github.com/juju/errors"
)

type Command struct {
	Cmd  string   `toml:"cmd"`
	Args []string `toml:"args"`
}

type TemplateSet struct {
	Name            string  `toml:"name"`
	TemplatePath    string  `toml:"template"`
	DestinationPath string  `toml:"dest"`
	Check           Command `toml:"check"`
	Run             Command `toml:"run"`
}

type Config struct {
	Repeat int           `toml:"repeat"`
	Host   string        `toml:"host"`
	Sets   []TemplateSet `toml:"set"`
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) Print() {
	printInfo("check every %d seconds", cnf.Repeat)
	printInfo("metadata host: %s", cnf.Host)
	printInfo(" %d template sets found", len(cnf.Sets))

	util.Inspect(cnf)
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) Check() error {
	if cnf.Host == "" || cnf.Repeat == 0 {
		return errors.New("invalid runtime options")
	}

	if len(cnf.Sets) == 0 {
		return errors.New("no template sets provided")
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////
func readConfig(configPath string) (*Config, error) {
	conf := new(Config)
	if _, err := os.Stat(configPath); err != nil {
		return nil, errors.New("no configuration file not found")
	}

	if _, err := toml.DecodeFile(configPath, conf); err != nil {
		return nil, errors.Annotate(err, "create config")
	}

	return conf, nil
}

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
				cli.StringFlag{"destination, d", "", "destination path", "RANCHER_META_DEST_PATH"},
			},
			Action: func(ctx *cli.Context) {
				printInfo("startup")

				confPath := ctx.String("config")
				cnf, err := readConfig(confPath)
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
				}

				if !govalidator.IsRequestURL(cnf.Host) {
					printError(errors.New("provide a valid host url"))
					return
				}

				if err := cnf.Check(); err != nil {
					printError(errors.Annotate(err, "check config"))
					return
				}

				cnf.Print()
				if err := processTemplates(cnf); err != nil {
					printError(errors.Annotate(err, "process templates"))
				}
			},
		},
	}

	app.RunAndExitOnError()
}
