package logger

import (
	"github.com/Sirupsen/logrus"
)

var (
	Server *logrus.Entry
	Client *logrus.Entry
)

func SetJSONFormatter() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// TODO(cesar0094): add processId/systemId as field so the logs can be filtered by
// the specific server/client
func init() {
	Server = logrus.WithFields(logrus.Fields{
		"source": "server",
	})
	Client = logrus.WithFields(logrus.Fields{
		"source": "client",
	})
}
