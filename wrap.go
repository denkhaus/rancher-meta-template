package main

import (
	"strings"

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
func (p ContainerWrap) portSelect(idx int) string {
	part := p.Ports[idx]
	if part == "" {
		return part
	}

	items := strings.Split(part, ":")
	if len(items) == 2 {
		part = items[idx]
	} else {
		part = items[0]
	}

	if strings.Contains(part, "/") {
		part = strings.Split(part, "/")[0]
	}

	return part
}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) PortInternal(idx int) (string, error) {
	if len(p.Ports) > idx {
		printTmplError("PortInternal: index %d out of range", idx)
		return "", errors.New("PortInternal: index %d out of range")
	}

	return p.portSelect(2), nil
}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) PortExternal(idx int) (string, error) {
	if len(p.Ports) > idx {
		printTmplError("PortExternal: index %d out of range", idx)
		return "", errors.New("PortExternal: index %d out of range")
	}

	return p.portSelect(1), nil
}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) LabelByKey(key string) (string, error) {
	if val, ok := p.Labels[key]; ok {
		return val, nil
	}

	return "", errors.Errorf("LabelByKey: label %s not found", key)
}

//////////////////////////////////////////////////////////////////////////////////
type HostWrap struct {
	metadata.Host
}
