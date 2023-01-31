package mod

import (
	"bufio"
	"encoding/json"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"terraria-run/internal/common/constant"
	"terraria-run/internal/common/model/model"
	"terraria-run/internal/common/util"
	"time"
)

type Handler struct {
	mods          []model.Mod
	modEnableData []string
}

func NewHandler() *Handler {
	h := Handler{}
	mods := viper.Get("mod").(map[int]model.Mod)
	for id := range mods {
		if mods[id].Enable {
			h.mods = append(h.mods, mods[id])
		}
	}
	return &h
}

func (h *Handler) Deploy() error {
	if len(h.mods) == 0 {
		return nil
	}
	zap.S().Infof("Begin deploy %d mods", len(h.mods))
	if err := h.downloadMods(); err != nil {
		zap.S().Error("Download mods failed", err)
		return err
	}
	if err := h.copyMods(); err != nil {
		zap.S().Error("Copy mods failed", err)
		return err
	}
	if err := h.writeEnableJson(); err != nil {
		zap.S().Error("Write enable.json failed", err)
		return err
	}
	zap.S().Info("Deploy mods success")
	return nil
}

func (h *Handler) downloadMods() error {
	cmd := exec.Command("steamcmd", "+login", "anonymous")
	for i := range h.mods {
		cmd.Args = append(cmd.Args, "+workshop_download_item", "1281930", strconv.Itoa(h.mods[i].ID))
	}
	cmd.Args = append(cmd.Args, "+quit")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	r := io.MultiReader(stdout, stderr)
	go func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			zap.S().Debug(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			zap.S().Error("Output scan failed", err)
		}
	}(r)
	return cmd.Run()
}

func (h *Handler) copyMods() error {
	if err := util.Remove(constant.ModDir); err != nil {
		return err
	}
	if err := os.Mkdir(constant.ModDir, 0755); err != nil {
		return err
	}
	for i := range h.mods {
		if err := h.copyMod(h.mods[i].ID); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) writeEnableJson() error {
	bytes, err := json.MarshalIndent(h.modEnableData, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(constant.EnableJson, bytes, 0644)
}

func (h *Handler) copyMod(id int) error {
	latestModDir, err := getLatestModDir(id)
	if err != nil {
		return err
	}
	files, err := os.ReadDir(latestModDir)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tmod") {
			if err := util.Copy(filepath.Join(latestModDir, file.Name()), constant.ModDir); err != nil {
				return err
			}
			h.modEnableData = append(h.modEnableData, strings.TrimSuffix(file.Name(), ".tmod"))
			zap.S().Infof("Copy mod from %s to %s", filepath.Join(latestModDir, file.Name()), constant.ModDir)
		}
	}
	return nil
}

type sorter []struct {
	dirname string
	time    time.Time
}

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Less(i, j int) bool {
	return s[i].time.Unix() < s[j].time.Unix()
}

func (s sorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func getLatestModDir(id int) (string, error) {
	modDir := filepath.Join(constant.SteamModDir, strconv.Itoa(id))
	if _, err := os.Stat(modDir); os.IsNotExist(err) {
		return "", err
	}
	files, err := os.ReadDir(modDir)
	if err != nil {
		return "", err
	}
	var s sorter
	for _, file := range files {
		if file.IsDir() {
			t, err := time.Parse("2006.2", file.Name())
			if err != nil {
				return "", err
			}
			s = append(s, sorter{{
				dirname: file.Name(),
				time:    t,
			}}...)
		}
	}
	sort.Sort(s)
	return filepath.Join(modDir, s[len(s)-1].dirname), nil
}
