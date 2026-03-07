package runner

import (
	"fmt"
	"os/exec"
	"syscall"
	"testing"
)

func TestStop_KillsProcess(t *testing.T) {
	// start a long-running process
	serverProcess = exec.Command("sleep", "100")
	serverProcess.SysProcAttr = pgidAttr()
	if err := serverProcess.Start(); err != nil {
		t.Fatal(err)
	}

	pid := serverProcess.Process.Pid
	Stop()

	// process should be dead
	err := exec.Command("kill", "-0", fmt.Sprintf("%d", pid)).Run()
	if err == nil {
		t.Error("process still running after Stop()")
	}
}
func pgidAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}
