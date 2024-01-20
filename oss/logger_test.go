package oss

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	l := NewLogger(0, nil)
	l.Debugf("%s", "123")
	l.Infof("%s", "123")
	l.Warnf("%s", "123")
	l.Errorf("%s", "123")

	buff := bytes.NewBuffer(nil)
	l = NewLogger(LogDebug, LogPrinterFunc(func(a ...any) {
		fmt.Fprint(buff, a...)
	}))
	buff.Reset()
	l.Debugf("%s", "123")
	assert.Equal(t, "DEBUG 123", buff.String())

	buff.Reset()
	l.Infof("%s", "123")
	assert.Equal(t, "INFO 123", buff.String())

	buff.Reset()
	l.Warnf("%s", "123")
	assert.Equal(t, "WARNING 123", buff.String())

	buff.Reset()
	l.Errorf("%s", "123")
	assert.Equal(t, "ERROR 123", buff.String())

	l = NewLogger(LogInfo, LogPrinterFunc(func(a ...any) {
		fmt.Fprint(buff, a...)
	}))
	buff.Reset()
	l.Debugf("%s", "123")
	assert.Equal(t, "", buff.String())

	buff.Reset()
	l.Infof("%s", "123")
	assert.Equal(t, "INFO 123", buff.String())

	buff.Reset()
	l.Warnf("%s", "123")
	assert.Equal(t, "WARNING 123", buff.String())

	buff.Reset()
	l.Errorf("%s", "123")
	assert.Equal(t, "ERROR 123", buff.String())

	l = NewLogger(LogError, LogPrinterFunc(func(a ...any) {
		fmt.Fprint(buff, a...)
	}))
	buff.Reset()
	l.Debugf("%s", "123")
	assert.Equal(t, "", buff.String())

	buff.Reset()
	l.Infof("%s", "123")
	assert.Equal(t, "", buff.String())

	buff.Reset()
	l.Warnf("%s", "123")
	assert.Equal(t, "", buff.String())

	buff.Reset()
	l.Errorf("%s", "123")
	assert.Equal(t, "ERROR 123", buff.String())

	l = NewLogger(LogDebug, nil)
	l.Debugf("%s", "123")
	l.Infof("%s", "123")
	l.Warnf("%s", "123")
	l.Errorf("%s", "123")
}
