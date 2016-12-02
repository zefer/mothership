package glog

import (
	"bytes"
	"fmt"

	golog "github.com/golang/glog"
)

const logDepth = 4

// logBridge provides the Write method that enables CopyStandardLogTo to connect
// Go's standard logs to the logs provided by this package.
type logBridge severity

// Write parses the standard logging line and passes its components to the
// logger for severity(lb).
func (lb logBridge) Write(b []byte) (n int, err error) {
	var text string
	// Split "d.go:23: message" into "d.go", "23", and "message".
	if parts := bytes.SplitN(b, []byte{':'}, 3); len(parts) != 3 || len(parts[0]) < 1 || len(parts[2]) < 1 {
		text = fmt.Sprintf("bad log format: %s", b)
	} else {
		text = string(parts[2][1:]) // skip leading space
	}

	switch severity(lb) {
	case infoLog:
		golog.InfoDepth(logDepth, text)
	case warningLog:
		golog.WarningDepth(logDepth, text)
	case errorLog:
		golog.ErrorDepth(logDepth, text)
	case fatalLog:
		golog.FatalDepth(logDepth, text)
	}
	notifyAirbrake(logDepth+1, severity(lb), text)

	return len(b), nil
}
