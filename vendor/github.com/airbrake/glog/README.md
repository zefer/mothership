# Glog

This fork of https://github.com/golang/glog provides all of glog's functionality
and adds the ability to send errors/logs to [Airbrake.io](https://airbrake.io).

## Logging

Please refer to the [glog](https://github.com/golang/glog) code & docs.

## Sending errors to Airbrake.io

A basic example of how to configure glog to send logged errors to Airbrake.io:

```go
package main

import (
	"errors"

	"github.com/airbrake/glog"
	"github.com/airbrake/gobrake"
)

var projectId int64 = 123
var apiKey string = "API_KEY"

func doSomeWork() error {
	return errors.New("hello from Go")
}

func main() {
	airbrake := gobrake.NewNotifier(projectId, apiKey)
	defer airbrake.Close()
	defer airbrake.NotifyOnPanic()

	airbrake.AddFilter(func(n *gobrake.Notice) *gobrake.Notice {
		n.Context["environment"] = "production"
		return n
	})
	glog.SetGobrakeNotifier(airbrake)

	if err := doSomeWork(); err != nil {
		glog.Errorf("doSomeWork failed: %s", err)
	}
}
```

## Configure severity

The default is to send only error logs to Airbrake.io. You can change the
severity threshold to also send lower severity logs too, such as warnings:

```go
glog.SetGobrakeSeverity("WARNING")
```
