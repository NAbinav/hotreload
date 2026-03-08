package runner

import (
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var (
	mu            sync.Mutex
	serverProcess *exec.Cmd
	lastStartTime time.Time
)

func Run(buildCmd, execCmd string) {
	mu.Lock()
	defer mu.Unlock()

	stop()

	slog.Info("building")
	build := exec.Command("sh", "-c", buildCmd)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	if err := build.Run(); err != nil {
		slog.Error("build failed", "err", err)
		return
	}

	if time.Since(lastStartTime) < time.Second {
		slog.Warn("server crashed too quickly, skipping restart")
		return
	}

	slog.Info("starting server")
	serverProcess = exec.Command("sh", "-c", execCmd)
	serverProcess.Stdout = os.Stdout
	serverProcess.Stderr = os.Stderr
	serverProcess.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := serverProcess.Start(); err != nil {
		slog.Error("start failed", "err", err)
		return
	}
	lastStartTime = time.Now()
}

func Stop() {
	mu.Lock()
	defer mu.Unlock()
	stop()
}

func stop() { // unexported, call with lock held
	if serverProcess != nil && serverProcess.Process != nil {
		slog.Info("stopping server")
		syscall.Kill(-serverProcess.Process.Pid, syscall.SIGTERM)
		serverProcess.Wait()
		serverProcess = nil
	}
}
