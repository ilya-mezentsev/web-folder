package app

import (
	"encoding/json"
	"fmt"
	"github.com/ilya-mezentsev/folder-info/internal/services"
	"github.com/ilya-mezentsev/folder-info/internal/transport"
	"log/slog"
	"os"
)

type Config struct {
	RootDirs []string `json:"root_dirs"`
	AppAddr  string   `json:"app_addr"`
}

func Main() {
	cfg := mustParseConfig(os.Getenv("CONFIG_PATH"))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	core := services.New(logger, cfg.RootDirs)

	transport.Run(
		cfg.AppAddr,
		core,
		logger,
	)
}

func mustParseConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("unable to read settings file: %s, err: %v", path, err))
	}

	var s Config
	err = json.Unmarshal(data, &s)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal settings to struct: %v", err))
	}

	return s
}
