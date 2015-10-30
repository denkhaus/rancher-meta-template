package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	"gopkg.in/pipe.v2"

	"github.com/juju/errors"
	"github.com/rancher/go-rancher-metadata/metadata"
)

//////////////////////////////////////////////////////////////////////////////
func printError(err error) {
	fmt.Printf("rancher-meta-template::error: %s\n", errors.ErrorStack(err))
}

//////////////////////////////////////////////////////////////////////////////
func printInfo(format string, args ...interface{}) {
	fmt.Printf("rancher-meta-template::info: %s\n", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func createTemplateCtx(meta *metadata.Client) (interface{}, error) {

	services, err := meta.GetServices()
	if err != nil {
		return nil, errors.Annotate(err, "get services")
	}

	containers, err := meta.GetContainers()
	if err != nil {
		return nil, errors.Annotate(err, "get containers")
	}

	hosts, err := meta.GetHosts()
	if err != nil {
		return nil, errors.Annotate(err, "get hosts")
	}

	ctx := map[string]interface{}{
		"Services":   services,
		"Containers": containers,
		"Hosts":      hosts,
	}

	return ctx, nil
}

//////////////////////////////////////////////////////////////////////////////////
func appendCommandPipe(cmd Command, pipes []pipe.Pipe) []pipe.Pipe {
	if cmd.Cmd != "" {
		if cmd.Args != nil {
			return append(pipes, pipe.Exec(cmd.Cmd, cmd.Args...))
		}
		return append(pipes, pipe.Exec(cmd.Cmd))
	}

	return pipes
}

//////////////////////////////////////////////////////////////////////////////////
func processTemplateSet(templ *template.Template, meta *metadata.Client, set TemplateSet) error {
	buf, err := ioutil.ReadFile(set.TemplatePath)
	if err != nil {
		return errors.Annotate(err, "read template file")
	}

	tmpl, err := templ.Parse(string(buf))
	if err != nil {
		return errors.Annotate(err, "parse template")
	}

	ctx, err := createTemplateCtx(meta)
	if err != nil {
		return errors.Annotate(err, "create template context")
	}

	f, err := os.Create(set.DestinationPath)
	if err != nil {
		return errors.Annotate(err, "create destination file")
	}

	if err := tmpl.Execute(f, ctx); err != nil {
		return errors.Annotate(err, "execute template")
	}
	f.Close()

	printInfo("process check & run")

	pipes := make([]pipe.Pipe, 0)
	pipes = appendCommandPipe(set.Check, pipes)
	pipes = appendCommandPipe(set.Run, pipes)

	script := pipe.Script(pipes...)
	output, err := pipe.CombinedOutput(script)
	if err != nil {
		printInfo(string(output))
		return errors.Annotate(err, "check & run")
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////
func processTemplates(cnf *Config) error {
	meta := metadata.NewClient(cnf.Host)

	printInfo("connect rancher metadata host: %q", cnf.Host)
	tmpl := template.New("rancher-proxy").Funcs(newFuncMap())
	version := "init"

	for {
		newVersion, err := meta.GetVersion()
		if err != nil {
			return errors.Annotate(err, "get version")
		}

		if newVersion == version {
			time.Sleep(5 * time.Second)
			continue
		}

		version = newVersion
		printInfo("metadata changed - refresh config")

		for _, set := range cnf.Sets {
			if err := processTemplateSet(tmpl, meta, set); err != nil {
				return errors.Annotate(err, "process template set")
			}
		}

		time.Sleep(time.Duration(cnf.Repeat) * time.Second)
	}

	return nil
}
