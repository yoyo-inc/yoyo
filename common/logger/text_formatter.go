package logger

import (
	"bytes"
	"time"

	"github.com/sirupsen/logrus"
)

const BASIC_TIME_LAYOUT = "2006/01/02 15:04:05.000"

type TextFormatter struct {
	Service string
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	// timestamp
	timestamp := time.Now().Format(BASIC_TIME_LAYOUT)
	buf.WriteString("[")
	buf.WriteString(timestamp)
	buf.WriteString("]")
	buf.WriteString(" ")

	// level
	buf.WriteString("[")
	buf.WriteString(entry.Level.String())
	buf.WriteString("]")
	buf.WriteString(" ")

	// service name
	buf.WriteString("[")
	buf.WriteString(f.Service)
	buf.WriteString("]")

	// message
	buf.WriteString(" ")
	buf.WriteString("-")
	buf.WriteString(" ")
	buf.WriteString(entry.Message)
	return buf.Bytes(), nil
}
