package jaal

import (
	"github.com/huguesalary/slack-go"
	"github.com/sirupsen/logrus"
)

type SlackNotifier struct {
}

var allLevels = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

func (SlackNotifier) Levels() []logrus.Level {
	return allLevels
}

func (SlackNotifier) Fire(entry *logrus.Entry) error {
	event := entry.Data["data"].(*Event)
	c := slack.NewClient("https://hooks.slack.com/services/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")
	msg := &slack.Message{
		Channel: "#general",
		Text:    event.Summary,
	}
	return c.SendMessage(msg)
}
