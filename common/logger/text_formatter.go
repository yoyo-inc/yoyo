package logger

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const basicTimeLayout = "2006/01/02 15:04:05.000"

// TextFormatter provides
type TextFormatter struct {
	Service string
}

// Format formats entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	// timestamp
	timestamp := time.Now().Format(basicTimeLayout)
	buf.WriteString("[")
	buf.WriteString(timestamp)
	buf.WriteString(" ]")

	// level
	buf.WriteString("[")
	buf.WriteString(strings.ToUpper(entry.Level.String()))
	buf.WriteString("] ")

	// service name
	buf.WriteString("[")
	buf.WriteString(f.Service)
	buf.WriteString("] ")

	// file information
	if entry.Caller != nil {
		buf.WriteString(entry.Caller.File)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(entry.Caller.Line))
	}

	entry.Message = strings.TrimSuffix(entry.Message, "\n")
	// message
	buf.WriteString(" ")
	buf.WriteString("-")
	buf.WriteString(" ")
	buf.WriteString(entry.Message)

	buf.WriteString("\n")
	return buf.Bytes(), nil
}
