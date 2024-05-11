package app

import (
	"flag"
	"github.com/ilya-mezentsev/folder-info/pkg/dto"
	"github.com/ilya-mezentsev/folder-view/internal/services/api"
	"github.com/ilya-mezentsev/folder-view/internal/transport"
	"log/slog"
	"os"
	"path"
)

var (
	addr         = flag.String("addr", "localhost:8001", "address for application to listen")
	infoApiAddr  = flag.String("info-api", "localhost:8000", "folder info API address")
	templatesDir = flag.String("templates-dir", "/dev/null", "html templates dir")
)

//goland:noinspection HttpUrlsUsage
const apiProtocol = "http://"

func init() {
	flag.Parse()
}

func Main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	rootsFetcher := api.MustNew[[]string](apiProtocol + path.Join(*infoApiAddr, "roots"))
	folderFetcher := api.MustNew[dto.DirInfo](apiProtocol + path.Join(*infoApiAddr, "folder"))
	fileFetcher := api.MustNew[dto.DetailedFile](apiProtocol + path.Join(*infoApiAddr, "file"))

	transport.Run(
		*addr,
		rootsFetcher,
		folderFetcher,
		fileFetcher,
		*templatesDir,
		logger,
	)
}
