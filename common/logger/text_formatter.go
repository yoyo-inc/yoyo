package logger

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const DateTimeLayout = "2006-01-02 15:04:05.000"

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
	if buf == nil {
		return nil, errors.New("entry buffer nil")
	}

	// timestamp
	timestamp := time.Now().Format(DateTimeLayout)
	buf.WriteString("[")
	buf.WriteString(timestamp)
	buf.WriteString("] ")

	// level
	buf.WriteString("[")
	buf.WriteString(strings.ToUpper(entry.Level.String()))
	buf.WriteString("] ")

	// service name
	buf.WriteString("[")
	buf.WriteString(f.Service)
	buf.WriteString("] ")

	// file information
	if entry.HasCaller() && entry.Level != logrus.InfoLevel {
		buf.WriteString(fmt.Sprintf("%s:%d ", entry.Caller.File, entry.Caller.Line))
	}

	entry.Message = strings.TrimSuffix(entry.Message, "\n")
	// message
	buf.WriteString("-")
	buf.WriteString(" ")
	buf.WriteString(entry.Message)
	buf.WriteString("\n")

	return buf.Bytes(), nil
}
