package model

type Config struct {
	ServerConfig ServerConfig   `json:"server_config"`
	IdToMod      map[string]Mod `json:"mod"`
}

type ServerConfig struct {
	AutoCreate int    `json:"auto_create"`
	Difficulty int    `json:"difficulty"`
	MaxPlayers int    `json:"max_players"`
	Password   string `json:"password"`
	Port       int    `json:"port"`
	Seed       string `json:"seed"`
	WorldName  string `json:"world_name"`
}

type Mod struct {
	Id     string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Enable bool   `json:"enable" binding:"required"`
}
