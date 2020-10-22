package load

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"

	"copy-basta/cmd/copy-basta/clients/github"
	"copy-basta/cmd/copy-basta/common"
	"copy-basta/cmd/copy-basta/common/log"
)

type Loader interface {
	Load() ([]File, error)
}

type File struct {
	Path   string
	Mode   os.FileMode
	Reader io.Reader
}

// =========== disk loader =========== //

type DiskLoader struct {
	root string
}

func NewDiskLoader(root string) (*DiskLoader, error) {
	return &DiskLoader{root: root}, nil
}

func (l *DiskLoader) Load() ([]File, error) {
	var files []File

	err := filepath.Walk(l.root, func(fPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		r, err := os.Open(fPath)
		if err != nil {
			log.L.DebugWithData("external error", log.Data{"error": err.Error()})
			return err
		}
		files = append(files, File{Path: common.TrimRootDir(fPath), Mode: info.Mode(), Reader: r})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// =========== github loader =========== //

type GithubLoader struct {
	ghc *github.Client
}

func NewGithubLoader(ghc *github.Client) (*GithubLoader, error) {
	return &GithubLoader{ghc: ghc}, nil
}

func (l *GithubLoader) Load() ([]File, error) {
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

	log.L.Debug(fmt.Sprintf("archive filename: %s", params["filename"]))

	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
		return nil, errors.New("github api response error: zip reader failed")
	}

	var files []File

	for _, zfile := range r.File {
		info := zfile.FileInfo()
		if info.IsDir() {
			continue
		}

		r, err := zfile.Open()
		if err != nil {
			log.L.DebugWithData("external error", log.Data{"error": err.Error()})
			return nil, err
		}
		files = append(files, File{
			Path:   common.TrimRootDir(zfile.Name),
			Mode:   info.Mode(),
			Reader: r,
		})
	}

	return files, nil
}
