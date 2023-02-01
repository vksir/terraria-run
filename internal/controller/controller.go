package controller

import (
	"fmt"
	"go.uber.org/zap"
	"terraria-run/internal/agent"
	"terraria-run/internal/controller/mod"
	"terraria-run/internal/controller/serverconfig"
)

const (
	StatusActive     = "ACTIVE"
	StatusInactive   = "INACTIVE"
	StatusStarting   = "STARTING"
	StatusStopping   = "STOPPING"
	StatusRestarting = "RESTARTING"
	StatusUpdating   = "UPDATING"
)

var log = zap.S()
var status = StatusInactive
var agt *agent.Agent

func Status() string {
	return status
}

func Start() error {
	if status != StatusInactive {
		return fmt.Errorf("status is %s, cannot start", status)
	}
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
	a = agt
	status = StatusStarting
	// TODO: watch starting complete
	return nil
}

func Stop() error {
	if status != StatusActive && status != StatusStarting {
		return fmt.Errorf("status is %s, cannot stop", status)
	}
	if agt == nil {
		return fmt.Errorf("agent is nil, stop failed")
	}
	status = StatusStopping
	go func() {
		if err := agt.Stop(); err != nil {
			panic(err)
		}
		agt = nil
		status = StatusInactive
	}()
	return nil
}
