package modresp

import (
	"terraria-run/internal/common/model"
)

type Response struct {
	Mods []model.Mod `json:"mods"`
}
