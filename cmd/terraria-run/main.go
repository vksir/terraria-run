package main

import (
	"go.uber.org/zap"
	"terraria-run/internal/common/config"
	_ "terraria-run/internal/common/config"
	_ "terraria-run/internal/common/log"
	"terraria-run/internal/controller"
	"terraria-run/internal/server"
)

var log = zap.S()

func main() {
	log.Info("Start terraria run")
	config.Read()

	//a
	//:= agent.NewAgent("")
	//err := a.Start()
	//if err != nil {
	//	panic(err)
	//}
	//time.Sleep(10 * time.Second)
	//err = a.Stop()
	//if err != nil {
	//	log.Error("Stop agent failed", err)
	//}

	//serverConfigHandler := serverconfig.NewHandler()
	//if err := serverConfigHandler.Deploy(); err != nil {
	//	panic(err)
	//}
	//
	//modHandler := mod.NewHandler()
	//err := modHandler.Deploy()
	//if err != nil {
	//	panic(err)
	//}
	if err := controller.Start(); err != nil {
		panic(err)
	}
	server.Run()
}
