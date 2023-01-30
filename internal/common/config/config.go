package config

import (
	"github.com/spf13/viper"
	"os"
	"terraria-run/internal/common/constant"
)

const (
	SizeSmall  = 1
	SizeMedium = 2
	SizeLarge  = 3

	DifficultyClassic = 0
	DifficultyExpert  = 1
	DifficultyMaster  = 2
	DifficultyJourney = 3
)

func init() {
	setDefault()
	Read()
}

func Read() {
	if _, err := os.Stat(constant.ConfigPath()); os.IsNotExist(err) {
		return
	}
	r, err := os.Open(constant.ConfigPath())
	if err != nil {
		panic(err)
	}
	err = viper.ReadConfig(r)
	if err != nil {
		panic(err)
	}
}

func Write() {
	err := viper.WriteConfigAs(constant.ConfigPath())
	if err != nil {
		panic(err)
	}
}

// serverconfig.txt: https://terraria.wiki.gg/wiki/Guide:Setting_up_a_Terraria_server
func setDefault() {
	viper.SetDefault("world_name", "NeutronStar")
	viper.SetDefault("password", "")
	viper.SetDefault("server_config", map[string]any{
		"auto_create": SizeLarge,
		"seed":        "",
		"world_name":  "NeutronStar",
		"difficulty":  DifficultyMaster,
		"max_players": 8,
		"port":        7777,
		"password":    "",
	})
}
