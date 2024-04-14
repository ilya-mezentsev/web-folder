package services

import "log/slog"

type Core struct {
	logger   *slog.Logger
	rootDirs []string
}

func New(
	logger *slog.Logger,
	rootDirs []string,
) Core {

	return Core{
		logger:   logger,
		rootDirs: rootDirs,
	}
}
