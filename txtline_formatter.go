package log

import (
	"bytes"
)

type TxtLineFormatter struct{}

func (t *TxtLineFormatter) Format(buffer *bytes.Buffer, l *LogInstance) (*bytes.Buffer, error) {
	if l.Flag&FlagTime > 0 {
		buffer.WriteString(l.Time)
		buffer.WriteString(FieldSplit)
	}
	if l.Flag&FlagLevel > 0 {
		buffer.WriteString(l.Level)
		buffer.WriteString(FieldSplit)
	}
	if l.Flag&FlagFile > 0 {
		buffer.WriteString(l.File)
		buffer.WriteString(FieldSplit)
	}

	buffer.WriteString(l.Msg)
	if l.Msg[len(l.Msg)-1] != '\n' {
		buffer.WriteByte('\n')
	}

	return buffer, nil
}
