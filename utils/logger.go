package utils

import "github.com/sirupsen/logrus"

func GetLogLevel(level uint32) []logrus.Level {
	var levels []logrus.Level
	for _, allLevel := range logrus.AllLevels {
		if allLevel > logrus.Level(level) {
			break
		}
		levels = append(levels, allLevel)
	}
	return levels
}
