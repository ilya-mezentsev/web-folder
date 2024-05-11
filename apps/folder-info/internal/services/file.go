package services

import (
	"github.com/ilya-mezentsev/folder-info/pkg/dto"
	"log/slog"
	"os"
	"time"
)

func (c Core) File(path string) (dto.DetailedFile, error) {
	if !c.isPathAllowed(path) {
		return dto.DetailedFile{}, ErrPathIsNotAllowed
	}

	stat, err := os.Stat(path)
	if err != nil {
		c.logger.Error("unable to get file stat", slog.String("path", path), slog.Any("err", err))

		return dto.DetailedFile{}, ErrUnknown
	}

	return dto.DetailedFile{
		Name:    stat.Name(),
		Size:    byteCountIEC(stat.Size()),
		Created: stat.ModTime().Format(time.DateTime),
	}, nil
}
