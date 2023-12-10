package ml_process

import (
	"io"
	"os/exec"
	"strings"
	"syscall"
)

func ExecWithStdIn(stdin io.Reader, args ...string) (string, error) {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdin = stdin
	c.Stdout = b
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := c.Run()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func Exec(args ...string) (string, string, error) {
	c := exec.Command(args[0], args[1:]...)
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	c.Stdout = stdout
	c.Stderr = stderr
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := c.Run()
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

func Create(args ...string) (*exec.Cmd, *strings.Builder) {
	c, b := exec.Command(args[0], args[1:]...), new(strings.Builder)
	c.Stdout = b
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return c, b
}
