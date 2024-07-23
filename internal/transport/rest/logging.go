package rest

import "github.com/sirupsen/logrus"

func logFields(handler, problem string) logrus.Fields {
	return logrus.Fields{
		"handler": handler,
		"problem": problem,
	}
}
func logError(handler, problem string, err error) {
	logrus.WithFields(logFields(handler, problem)).Error(err)
}
