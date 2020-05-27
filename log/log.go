package log

import (
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/ca17/go-common/log/dailyrotate"
)

const ModuleSystem = "System"

var log = logging.MustGetLogger(ModuleSystem)

func SetupLog(level logging.Level, syslogaddr string, logdir string, module string) {

	var format = logging.MustStringFormatter(
		`%{color} %{time:15:04:05.000} %{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	Backends := make([]logging.Backend, 0)
	backendStderr := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backendStderr, format)
	Backends = append(Backends, backendFormatter)
	bs := _setupSyslog(level, syslogaddr, module)
	bf := _fileSyslog(level, logdir, module)

	if bs != nil {
		Backends = append(Backends, bs)
	}
	if bf != nil {
		Backends = append(Backends, bf)
	}
	logging.SetBackend(Backends...)
	logging.SetLevel(level, module)
	log = logging.MustGetLogger(module)
}

func _clearLogs(logsdir string, prefix string) {
	daydirs, err := ioutil.ReadDir(logsdir)
	if err != nil {
		log.Errorf("read day logs dir error, %s", err.Error())
		return
	}

	for _, item := range daydirs {
		if !item.IsDir() && strings.HasPrefix(item.Name(), prefix) && item.ModTime().Before(time.Now().Add(-(time.Hour * 24 * 7))) {
			fpath := filepath.Join(logsdir, item.Name())
			err = os.Remove(fpath)
			if err != nil {
				log.Errorf("remove logfile %s error", fpath)
			}
		}
	}
}

func _fileSyslog(level logging.Level, logdir string, module string) logging.LeveledBackend {
	if logdir == "N/A" {
		return nil
	}
	var format = logging.MustStringFormatter(
		`%{time:15:04:05.000} %{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
	)

	logfile, err := dailyrotate.NewFile(filepath.Join(logdir, module+"-daily-2006-01-02.log"), func(path string, didRotate bool) {
		fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
		if !didRotate {
			return
		}
		// process just closed file e.g. upload to backblaze storage for backup
		go _clearLogs(logdir, module+"-daily-")
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backendFile := logging.NewLogBackend(logfile, "", 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backend2Formatter := logging.NewBackendFormatter(backendFile, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(level, module)
	return backend1Leveled
}

func _setupSyslog(level logging.Level, syslogaddr string, module string) logging.LeveledBackend {
	var format = logging.MustStringFormatter(
		`%{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
	)
	backend, err := NewSyslogBackend("", syslogaddr, syslog.LOG_INFO)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backend2Formatter := logging.NewBackendFormatter(backend, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(level, module)
	return backend1Leveled
}

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
