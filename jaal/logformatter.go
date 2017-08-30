package jaal

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct {
	indent string
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var serialized []byte
	var err error
	if len(f.indent) > 0 {
		serialized, err = json.MarshalIndent(entry.Data["data"], "", f.indent)
	} else {
		serialized, err = json.Marshal(entry.Data["data"])
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return serialized, nil
}
