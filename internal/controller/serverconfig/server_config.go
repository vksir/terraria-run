package serverconfig

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"regexp"
	"terraria-run/assets"
	"terraria-run/internal/common/constant"
)

type Handler struct {
	content string
}

func NewHandler() *Handler {
	bytes, err := assets.FS.ReadFile("asserts/serverconfig.txt")
	if err != nil {
		panic(err)
	}
	s := Handler{content: string(bytes)}
	s.setRequiredOptions()
	s.setOptionalOptions()
	return &s
}

func (h *Handler) Deploy() error {
	return os.WriteFile(constant.ServerConfigPath, []byte(h.content), 0644)
}

func (h *Handler) setRequiredOptions() {
	h.set("world", filepath.Join(constant.WorldDir, fmt.Sprintf("%h.wld", viper.GetString("server_config.world_name"))))
	h.set("autocreate", viper.GetString("server_config.auto_create"))
	h.set("worldname", viper.GetString("server_config.world_name"))
	h.set("difficulty", viper.GetString("server_config.difficulty"))
	h.set("worldpath", constant.WorldDir)
}

func (h *Handler) setOptionalOptions() {
	h.setIfNotEmpty("seed", viper.GetString("server_config.seed"))
	h.setIfNotEmpty("maxplayers", viper.GetString("server_config.max_players"))
	h.setIfNotEmpty("password", viper.GetString("server_config.password"))
}

func (h *Handler) set(key, value string) {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^#%s=.*$`, key))
	h.content = pattern.ReplaceAllString(h.content, fmt.Sprintf("%s=%s", key, value))
}

func (h *Handler) setIfNotEmpty(key, value string) {
	if value == "" {
		return
	}
	h.set(key, value)
}
