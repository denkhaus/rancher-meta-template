package logging

import (
	log "github.com/Sirupsen/logrus"
	"github.com/juju/errors"
)

//////////////////////////////////////////////////////////////////////////////
func SetLogLevel(level string) error {
	lv, err := log.ParseLevel(level)
	if err != nil {
		return errors.Annotate(err, "parse log level")
	}

	log.SetLevel(lv)
	return nil
}

//////////////////////////////////////////////////////////////////////////////
func PrintTmplError(format string, args ...interface{}) {
	log.WithField("ctx", "template").Errorf(format, args...)
}

//////////////////////////////////////////////////////////////////////////////
func PrintError(err error) {
	log.Errorf(errors.ErrorStack(err))
}

//////////////////////////////////////////////////////////////////////////////
func PrintInfo(format string, args ...interface{}) {
	log.Infof(format, args...)
}

//////////////////////////////////////////////////////////////////////////////
func PrintDebug(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

//////////////////////////////////////////////////////////////////////////////
func PrintWarning(format string, args ...interface{}) {
	log.Warningf(format, args...)
}
