package httpclient

import (
	"fmt"
	"github.com/andibalo/meowhasiswa-be/pkg/logger"
)

type restyLogger struct {
	l logger.Logger
}

func (r restyLogger) Errorf(format string, v ...interface{}) {
	r.l.Error(fmt.Sprintf(format, v...))
}

func (r restyLogger) Warnf(format string, v ...interface{}) {
	r.l.Warn(fmt.Sprintf(format, v...))
}

func (r restyLogger) Debugf(format string, v ...interface{}) {
	r.l.Info(fmt.Sprintf(format, v...))
}
