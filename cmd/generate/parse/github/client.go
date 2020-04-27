package github

import (
	"copy-basta/cmd/common/log"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GithubRepositoryContent represents a LoadedFile or directory in a github repository.
type RepositoryEntry struct {
	Type string `json:"type"`
	Path string `json:"path"`

	// Target is only set if the type is "symlink" and the target is not a normal LoadedFile.
	// If Target is set, Path will be the symlink path.
	Target   *string `json:"target,omitempty"`
	Encoding *string `json:"encoding,omitempty"`
	Size     *int    `json:"size,omitempty"`
	Name     *string `json:"name,omitempty"`

	// Content contains the actual LoadedFile content, which may be encoded.
	// Callers should call GetContent which will decode the content if
	// necessary.
	Content     *string `json:"content,omitempty"`
	SHA         *string `json:"sha,omitempty"`
	URL         *string `json:"url,omitempty"`
	GitURL      *string `json:"git_url,omitempty"`
	HTMLURL     *string `json:"html_url,omitempty"`
	DownloadURL *string `json:"download_url,omitempty"`
}

const (
	contentTypeFile    = "file"
	contentTypeDir     = "dir"
	contentTypeSymlink = "symlink"
)

const (
	defaultMode = 0666
)

type Client struct {
	repoNamespace string
	repoID        string
	branch        string
}

func NewGitHubAPIClient(repoRef string) (*Client, error) {
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

func (ghc *Client) apiContentsURL() string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", ghc.repoNamespace, ghc.repoID)
}

func (ghc *Client) zipArchiveURL() string {
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
	url := fmt.Sprintf("%s/%s", ghc.apiContentsURL(), path)
	_, data, err := ghc.DoGetRequest(url)
	return data, err
}
