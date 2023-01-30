package constant

import (
	"os"
	"path/filepath"
)

func Home() string {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return h
}

func InstallDir() string {
	return filepath.Join(Home(), "neutron-star", "terraria-run")
}

func Workspace() string {
	p := filepath.Join(Home(), ".terraria-run")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.Mkdir(p, 0755); err != nil {
			panic(err)
		}
	}
	return p
}

func ServerLogPath() string {
	return filepath.Join(Workspace(), "server.log")
}

func DotnetPath() string {
	return filepath.Join(InstallDir(), "tModLoader/dotnet/6.0.0/dotnet")
}

func ServerConfigPath() string {
	return filepath.Join(Workspace(), "serverconfig.txt")
}

func TModLoaderLogPath() string {
	return filepath.Join(Workspace(), "tModLoader.log")
}

func ConfigPath() string {
	return filepath.Join(Workspace(), "config.json")
}
