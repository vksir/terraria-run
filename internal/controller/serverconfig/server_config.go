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

type ServerConfig struct {
	content string
}

func NewServerConfig() *ServerConfig {
	bytes, err := assets.FS.ReadFile("asserts/serverconfig.txt")
	if err != nil {
		panic(err)
	}
	s := ServerConfig{content: string(bytes)}
	s.setRequiredOptions()
	s.setOptionalOptions()
	return &s
}

func (s *ServerConfig) Deploy() error {
	return os.WriteFile(constant.ServerConfigPath, []byte(s.content), 0644)
}

func (s *ServerConfig) setRequiredOptions() {
	s.set("world", filepath.Join(constant.WorldDir, fmt.Sprintf("%s.wld", viper.GetString("server_config.world_name"))))
	s.set("autocreate", viper.GetString("server_config.auto_create"))
	s.set("worldname", viper.GetString("server_config.world_name"))
	s.set("difficulty", viper.GetString("server_config.difficulty"))
	s.set("worldpath", constant.WorldDir)
}

func (s *ServerConfig) setOptionalOptions() {
	s.setIfNotEmpty("seed", viper.GetString("server_config.seed"))
	s.setIfNotEmpty("maxplayers", viper.GetString("server_config.max_players"))
	s.setIfNotEmpty("password", viper.GetString("server_config.password"))
}

func (s *ServerConfig) set(key, value string) {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^#%s=.*$`, key))
	s.content = pattern.ReplaceAllString(s.content, fmt.Sprintf("%s=%s", key, value))
}

func (s *ServerConfig) setIfNotEmpty(key, value string) {
	if value == "" {
		return
	}
	s.set(key, value)
}
