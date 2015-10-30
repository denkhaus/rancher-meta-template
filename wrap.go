package main

import (
	"strings"

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
func (p ContainerWrap) PortInternal(idx int) string {
	if len(p.Ports) > idx {
		printTmplError("PortInternal: index %d out of range", idx)
		return ""
	}

	return p.portSelect(2)
}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) PortExternal(idx int) string {
	if len(p.Ports) > idx {
		printTmplError("PortExternal: index %d out of range", idx)
		return ""
	}

	return p.portSelect(1)
}

//////////////////////////////////////////////////////////////////////////////////
func (p ContainerWrap) LabelByKey(key string) string {
	if val, ok := p.Labels[key]; ok {
		return val
	}

	return ""
}

//////////////////////////////////////////////////////////////////////////////////
type HostWrap struct {
	metadata.Host
}
