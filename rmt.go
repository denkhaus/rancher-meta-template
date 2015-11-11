package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/denkhaus/rancher-meta-template/scratch"
	"github.com/juju/errors"
	"github.com/rancher/go-rancher-metadata/metadata"
	"gopkg.in/pipe.v2"
)

const (
	DEFAULT_TEMPLATE_DIR = "/etc/rancher-meta-template/templates"
)

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

	containersW := make([]ContainerWrap, 0)
	for _, container := range containers {
		cw := ContainerWrap{container}
		containersW = append(containersW, cw)
	}

	servicesW := make([]ServiceWrap, 0)
	for _, service := range services {
		sw := ServiceWrap{service}
		servicesW = append(servicesW, sw)
	}

	ctx := map[string]interface{}{
		"Services":   servicesW,
		"Containers": containersW,
	}

	return ctx, nil
}

//////////////////////////////////////////////////////////////////////////////
func computeMd5(filePath string) (string, error) {
	if _, err := os.Stat(filePath); err != nil {
		return "", nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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
func processTemplateSet(meta *metadata.Client, set TemplateSet) error {

	if _, err := os.Stat(set.TemplatePath); err != nil {
		printWarning("template path %q is not available: skip", set.TemplatePath)
		return nil
	}

	buf, err := ioutil.ReadFile(set.TemplatePath)
	if err != nil {
		return errors.Annotate(err, "read template file")
	}

	templ := template.New(set.Name).Funcs(newFuncMap())
	tmpl, err := templ.Parse(string(buf))
	if err != nil {
		return errors.Annotate(err, "parse template")
	}

	lastMd5 := ""
	if _, err = os.Stat(set.DestinationPath); err == nil {
		lmd5, err := computeMd5(set.DestinationPath)
		if err != nil {
			return errors.Annotate(err, "get last md5 hash")
		}
		lastMd5 = lmd5
	}

	ctx, err := createTemplateCtx(meta)
	if err != nil {
		return errors.Annotate(err, "create template context")
	}

	var newBuf bytes.Buffer
	wr := bufio.NewWriter(&newBuf)
	if err := tmpl.Execute(wr, ctx); err != nil {
		return errors.Annotate(err, "execute template")
	}

	if err := wr.Flush(); err != nil {
		return errors.Annotate(err, "flush tmpl writer")
	}

	hash := md5.New()
	hash.Write(newBuf.Bytes())
	currentMd5 := fmt.Sprintf("%x", hash.Sum(nil))

	if lastMd5 == currentMd5 {
		return nil
	}

	if lastMd5 == "" {
		printInfo("create output file")
	} else {
		printInfo("last md5 sum is %q", lastMd5)
		printInfo("current md5 sum is %q", currentMd5)
		printInfo("output file needs refresh")
	}

	f, err := os.Create(set.DestinationPath)
	if err != nil {
		return errors.Annotate(err, "create destination file")
	}

	w := bufio.NewWriter(f)
	if _, err := w.Write(newBuf.Bytes()); err != nil {
		return errors.Annotate(err, "write to output")
	}

	if err := w.Flush(); err != nil {
		return errors.Annotate(err, "flush out writer")
	}

	if err := f.Close(); err != nil {
		return errors.Annotate(err, "outfile close")
	}

	printInfo("process check & run")

	pipes := make([]pipe.Pipe, 0)
	pipes = appendCommandPipe(set.Check, pipes)
	pipes = appendCommandPipe(set.Run, pipes)

	script := pipe.Script(pipes...)
	if output, err := pipe.CombinedOutput(script); err != nil {
		printInfo(string(output))
		return errors.Annotate(err, "check & run")
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////
func processTemplates(cnf *Config) error {

	apiURL := fmt.Sprintf("%s%s", cnf.Host, cnf.Prefix)
	meta := metadata.NewClient(apiURL)

	printInfo("connect rancher metadata url: %q", apiURL)

	//expand template paths
	printDebug("expand template paths")
	for idx, set := range cnf.Sets {
		if !path.IsAbs(set.TemplatePath) {
			cnf.Sets[idx].TemplatePath = path.Join(DEFAULT_TEMPLATE_DIR, set.TemplatePath)
		}
	}

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
			scratch.Reset()
			if err := processTemplateSet(meta, set); err != nil {
				return errors.Annotate(err, "process template set")
			}
		}

		time.Sleep(time.Duration(cnf.Repeat) * time.Second)
	}

	return nil
}
