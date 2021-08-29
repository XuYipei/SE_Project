package tools

import (
	"os/exec"
)

func SysPwsh(opt string) error {
	// log.print(opt)

	cmd := exec.Command("pwsh", "-c", opt)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
