package agent

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"terraria-run/internal/agent/output"
	"terraria-run/internal/common/constant"
	"time"
)

var log = zap.S()

type Agent struct {
	Name     string
	Listener *output.Listener

	cmd              *exec.Cmd
	stdin            io.Writer
	stdout           io.Reader
	stderr           io.Reader
	ctx              context.Context
	stopListenOutput context.CancelFunc
}

func NewAgent(name string) *Agent {
	cmd := exec.Command(constant.DotnetPath, "tModLoader.dll", "-server", "-config", constant.ServerConfigPath)
	cmd.Dir = filepath.Join(constant.InstallDir, "tModLoader")
	a := Agent{
		Name:     name,
		cmd:      cmd,
		Listener: output.NewListener(),
	}
	a.ctx, a.stopListenOutput = context.WithCancel(context.Background())
	log.Infof("New agent: %s", strings.Join(a.cmd.Args, " "))
	return &a
}

func (a *Agent) Start() (err error) {
	log.Info("Start agent")
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
		err := a.Listener.Run(a.ctx, io.MultiReader(a.stdout, a.stderr))
		if err != nil {
			log.Error("Listener stopped", err)
		}
	}()
	err = a.cmd.Start()
	if err != nil {
		return err
	}
	return err
}

func (a *Agent) Stop() error {
	log.Info("Stop agent")
	defer a.stopListenOutput()
	return a.stopProcess()
}

func (a *Agent) RunCmd(cmd string) (string, error) {
	// TODO: Read output
	cmd = fmt.Sprintf("%s\n", strings.TrimRight(cmd, "\n"))
	_, err := a.stdin.Write([]byte(cmd))
	if err != nil {
		return "", err
	}
	return "", nil
}

func (a *Agent) stopProcess() error {
	_, err := a.RunCmd("exit")
	if err != nil {
		return err
	}
	t := time.Now()
	for time.Since(t).Seconds() < 15 {
		if a.cmd.ProcessState != nil {
			log.Info("Process exit normally")
			return nil
		}
		time.Sleep(time.Second)
	}
	err = a.cmd.Process.Kill()
	if err != nil {
		return err
	}
	log.Error("Process been killed", nil)
	return nil
}
