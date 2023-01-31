package modreq

import "terraria-run/internal/common/model"

type ModIds struct {
	ModIds []string `json:"mod_ids"`
}

type Mods struct {
	Mods []model.Mod `json:"mods"`
}
