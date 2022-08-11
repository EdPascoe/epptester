package log

// Just exists to configure easy logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

func init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05", // the "time" field configuratiom
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// this function is required when you want to introduce your custom format.
			// In my case I wanted file and line to look like this `file="engine.go:141`
			// but f.File provides a full path along with the file name.
			// So in `formatFilePath()` function I just trimmet everything before the file name
			// and added a line number in the end
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
