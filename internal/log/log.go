package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

const pkgLen = len("github.com/ronazst/notion-ical-syncer/")

func SetupLogrus(requestId string) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&formatter{
		requestId: requestId,
		formatter: &logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
				return "", fmt.Sprintf("%s:%d", frame.File[pkgLen:], frame.Line)
			},
		},
	})
}

type formatter struct {
	requestId string
	formatter logrus.Formatter
}

func (f formatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["request_id"] = f.requestId
	return f.formatter.Format(entry)
}
