package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

func printTmplError(format string, args ...interface{}) {
	log.Errorf("rancher-meta-template:template: %s", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func printError(err error) {
	log.Errorf("rancher-meta-template: %s", errors.ErrorStack(err))
}

//////////////////////////////////////////////////////////////////////////////
func printInfo(format string, args ...interface{}) {
	log.Infof("rancher-meta-template: %s", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func printDebug(format string, args ...interface{}) {
	log.Debugf("rancher-meta-template: %s", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func printWarning(format string, args ...interface{}) {
	log.Warningf("rancher-meta-template: %s", fmt.Sprintf(format, args...))
}
