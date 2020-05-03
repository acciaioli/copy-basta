package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"copy-basta/cmd/copy-basta/common/log"
)

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
