package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

//////////////////////////////////////////////////////////////////////////////
func setLogLevel(level string) error {
	lv, err := log.ParseLevel(level)
	if err != nil {
		return errors.Annotate(err, "parse log level")
	}

	log.SetLevel(lv)
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func printTmplError(format string, args ...interface{}) {
	log.Errorf("rancher-mt:template: %s", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func printError(err error) {
	log.Errorf("rancher-mt: %s", errors.ErrorStack(err))
}

//////////////////////////////////////////////////////////////////////////////
func printInfo(format string, args ...interface{}) {
	log.Infof("rancher-mt: %s", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func printDebug(format string, args ...interface{}) {
	log.Debugf("rancher-mt: %s", fmt.Sprintf(format, args...))
}

//////////////////////////////////////////////////////////////////////////////
func printWarning(format string, args ...interface{}) {
	log.Warningf("rancher-mt: %s", fmt.Sprintf(format, args...))
}
