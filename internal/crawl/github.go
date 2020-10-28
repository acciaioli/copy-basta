package crawl

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"mime"

	"copy-basta/internal/clients/github"
	"copy-basta/internal/common"

	"copy-basta/internal/common/log"
)

type githubCrawler struct {
	ghc *github.Client
}

func NewGithubCrawler(ghc *github.Client) Crawler {
	return &githubCrawler{ghc: ghc}
}

func (c *githubCrawler) Crawl() ([]File, error) {
	url := c.ghc.ZipArchiveURL()
	headers, data, err := c.ghc.DoGetRequest(url)
	if err != nil {
		return nil, err
	}

	if len(headers["Content-Disposition"]) != 1 {
		log.L.DebugWithData(
			"github response error: Content-Disposition",
			log.Data{"url": url, "content-disposition": headers["Content-Disposition"]},
		)
		return nil, errors.New("failed to download github zip archive")
	}
	_, params, err := mime.ParseMediaType(headers["Content-Disposition"][0])
	if err != nil {
		log.L.DebugWithData(
			"github response error: Content-Disposition",
			log.Data{"url": url, "content-disposition": headers["Content-Disposition"], "error": err.Error()},
		)
		return nil, errors.New("failed to download github zip archive")
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
