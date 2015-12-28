package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/asaskevich/govalidator"
	"github.com/codegangsta/cli"
	"github.com/denkhaus/rancher-meta-template/logging"
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
	User            string  `toml:"user"`
	Group           string  `toml:"group"`
	Check           Command `toml:"check"`
	Run             Command `toml:"run"`
}

type Config struct {
	Repeat   int           `toml:"repeat"`
	Host     string        `toml:"host"`
	Prefix   string        `toml:"prefix"`
	User     string        `toml:"user"`
	Group    string        `toml:"group"`
	LogLevel string        `toml:"loglevel"`
	Sets     []TemplateSet `toml:"set"`
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) Print() {
	logging.PrintInfo("check every %d seconds", cnf.Repeat)
	logging.PrintInfo("metadata host: %s", cnf.Host)
	logging.PrintInfo("prefix is: %s", cnf.Prefix)
	logging.PrintInfo("loglevel is: %s", cnf.LogLevel)
	logging.PrintInfo("run as %s:%s", cnf.User, cnf.Group)
	logging.PrintInfo(" %d template sets found", len(cnf.Sets))
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) Validate() error {
	if !govalidator.IsRequestURL(cnf.Host) {
		return errors.New("invalid host url")
	}

	if cnf.Repeat == 0 ||
		cnf.Prefix == "" ||
		cnf.User == "" ||
		cnf.Group == "" ||
		cnf.LogLevel == "" {
		return errors.New("invalid runtime options")
	}

	if len(cnf.Sets) == 0 {
		return errors.New("no template sets found")
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) OverrideFromCtx(ctx *cli.Context) {
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

//////////////////////////////////////////////////////////////////////////////
func NewFromCtx(ctx *cli.Context) (*Config, error) {
	templatePath := ctx.String("template")
	destinationPath := ctx.String("destination")

	if templatePath == "" {
		return nil, errors.New("no template path provided")
	}

	if destinationPath == "" {
		return nil, errors.New("no destination path provided")
	}

	cnf := new(Config)
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

	return cnf, nil
}

//////////////////////////////////////////////////////////////////////////////
func NewFromFile(configPath string) (*Config, error) {
	conf := new(Config)
	if _, err := os.Stat(configPath); err != nil {
		return nil, errors.New("config file not found")
	}

	if _, err := toml.DecodeFile(configPath, conf); err != nil {
		return nil, errors.Annotate(err, "decode config")
	}

	return conf, nil
}
