package main

import (
	"os"

	"github.com/BurntSushi/toml"
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
	Repeat int           `toml:"repeat"`
	Host   string        `toml:"host"`
	Prefix string        `toml:"prefix"`
	User   string        `toml:"user"`
	Group  string        `toml:"group"`
	Sets   []TemplateSet `toml:"set"`
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) Print() {
	printInfo("check every %d seconds", cnf.Repeat)
	printInfo("metadata host: %s", cnf.Host)
	printInfo("prefix is: %s", cnf.Prefix)
	printInfo("run as %s:%s", cnf.User, cnf.Group)
	printInfo(" %d template sets found", len(cnf.Sets))
}

//////////////////////////////////////////////////////////////////////////////
func (cnf *Config) Check() error {
	if cnf.Host == "" || cnf.Repeat == 0 || cnf.Prefix == "" {
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
