package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (e *Engine) mainLog(format string, v ...interface{}) {
	e.logger.main()(format, v...)
}

func (e *Engine) buildLog(format string, v ...interface{}) {
	e.logger.build()(format, v...)
}

func (e *Engine) runnerLog(format string, v ...interface{}) {
	e.logger.runner()(format, v...)
}

func (e *Engine) watcherLog(format string, v ...interface{}) {
	e.logger.watcher()(format, v...)
}

func (e *Engine) appLog(format string, v ...interface{}) {
	e.logger.app()(format, v...)
}

func (e *Engine) isTmpDir(path string) bool {
	return e.config.fullPath(path) == e.config.tmpPath()
}

func isHiddenDirectory(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

func cleanPath(path string) string {
	return strings.TrimSuffix(strings.TrimSpace(path), "/")
}

func (e *Engine) isExcludeDir(path string) bool {
	rp := e.config.rel(path)
	for _, d := range e.config.Build.ExcludeDir {
		if cleanPath(rp) == d {
			return true
		}
	}
	return false
}

func (e *Engine) isIncludeExt(path string) bool {
	ext := filepath.Ext(path)
	for _, v := range e.config.Build.IncludeExt {
		if ext == "."+strings.TrimSpace(v) {
			return true
		}
	}
	return false
}

func (e *Engine) writeBuildErrorLog(msg string) error {
	var err error
	f, err := os.OpenFile(e.config.buildLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err = f.Write([]byte(msg)); err != nil {
		return err
	}
	return f.Close()
}

func (e *Engine) withLock(f func()) {
	e.mu.Lock()
	f()
	e.mu.Unlock()
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home := os.Getenv("HOME")
		return home + path[1:], nil
	}
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if path == "." {
		return wd, nil
	}
	if strings.HasPrefix(path, "./") {
		return wd + path[2:], nil
	}
	return path, nil
}

func killCmd(cmd *exec.Cmd) (int, error) {
	pid := cmd.Process.Pid
	return pid, cmd.Process.Kill()
}
