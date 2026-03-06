package runner

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"
)

var serverProcess *exec.Cmd

func Run(buildCmd, execCmd string) {

	stop()

	slog.Info("building")

	build := exec.Command("sh", "-c", buildCmd)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	if err := build.Run(); err != nil {
		slog.Error("build failed", "err", err)
		return
	}

	slog.Info("starting server")

	serverProcess = exec.Command("sh", "-c", execCmd)
	serverProcess.Stdout = os.Stdout
	serverProcess.Stderr = os.Stderr
	serverProcess.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	err := serverProcess.Start()
	if err != nil {
		slog.Error("start failed", "err", err)
	}
}

func stop() {

	if serverProcess != nil && serverProcess.Process != nil {

		slog.Info("stopping server")

		syscall.Kill(-serverProcess.Process.Pid, syscall.SIGKILL)

		serverProcess.Wait()

		serverProcess = nil
	}
}
