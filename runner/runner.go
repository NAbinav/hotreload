package runner

import (
	"log/slog"
	"os"
	"os/exec"
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
	serverProcess.Start()
}

func stop() {
	if serverProcess != nil && serverProcess.Process != nil {
		slog.Info("stopping server")
		serverProcess.Process.Kill()
		serverProcess.Wait()
		serverProcess = nil
	}
}
