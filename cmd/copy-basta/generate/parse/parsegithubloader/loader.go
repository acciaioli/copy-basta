package parsegithubloader

import (
	"archive/zip"
	"bytes"
	"errors"
	"mime"

	"copy-basta/cmd/copy-basta/clients/github"
	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/common/log"
	"copy-basta/cmd/copy-basta/generate/parse"
	"copy-basta/cmd/copy-basta/generate/parse/ignore"
)

type Loader struct {
	ghc *github.Client
}

func New(ghc *github.Client) (*Loader, error) {
	return &Loader{ghc: ghc}, nil
}

func (l *Loader) LoadFiles() ([]parse.LoadedFile, error) {
	url := l.ghc.ZipArchiveURL()
	headers, data, err := l.ghc.DoGetRequest(url)
	if len(headers["Content-Disposition"]) != 1 {
		log.L.DebugWithData(
			"github response error: Content-Disposition",
			log.Data{"url": url, "content-disposition": headers["Content-Disposition"]},
		)
		return nil, errors.New("github api response error: invalid `Content-Disposition` header")
	}
	_, params, err := mime.ParseMediaType(headers["Content-Disposition"][0])
	if err != nil {
		log.L.DebugWithData(
			"github response error: Content-Disposition",
			log.Data{"url": url, "content-disposition": headers["Content-Disposition"], "error": err.Error()},
		)
		return nil, errors.New("github api response error: invalid `Content-Disposition` header")
	}

	ignorer, err := l.getIgnorer(params["filename"])
	if err != nil {
		return nil, err
	}

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
		return nil, errors.New("github api response error: zip reader failed")
	}

	var files []parse.LoadedFile

	for _, zfile := range r.File {
		if ignorer.Ignore(zfile.Name) {
			continue
		}

		info := zfile.FileInfo()
		if info.IsDir() {
			continue
		}

		r, err := zfile.Open()
		if err != nil {
			log.L.DebugWithData("external error", log.Data{"error": err.Error()})
			return nil, err
		}
		files = append(files, parse.LoadedFile{
			Path:   zfile.Name,
			Mode:   info.Mode(),
			Reader: r,
		})
	}

	return files, nil
}

func (l *Loader) getIgnorer(root string) (*ignore.Ignorer, error) {
	b, err := l.ghc.GetContentsFileData(common.IgnoreFile)
	if err != nil {
		return ignore.New(root, nil)
	}

	return ignore.New(root, bytes.NewReader(b))
}
