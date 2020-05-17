package log

import (
	"os"

	"github.com/op/go-logging"
)

const SystemLogModule = "System"

func init() {
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backend2Formatter)
}

var log = logging.MustGetLogger(SystemLogModule)

var (
	Error    = log.Error
	Errorf   = log.Errorf
	Info     = log.Info
	Infof    = log.Infof
	Warning  = log.Warning
	Warningf = log.Warningf
	Fatal    = log.Fatal
	Fatalf   = log.Fatalf
	Debug    = log.Debug
	Debugf   = log.Debugf

	IsDebug = log.IsEnabledFor(logging.DEBUG)
)
