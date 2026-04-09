package app

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func runCommand(command string) tea.Cmd {
	return func() tea.Msg {
		cmd := buildShellCommand(command)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		return commandFinishedMsg{
			command: command,
			output:  strings.TrimRight(stdout.String(), "\r\n"),
			errText: strings.TrimRight(stderr.String(), "\r\n"),
			err:     err,
		}
	}
}

func buildShellCommand(command string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", "-NoLogo", "-NoProfile", "-Command", command)
	}
	return exec.Command("sh", "-lc", command)
}

func currentShellLabel() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	return "sh"
}
