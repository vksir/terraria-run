package config

import (
	"github.com/spf13/viper"
	"os"
	"terraria-run/internal/common/constant"
	"terraria-run/internal/common/model/model"
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
}

func Read() {
	if _, err := os.Stat(constant.ConfigPath); os.IsNotExist(err) {
		return
	}
	r, err := os.Open(constant.ConfigPath)
	if err != nil {
		panic(err)
	}
	err = viper.ReadConfig(r)
	if err != nil {
		panic(err)
	}
}

func Write() {
	err := viper.WriteConfigAs(constant.ConfigPath)
	if err != nil {
		panic(err)
	}
}

func setDefault() {
	// https://terraria.wiki.gg/wiki/Guide:Setting_up_a_Terraria_server
	viper.SetDefault("server_config", map[string]any{
		"auto_create": SizeLarge,
		"seed":        "",
		"world_name":  "NeutronStar",
		"difficulty":  DifficultyMaster,
		"max_players": 8,
		"port":        7777,
		"password":    "",
	})
	viper.SetDefault("mod", map[int]model.Mod{
		2619954303: {
			ID:     2619954303,
			Name:   "Recipe Browser",
			Enable: true,
		},
	})
}
