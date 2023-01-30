package main

import (
	"golang.org/x/exp/slog"
	"terraria-run/internal/agent"
	_ "terraria-run/internal/common/config"
	_ "terraria-run/internal/common/log"
	"time"
)

func main() {
	slog.Info("Start terraria run")
	a := agent.NewAgent("")
	err := a.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	err = a.Stop()
	if err != nil {
		slog.Error("Stop agent failed", err)
	}
}
