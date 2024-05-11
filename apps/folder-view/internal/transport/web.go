package transport

import (
	"github.com/ilya-mezentsev/folder-info/pkg/dto"
	"github.com/ilya-mezentsev/folder-view/internal/services/api"
	"html/template"
	"log/slog"
	"net/http"
	"path"
	"slices"
)

func Run(
	addr string,
	rootsFetcher api.Service[[]string],
	folderFetcher api.Service[dto.DirInfo],
	fileFetcher api.Service[dto.DetailedFile],
	templatesDir string,
	logger *slog.Logger,
) {
	errTemplate := template.Must(template.ParseFiles(path.Join(templatesDir, "error.tmpl")))

	http.HandleFunc("/", mustMakeRootsHandler(
		rootsFetcher,
		templatesDir,
		errTemplate,
		logger,
	))

	http.HandleFunc("/folder", mustMakeFolderHandler(
		rootsFetcher,
		folderFetcher,
		templatesDir,
		errTemplate,
		logger,
	))

	http.HandleFunc("/file", mustMakeFileHandler(
		fileFetcher,
		templatesDir,
		errTemplate,
		logger,
	))

	fs := http.FileServer(http.Dir(templatesDir))
	http.Handle("/css/", http.StripPrefix("/css", fs))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("unable to ListenAndServe", slog.Any("err", err), slog.String("addr", addr))
	}
}

func mustMakeRootsHandler(
	rootsFetcher api.Service[[]string],
	templatesDir string,
	errTemplate *template.Template,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	rootsTemplate := template.Must(template.ParseFiles(path.Join(templatesDir, "roots.tmpl")))

	return func(w http.ResponseWriter, r *http.Request) {
		roots, err := rootsFetcher.Fetch()
		if err != nil {
			handleError(
				"roots fetch error",
				err,
				errTemplate,
				w,
				logger,
			)
			return
		}

		err = rootsTemplate.Execute(w, map[string][]string{
			"Roots": roots,
		})
		if err != nil {
			handleError(
				"roots render error",
				err,
				errTemplate,
				w,
				logger,
			)
		}
	}
}

func mustMakeFolderHandler(
	rootsFetcher api.Service[[]string],
	folderFetcher api.Service[dto.DirInfo],
	templatesDir string,
	errTemplate *template.Template,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	folderTemplate := template.Must(template.ParseFiles(path.Join(templatesDir, "folder.tmpl")))

	return func(w http.ResponseWriter, r *http.Request) {
		roots, err := rootsFetcher.Fetch()
		if err != nil {
			handleError(
				"roots fetch error",
				err,
				errTemplate,
				w,
				logger,
			)
			return
		}

		folderPath := r.URL.Query().Get("path")
		dirInfo, err := folderFetcher.FetchWithQuery("path", folderPath)
		if err != nil {
			handleError(
				"folder fetch error",
				err,
				errTemplate,
				w,
				logger.With(slog.String("path", folderPath)),
			)
			return
		}

		fromRoots := slices.Contains(roots, folderPath)
		err = folderTemplate.Execute(w, map[string]any{
			"Path":  dirInfo.Path,
			"Files": dirInfo.Files,
			"Dirs":  dirInfo.Dirs,

			"ShowParent": !fromRoots,
			"Parent": makePrevPath(
				folderPath,
				fromRoots,
			),
		})
		if err != nil {
			handleError(
				"folder render error",
				err,
				errTemplate,
				w,
				logger.With(slog.String("path", folderPath)),
			)
		}
	}
}

func makePrevPath(folder string, fromRoots bool) string {
	if fromRoots {
		return "/"
	}

	return "/folder?path=" + path.Dir(folder)
}

func mustMakeFileHandler(
	fileFetcher api.Service[dto.DetailedFile],
	templatesDir string,
	errTemplate *template.Template,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileTemplate := template.Must(template.ParseFiles(path.Join(templatesDir, "file.tmpl")))

		filePath := r.URL.Query().Get("path")
		fileInfo, err := fileFetcher.FetchWithQuery("path", filePath)
		if err != nil {
			handleError(
				"file fetch error",
				err,
				errTemplate,
				w,
				logger.With(slog.String("path", filePath)),
			)
			return
		}

		err = fileTemplate.Execute(w, map[string]string{
			"Name":    fileInfo.Name,
			"Size":    fileInfo.Size,
			"Created": fileInfo.Created,
			"Dir":     path.Dir(filePath),
		})
		if err != nil {
			handleError(
				"folder render error",
				err,
				errTemplate,
				w,
				logger.With(slog.String("path", filePath)),
			)
		}
	}
}

func handleError(
	msg string,
	err error,
	errTemplate *template.Template,
	w http.ResponseWriter,
	logger *slog.Logger,
) {
	logger.Error(msg, slog.Any("err", err))

	err = errTemplate.Execute(w, nil)
	if err != nil {
		logger.Error("unable to execute error template", slog.Any("err", err))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal error"))
	}
}
