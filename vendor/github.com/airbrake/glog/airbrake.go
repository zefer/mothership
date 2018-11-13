package glog

import (
	"fmt"
	"net/http"

	"github.com/airbrake/gobrake"
)

// Gobrake is an instance of Airbrake Go Notifier that is used to send
// logs to Airbrake.
var Gobrake *gobrake.Notifier

var gobrakeSeverity = errorLog

func SetGobrakeNotifier(notifier *gobrake.Notifier) {
	Gobrake = notifier
}

// SetGobrakeSeverity sets minimum log severity that will be sent to Airbrake.
//
// Valid names are "INFO", "WARNING", "ERROR", and "FATAL".  If the name is not
// recognized, "ERROR" severity is used.
func SetGobrakeSeverity(severity string) {
	v, ok := severityByName(severity)
	if ok {
		gobrakeSeverity = v
	}
}

type requester interface {
	Request() *http.Request
}

// Implemented by context.Context
type valuer interface {
	Value(key interface{}) interface{}
}

func notifyAirbrake(depth int, s severity, format string, args ...interface{}) {
	if Gobrake == nil {
		return
	}
	if s < gobrakeSeverity {
		return
	}

	var msg string
	if format != "" {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = fmt.Sprint(args...)
	}

	var theErr error
	var req *http.Request
	var values valuer
	for _, arg := range args {
		if v, ok := arg.(error); ok {
			theErr = v
		}
		if v, ok := arg.(*http.Request); ok {
			req = v
		}
		if v, ok := arg.(requester); ok {
			req = v.Request()
		}
		if v, ok := arg.(valuer); ok {
			values = v
		}
	}

	var notice *gobrake.Notice
	if theErr != nil {
		notice = Gobrake.Notice(theErr, req, depth)
		notice.Errors[0].Message = msg
	} else {
		notice = Gobrake.Notice(msg, req, depth)
	}
	notice.Context["severity"] = severityName[s]

	if values != nil {
		route, _ := values.Value("route").(string)
		if route != "" {
			notice.Context["route"] = route
		}
	}

	Gobrake.SendNoticeAsync(notice)
}
