package controller

import (
	"fmt"
	"go.uber.org/zap"
	"path"
	"runtime"
	"terraria-run/internal/agent"
	"terraria-run/internal/agent/handler"
	"terraria-run/internal/controller/mod"
	"terraria-run/internal/controller/serverconfig"
	"time"
)

const (
	StatusActive     = "ACTIVE"
	StatusInactive   = "INACTIVE"
	StatusStarting   = "STARTING"
	StatusStopping   = "STOPPING"
	StatusRestarting = "RESTARTING"
	StatusUpdating   = "UPDATING"
)

// TODO: Use lock

var log = zap.S()
var status = StatusInactive
var agt *agent.Agent

func Status() string {
	return status
}

func Agent() *agent.Agent {
	return agt
}

func Start() error {
	if err := serverconfig.NewHandler().Deploy(); err != nil {
		return err
	}
	if err := mod.NewHandler().Deploy(); err != nil {
		return err
	}
	a := agent.NewAgent("")
	if err := a.Start(); err != nil {
		return err
	}
	agt = a
	changeStatus(StatusStarting)
	go watchAgentStaringComplete(a)
	return nil
}

func Stop() error {
	if agt == nil {
		return fmt.Errorf("agent is nil, stop failed")
	}
	changeStatus(StatusStopping)
	go func() {
		if err := agt.Stop(); err != nil {
			panic(err)
		}
		agt = nil
		changeStatus(StatusInactive)
	}()
	return nil
}

func watchAgentStaringComplete(a *agent.Agent) {
	defer func() {
		changeStatus(StatusActive)
	}()
	for !a.Listener.Report.Ready() {
		time.Sleep(500 * time.Millisecond)
	}
	log.Info("Report is ready, begin watch")
	for {
		// TODO: Think about process dead
		time.Sleep(500 * time.Millisecond)
		events, err := a.Listener.Report.GetEvents()
		if err != nil {
			log.Error("Get report events failed: ", err)
			return
		}
		for i := range events {
			if events[i].Type == handler.TypeServerActive {
				return
			}
		}
	}
}

func changeStatus(newStatus string) {
	_, filePath, line, _ := runtime.Caller(1)
	_, filename := path.Split(filePath)
	log.Warnf("[%s:%d] Status change from %s to %s", filename, line, status, newStatus)
	status = newStatus
}
