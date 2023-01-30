package agent

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"terraria-run/internal/common/constant"
	"time"
)

type Agent struct {
	Name             string
	cmd              *exec.Cmd
	stdin            io.Writer
	stdout           io.Reader
	stderr           io.Reader
	ctx              context.Context
	stopListenOutput context.CancelFunc
}

func NewAgent(name string) *Agent {
	a := Agent{Name: name}
	a.cmd = exec.Command(constant.DotnetPath, "tModLoader.dll", "-server", "-config", constant.ServerConfigPath)
	a.cmd.Dir = filepath.Join(constant.InstallDir, "tModLoader")
	a.ctx, a.stopListenOutput = context.WithCancel(context.Background())
	return &a
}

func (a *Agent) Start() (err error) {
	slog.Info("Start agent")
	a.stdin, err = a.cmd.StdinPipe()
	if err != nil {
		return err
	}
	a.stdout, err = a.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	a.stderr, err = a.cmd.StderrPipe()
	if err != nil {
		return err
	}
	go func() {
		err := a.listenOutput()
		if err != nil {
			slog.Error("Stop listen output", err)
		}
	}()
	err = a.cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (a *Agent) Stop() error {
	slog.Info("Stop agent")
	defer a.stopListenOutput()
	return a.stopProcess()
}

func (a *Agent) RunCmd(cmd string) (string, error) {
	cmd = fmt.Sprintf("%s\n", strings.TrimRight(cmd, "\n"))
	_, err := a.stdin.Write([]byte(cmd))
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *Agent) listenOutput() error {
	slog.Info("Begin listen output")
	w, err := os.OpenFile(constant.TModLoaderLogPath, os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		err := w.Close()
		if err != nil {
			slog.Error("Close file failed", err)
		}
	}(w)
	r := io.MultiReader(a.stdout, a.stderr)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		_, err = w.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return err
		}
		select {
		case <-a.ctx.Done():
			return nil
		default:
		}
	}
	return scanner.Err()
}

func (a *Agent) stopProcess() error {
	_, err := a.RunCmd("exit")
	if err != nil {
		return err
	}
	t := time.Now()
	for time.Since(t).Seconds() < 15 {
		if a.cmd.ProcessState != nil {
			slog.Info("Process exit normally")
			return nil
		}
		time.Sleep(time.Second)
	}
	err = a.cmd.Process.Kill()
	if err != nil {
		return err
	}
	slog.Error("Process been killed", nil)
	return nil
}
