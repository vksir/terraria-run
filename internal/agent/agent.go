package agent

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"regexp"
	"strings"
	"terraria-run/internal/agent/handler"
	"terraria-run/internal/agent/listener"
	"terraria-run/internal/agent/process"
	"time"
)

var log = zap.S()

type Agent struct {
	Name             string
	Listener         *listener.Listener
	Process          *process.Process
	Record           *handler.Record
	Report           *handler.Report
	Cut              *handler.Cut
	ctx              context.Context
	stopListenOutput context.CancelFunc
}

func NewAgent(name string) *Agent {
	a := Agent{
		Name:     name,
		Process:  process.NewProcess(""),
		Listener: listener.NewListener(),
		Record:   handler.NewRecord(),
		Report:   handler.NewReport(),
		Cut:      handler.NewCut(),
	}
	a.ctx, a.stopListenOutput = context.WithCancel(context.Background())
	a.Listener.RegisterHandler(a.Record)
	a.Listener.RegisterHandler(a.Report)
	a.Listener.RegisterHandler(a.Cut)
	return &a
}

func (a *Agent) Start() error {
	log.Info("Start agent")
	if err := a.Process.Start(); err != nil {
		return err
	}
	if err := a.Record.Start(a.ctx); err != nil {
		return err
	}
	if err := a.Report.Start(a.ctx); err != nil {
		return err
	}
	if err := a.Cut.Start(a.ctx); err != nil {
		return err
	}
	if err := a.Listener.Start(a.ctx, io.MultiReader(a.Process.Stdout, a.Process.Stderr)); err != nil {
		return err
	}
	return nil
}

func (a *Agent) Stop() error {
	log.Info("Stop agent")
	defer a.stopListenOutput()
	return a.Process.Stop()
}

func (a *Agent) RunCmd(cmd string) (string, error) {
	cmd = strings.TrimRight(cmd, "\n")
	if err := a.Cut.BeginCut(); err != nil {
		return "", err
	}
	if _, err := a.Process.Stdin.Write([]byte(fmt.Sprintf("[BEGIN_CMD]\n%s\n[END_CMD]\n", cmd))); err != nil {
		_, _ = a.Cut.StopCut()
		return "", err
	}
	time.Sleep(500 * time.Millisecond)
	rawOut, err := a.Cut.StopCut()
	if err != nil {
		return "", err
	}
	pattern := regexp.MustCompile(`(?ms)Invalid command.*^(.*)$\n.*Invalid command`)
	res := pattern.FindStringSubmatch(rawOut)
	if len(res) == 0 {
		return "", fmt.Errorf("regex cmd out failed: %s", rawOut)
	}
	log.Infof("Run cmd %s success: %s", cmd, res[1])
	return res[1], nil
}
