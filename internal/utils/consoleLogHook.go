package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type ConsoleLogHook struct {
}

func (hook *ConsoleLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ConsoleLogHook) Fire(entry *logrus.Entry) error {

	consoleFormat := &TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	rawByte, err := consoleFormat.Format(entry)
	msg := string(rawByte)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Print(msg)
	return nil

}
