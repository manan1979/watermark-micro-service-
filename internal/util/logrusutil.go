package util

import (
	"github.com/sirupsen/logrus"
)

type DefaultFieldFormatter struct {
	WrappedFormatter logrus.Formatter
	DefaultFields    logrus.Fields
	PrintLineNumber  bool
}

func Init(formatter *DefaultFieldFormatter) {
	if formatter == nil {
		return
	}

	if formatter.WrappedFormatter == nil {

		formatter.WrappedFormatter = &logrus.JSONFormatter{}
	}

	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(formatter.PrintLineNumber)

}

func ComponentInit(component string) {
	Init(

		&DefaultFieldFormatter{
			PrintLineNumber: true,
			DefaultFields:   logrus.Fields{"component": component},
		},
	)
}

func (f *DefaultFieldFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields, len(entry.Data)+len(f.DefaultFields))

	for k, v := range f.DefaultFields {
		data[k] = v
	}

	for k, v := range entry.Data {
		data[k] = v
	}

	return f.WrappedFormatter.Format(&logrus.Entry{
		Logger:  entry.Logger,
		Data:    data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
		Caller:  entry.Caller,
	})
}
