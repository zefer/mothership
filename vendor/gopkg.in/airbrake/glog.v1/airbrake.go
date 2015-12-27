package glog

import (
	"fmt"
	"net/http"

	"gopkg.in/airbrake/gobrake.v1"
)

var Gobrake *gobrake.Notifier

// Minimum log severity that will be sent to Airbrake.
var GobrakeSeverity = ErrorLog

func notifyAirbrake(s severity, format string, args ...interface{}) {
	if Gobrake == nil {
		return
	}
	if s < GobrakeSeverity {
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

	foundErr := false
	for _, arg := range args {
		err, ok := arg.(error)
		if !ok {
			continue
		}
		foundErr = true

		notice := Gobrake.Notice(err, req, 5)
		notice.Env["glog_message"] = msg
		go Gobrake.SendNotice(notice)
	}

	if !foundErr {
		notice := Gobrake.Notice(msg, req, 5)
		go Gobrake.SendNotice(notice)
	}
}

type requester interface {
	Request() *http.Request
}
