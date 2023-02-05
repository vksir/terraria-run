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
	listener         *listener.Listener
	process          *process.Process
	record           *handler.Record
	report           *handler.Report
	cut              *handler.Cut
	ctx              context.Context
	stopListenOutput context.CancelFunc
}

func NewAgent(name string, report *handler.Report) *Agent {
	a := Agent{
		Name:     name,
		process:  process.NewProcess(""),
		listener: listener.NewListener(),
		record:   handler.NewRecord(),
		report:   report,
		cut:      handler.NewCut(),
	}
	a.ctx, a.stopListenOutput = context.WithCancel(context.Background())
	a.listener.RegisterHandler(a.record)
	a.listener.RegisterHandler(a.report)
	a.listener.RegisterHandler(a.cut)
	return &a
}

func (a *Agent) Start() error {
	log.Info("Start agent")
	if err := a.process.Start(); err != nil {
		return err
	}
	if err := a.record.Start(a.ctx); err != nil {
		return err
	}
	if err := a.report.Start(a.ctx); err != nil {
		return err
	}
	if err := a.cut.Start(a.ctx); err != nil {
		return err
	}
	if err := a.listener.Start(a.ctx, io.MultiReader(a.process.Stdout, a.process.Stderr)); err != nil {
		return err
	}
	return nil
}

func (a *Agent) Stop() error {
	log.Info("Stop agent")
	defer a.stopListenOutput()
	return a.process.Stop()
}

func (a *Agent) RunCmd(cmd string) (string, error) {
	cmd = strings.TrimRight(cmd, "\n")
	if err := a.cut.BeginCut(); err != nil {
		return "", err
	}
	if _, err := a.process.Stdin.Write([]byte(fmt.Sprintf("[BEGIN_CMD]\n%s\n[END_CMD]\n", cmd))); err != nil {
		_, _ = a.cut.StopCut()
		return "", err
	}
	time.Sleep(500 * time.Millisecond)
	rawOut, err := a.cut.StopCut()
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
