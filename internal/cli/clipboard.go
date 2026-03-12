package cli

import (
	"errors"
	"os/exec"
	"runtime"
)

var ErrNoClipboard = errors.New("clipboard tool not available")

func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else {
			return ErrNoClipboard
		}
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err = stdin.Write([]byte(text))
	if closeErr := stdin.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if waitErr := cmd.Wait(); waitErr != nil && err == nil {
		err = waitErr
	}
	return err
}
