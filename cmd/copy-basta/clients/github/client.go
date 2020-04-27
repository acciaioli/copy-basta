package github

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"copy-basta/cmd/copy-basta/common/log"
)

/*
reference: https://developer.github.com/v3/repos/contents/#get-contents
*/

const (
	ContentTypeFile    = "file"
	ContentEncodingB64 = "base64"
)

type RepoContent struct {
	Type     string  `json:"type"`
	Path     string  `json:"path"`
	Encoding *string `json:"encoding"`
	Content  *string `json:"content"`
}

type Client struct {
	repoNamespace string
	repoID        string
	branch        string
}

func NewClient(repoRef string) (*Client, error) {
	// repo ref is expected to be something like "{namespace}/{repo-id}" (example `acciaioli/copy-basta`)
	repo := strings.Split(repoRef, "/")
	if len(repo) != 2 {
		log.L.DebugWithData("invalid repo: split error", log.Data{"repo-ref": repoRef})
		return nil, fmt.Errorf("github client error: invalid repo reference `%s`", repoRef)
	}

	ghc := Client{
		repoNamespace: repo[0],
		repoID:        repo[1],
		branch:        "master",
	}
	return &ghc, nil
}

func (ghc *Client) ApiContentsURL() string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", ghc.repoNamespace, ghc.repoID)
}

func (ghc *Client) ZipArchiveURL() string {
	return fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", ghc.repoNamespace, ghc.repoID, ghc.branch)
}

func (ghc *Client) DoGetRequest(url string) (http.Header, []byte, error) {
	log.L.DebugWithData("github api request", log.Data{"url": url})
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"url": url, "error": err.Error()})
		return nil, nil, errors.New("failed to create github api request")
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.L.DebugWithData("failed to close response body", log.Data{"url": url})
		}
	}()
	if resp.StatusCode != http.StatusOK {
		log.L.DebugWithData("github api status code not ok", log.Data{"url": url, "status-code": resp.StatusCode})
		return nil, nil, errors.New("github api response status error")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"url": url, "error": err.Error()})
		return nil, nil, errors.New("failed to read github api response")
	}

	return resp.Header, data, nil
}

func (ghc *Client) GetContents(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", ghc.ApiContentsURL(), path)
	_, data, err := ghc.DoGetRequest(url)
	return data, err
}

func (ghc *Client) GetContentsFileData(fpath string) ([]byte, error) {
	b, err := ghc.GetContents(fpath)
	if err != nil {
		log.L.Warn(fmt.Sprintf("failed to find %s in this repo. continuing without it.", fpath))
		return nil, err
	}

	var entry RepoContent
	err = json.Unmarshal(b, &entry)
	if err != nil {
		log.L.DebugWithData(
			"external error", log.Data{"error": err.Error()},
		)
		return nil, errors.New("failed to decoded github api json response")
	}
	if entry.Type != ContentTypeFile {
		log.L.DebugWithData(
			"expected to be a file but ins't", log.Data{"path": entry.Path, "type": entry.Type},
		)
		return nil, errors.New("expected github api json response to be a file but ins't")
	}
	if entry.Encoding == nil {
		log.L.DebugWithData(
			"github content error: nil encoding", log.Data{"path": entry.Path},
		)
		return nil, errors.New("failed to get necessary encoding from the github api json response")
	}
	if *entry.Encoding != ContentEncodingB64 {
		log.L.DebugWithData(
			"github content error: unknown encoding", log.Data{"path": entry.Path, "encoding": entry.Encoding},
		)
		return nil, errors.New("failed to get necessary encoding from the github api json response")
	}
	if entry.Content == nil {
		log.L.DebugWithData(
			"github content error: nil content", log.Data{"path": entry.Path},
		)
		return nil, errors.New("failed to get necessary content from the github api json response")
	}

	decodedContent, err := base64.StdEncoding.DecodeString(*entry.Content)
	if err != nil {
		log.L.DebugWithData(
			"external error", log.Data{"path": entry.Path, "error": err.Error()},
		)
		return nil, errors.New("failed decode content from the github api json response")
	}

	return decodedContent, nil
}
