package services

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const extPrefix = "."

type (
	Result struct {
		Path string `json:"path"`

		Files []File `json:"files"`
		Dirs  []Dir  `json:"dirs"`
	}

	File struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Size string `json:"size"`
	}

	Dir struct {
		Name string `json:"name"`
		Size string `json:"size"`
	}
)

// Folder collects folder info
func (c Core) Folder(path string) (Result, error) {
	if !c.isPathAllowed(path) {
		return Result{}, ErrPathIsNotAllowed
	}

	f, err := os.Open(path)
	if err != nil {
		c.logger.Error("unable to open path", slog.String("path", path), slog.Any("err", err))

		return Result{}, ErrUnknown
	}

	files, err := f.ReadDir(0)
	if err != nil {
		c.logger.Error("unable to read dir by path", slog.String("path", path), slog.Any("err", err))

		return Result{}, ErrUnknown
	}

	var (
		eg         errgroup.Group
		result     Result
		resultLock sync.Mutex
	)

	for _, file := range files {
		lFile := file
		eg.Go(func() error {
			return c.processFile(
				lFile,
				&result,
				&resultLock,
				path,
			)
		})
	}

	err = eg.Wait()
	if err != nil {
		c.logger.Error("unable to process files by path", slog.String("path", path), slog.Any("err", err))

		return Result{}, ErrUnknown
	}

	result.Path = path

	return result, nil
}

func (c Core) processFile(
	file os.DirEntry,
	result *Result,
	resultLock *sync.Mutex,
	rootPath string,
) error {

	var (
		fileInfo os.FileInfo
		err      error
	)
	fileInfo, err = file.Info()
	if err != nil {
		return err
	}

	name := file.Name()
	size := fileInfo.Size()

	if file.IsDir() {
		size, err = dirSize(path.Join(rootPath, name))
		if err != nil {
			return err
		}

		resultLock.Lock()
		result.Dirs = append(result.Dirs, Dir{
			Name: name,
			Size: byteCountIEC(size),
		})
		resultLock.Unlock()
	} else {
		resultLock.Lock()
		result.Files = append(result.Files, File{
			Name: name,
			Type: strings.TrimPrefix(filepath.Ext(name), extPrefix),
			Size: byteCountIEC(size),
		})
		resultLock.Unlock()
	}

	return nil
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
