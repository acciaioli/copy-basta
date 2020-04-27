package github

import (
	"archive/zip"
	"bytes"
	"copy-basta/cmd/common"
	"copy-basta/cmd/common/log"
	"copy-basta/cmd/generate/parse"
	"copy-basta/cmd/generate/parse/ignore"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"strings"
)

const (
	SourcePrefix = "https://github.com/"
)

type Loader struct {
	ghc *Client
}

func NewLoader(repo string) (*Loader, error) {
	ghc, err := NewGitHubAPIClient(repo)
	if err != nil {
		return nil, err
	}
	return &Loader{ghc: ghc}, nil
}

func (l *Loader) LoadFiles() ([]parse.LoadedFile, error) {
	url := l.ghc.zipArchiveURL()
	headers, data, err := l.ghc.DoGetRequest(url)
	if len(headers["Content-Disposition"]) != 1 {
		log.L.DebugWithData("github response error: Content-Disposition", log.Data{"url": url, "content-disposition": headers["Content-Disposition"]})
		return nil, errors.New("github api response error: invalid `Content-Disposition` header")
	}
	_, params, err := mime.ParseMediaType(headers["Content-Disposition"][0])
	if err != nil {
		log.L.DebugWithData("github response error: Content-Disposition", log.Data{"url": url, "content-disposition": headers["Content-Disposition"], "error": err.Error()})
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
	b, err := l.ghc.GetContents(common.IgnoreFile)
	if err != nil {
		log.L.Warn(fmt.Sprintf("failed to find %s in this repo. continuing without it.", common.IgnoreFile))
		return ignore.New(root, nil)
	}

	var entry RepositoryEntry
	err = json.Unmarshal(b, &entry)
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
		return nil, errors.New("failed to decoded github api json response")
	}
	if entry.Type != contentTypeFile {
		log.L.DebugWithData("ignore file is not a file", log.Data{"path": entry.Path, "type": entry.Type})
		return nil, errors.New("the ignore file of this repo is not a file")
	}
	if entry.Content == nil {
		log.L.DebugWithData("github content error: nil content", log.Data{"path": entry.Path})
		return nil, errors.New("failed to get necessary data from the github api json response")
	}

	return ignore.New(root, strings.NewReader(*entry.Content))
}
