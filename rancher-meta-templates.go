package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"

	"gopkg.in/pipe.v2"

	"github.com/juju/errors"
	"github.com/rancher/go-rancher-metadata/metadata"
)

const (
	PROXY_TEMPLATE_PATH = "/etc/rancher-proxy/templates/nginx.tmpl"
	PROXY_CONFIG_PATH   = "/etc/nginx/sites-enabled/rancher-proxy.conf"
)

//////////////////////////////////////////////////////////////////////////////
func printError(err error) {
	fmt.Printf("rancher-proxy::error: %s\n", errors.ErrorStack(err))
}

//////////////////////////////////////////////////////////////////////////////
func printInfo(format string, args ...interface{}) {
	fmt.Printf("rancher-proxy::info: %s\n", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func createTemplateCtx(meta *metadata.Client) (map[string][]interface{}, error) {

	data, err := meta.GetServices()
	if err != nil {
		return nil, errors.Annotate(err, "get services")
	}

	tmplData := make(map[string][]interface{})
	for _, service := range data {
		labels := service.Labels

		for _, value := range labels {

			port, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.Annotate(err, "convert port from metadata")
			}

			for _, containerName := range service.Containers {
				containers, err := meta.GetContainers()
				if err != nil {
					return nil, errors.Annotate(err, "get containers")
				}

				for _, container := range containers {
					if container.Name != containerName {
						continue
					}

					d := map[string]interface{}{
						"ip":            container.PrimaryIp,
						"port":          port,
						"containerName": containerName,
					}

					serviceName := service.Name
					printInfo("expose service %q:%d in container %q",
						serviceName, port, containerName)

					if _, ok := tmplData[serviceName]; ok {
						tmplData[serviceName] =
							append(tmplData[serviceName], d)
					} else {
						dat := []interface{}{d}
						tmplData[serviceName] = dat
					}
				}
			}
		}
	}

	return tmplData, nil
}

//////////////////////////////////////////////////////////////////////////////
func processTemplateSet(meta *metadata.Client, set TemplateSet) error {
	buf, err := ioutil.ReadFile(set.TemplatePath)
	if err != nil {
		return errors.Annotate(err, "read template file")
	}

	tmpl, err := template.New("rancher-proxy").Parse(string(buf))
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

	printInfo("restart nginx")
	script := pipe.Script(
		pipe.Exec("nginx", "-t"),
		pipe.Exec("nginx", "-s", "reload"),
	)

	output, err := pipe.CombinedOutput(script)
	if err != nil {
		printInfo(string(output))
		return errors.Annotate(err, "restart nginx")
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////
func processTemplates(cnf *Config) error {
	host := path.Join(cnf.Host, cnf.Prefix)
	meta := metadata.NewClient(host)

	printInfo("connect rancher metadata host: %q", host)

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
			if err := processTemplateSet(meta, set); err != nil {
				return errors.Annotate(err, "process template set")
			}
		}

		time.Sleep(time.Duration(cnf.Repeat) * time.Second)
	}

	return nil
}
