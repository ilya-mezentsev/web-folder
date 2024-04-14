package services

import (
	"log/slog"
	"os"
	"time"
)

type DetailedFile struct {
	Name    string `json:"name"`
	Size    string `json:"size"`
	Created string `json:"created"`
}

func (c Core) File(path string) (DetailedFile, error) {
	if !c.isPathAllowed(path) {
		return DetailedFile{}, ErrPathIsNotAllowed
	}

	stat, err := os.Stat(path)
	if err != nil {
		c.logger.Error("unable to get file stat", slog.String("path", path), slog.Any("err", err))

		return DetailedFile{}, ErrUnknown
	}

	return DetailedFile{
		Name:    stat.Name(),
		Size:    byteCountIEC(stat.Size()),
		Created: stat.ModTime().Format(time.DateTime),
	}, nil
}
