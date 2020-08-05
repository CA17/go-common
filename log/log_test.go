package log

import (
	llog "log"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

func TestFileLog(t *testing.T) {

	_log := llog.New(&lumberjack.Logger{
		Filename:   "/var/zpm/logs/test.log.txt",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}, "p1", 0)

	_log.Println("asdasdasd")

	_log2 := llog.New(&lumberjack.Logger{
		Filename:   "/var/zpm/logs/test.log2.txt",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}, "p1:", 0)

	_log2.Println("asdasdasd")
	_log2.Println("34563465364")

}
