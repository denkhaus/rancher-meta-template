package main

import (
	"fmt"
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
func (p ContainerWrap) PortInternal(idx int) (string, error) {
	if len(p.Ports) == 0 {
		return "", nil
	}
	if idx >= len(p.Ports) {
		return "", errors.Errorf("PortInternal: index %d out of range", idx)
	}

	part := p.Ports[idx]
	if part == "" {
		return part
	}

	fmt.Println(part)
	items := strings.Split(part, ":")
	part = items[len(items)-1]
	if strings.Contains(part, "/") {
		part = strings.Split(part, "/")[0]
	}

	return part

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
