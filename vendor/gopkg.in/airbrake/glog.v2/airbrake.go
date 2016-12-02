package glog

import (
	"fmt"
	"net/http"

	"gopkg.in/airbrake/gobrake.v2"
)

// Gobrake is an instance of Airbrake Go Notifier that is used to send
// logs to Airbrake.
var Gobrake *gobrake.Notifier

// Minimum log severity that will be sent to Airbrake.
//
// Valid names are "INFO", "WARNING", "ERROR", and "FATAL".  If the name is not
// recognized, "ERROR" severity is used.
//
// TODO: replace with SetGobrakeSeverity
var GobrakeSeverity = "ERROR"

type requester interface {
	Request() *http.Request
}

func notifyAirbrake(depth int, s severity, format string, args ...interface{}) {
	if Gobrake == nil {
		return
	}

	severity, ok := severityByName(GobrakeSeverity)
	if !ok {
		severity = errorLog
	}
	if s < severity {
		return
	}

	var msg string
	if format != "" {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = fmt.Sprint(args...)
	}

	var req *http.Request
	for _, arg := range args {
		if v, ok := arg.(requester); ok {
			req = v.Request()
			break
		}
	}

	for _, arg := range args {
		err, ok := arg.(error)
		if !ok {
			continue
		}

		notice := Gobrake.Notice(err, req, depth)
		notice.Errors[0].Message = msg
		Gobrake.SendNoticeAsync(notice)
		return
	}

	notice := Gobrake.Notice(msg, req, depth)
	Gobrake.SendNoticeAsync(notice)
}
