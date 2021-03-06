package main

import (
	"strings"

	"github.com/denkhaus/rancher-meta-template/logging"
	"github.com/juju/errors"

	"github.com/rancher/go-rancher-metadata/metadata"
)

//////////////////////////////////////////////////////////////////////////////////
type ServiceWrap struct {
	metadata.Service
}

//////////////////////////////////////////////////////////////////////////////////
type ContainerWrap struct {
	metadata.Container
}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) PortInternal(idx int) (string, error) {
	if len(p.Ports) == 0 {
		return "", nil
	}

	if idx >= len(p.Ports) {
		return "", errors.Errorf("PortInternal:: index %d out of range", idx)
	}

	part := p.Ports[idx]
	if part == "" {
		return part, nil
	}

	logging.Debug(part)
	items := strings.Split(part, ":")
	part = items[len(items)-1]
	if strings.Contains(part, "/") {
		part = strings.Split(part, "/")[0]
	}

	return part, nil

}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) LabelByKey(key string) (string, error) {
	if val, ok := p.Labels[key]; ok {
		return val, nil
	}

	logging.Debug("LabelByKey:: container: %s: key %q is not present", p.Name, key)
	return "", nil
}

//////////////////////////////////////////////////////////////////////////////////
type HostWrap struct {
	metadata.Host
}
