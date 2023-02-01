package util

import (
	"fmt"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

var log = zap.S()

func Copy(src, dst string) error {
	return RunCmd(fmt.Sprintf("cp -rf %s %s", src, dst))
}

func Remove(path string) error {
	return RunCmd(fmt.Sprintf("rm -rf %s", path))
}

func RunCmd(cmd string) error {
	if cmd == "" {
		return fmt.Errorf("empty cmd")
	}
	cmdLst := strings.Split(cmd, " ")
	c := exec.Command(cmdLst[0])
	if len(cmdLst) > 1 {
		c.Args = append(c.Args, cmdLst[1:]...)
	}
	out, err := c.CombinedOutput()
	if err != nil {
		log.Error("Run cmd failed", err, "cmd", cmd, "out", out)
		return err
	}
	return nil
}
