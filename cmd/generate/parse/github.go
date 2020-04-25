package parse

import (
	"bytes"
	"copy-basta/cmd/common/log"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// GithubRepositoryContent represents a file or directory in a github repository.
type RepositoryEntry struct {
	Type string `json:"type"`
	Path string `json:"path"`

	// Target is only set if the type is "symlink" and the target is not a normal file.
	// If Target is set, Path will be the symlink path.
	Target   *string `json:"target,omitempty"`
	Encoding *string `json:"encoding,omitempty"`
	Size     *int    `json:"size,omitempty"`
	Name     *string `json:"name,omitempty"`

	// Content contains the actual file content, which may be encoded.
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

type GitHubAPIClient struct {
	repoNamespace string
	repoID        string
}

func NewGitHubAPIClient(repoRef string) (*GitHubAPIClient, error) {
	// repo ref is expected to be something like "{host}/{namespace}/{repo-id}" (example `github.com/acciaioli/copy-basta`)
	repo := strings.Split(repoRef, "/")
	if len(repo) != 3 {
		log.L.DebugWithData("invalid repo: split error", log.Data{"repo-ref": repoRef})
		return nil, fmt.Errorf("github client error: invalid repo reference `%s`", repoRef)
	}

	ghc := GitHubAPIClient{repoNamespace: repo[1], repoID: repo[2]}
	return &ghc, nil
}

func (ghc *GitHubAPIClient) contentsURL() string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", ghc.repoNamespace, ghc.repoID)
}

func (ghc *GitHubAPIClient) DoGetRequest(url string) ([]byte, error) {
	log.L.DebugWithData("github api request", log.Data{"url": url})
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"url": url, "error": err.Error()})
		return nil, errors.New("failed to create github api request")
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.L.DebugWithData("failed to close response body", log.Data{"url": url})
		}
	}()
	if resp.StatusCode != http.StatusOK {
		log.L.DebugWithData("github api status code not ok", log.Data{"url": url, "status-code": resp.StatusCode})
		return nil, errors.New("github api response status error")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil
	}

	return b, nil
}

func (ghc *GitHubAPIClient) GetContents(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", ghc.contentsURL(), path)
	return ghc.DoGetRequest(url)
}

type GitHubParser struct {
	ghc *GitHubAPIClient
}

func NewGitHubParser(repo string) (*GitHubParser, error) {
	ghc, err := NewGitHubAPIClient(repo)
	if err != nil {
		return nil, err
	}
	return &GitHubParser{ghc: ghc}, nil
}

func (p *GitHubParser) Parse() ([]file, error) {
	var files []file
	err := p.parse("", files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (p *GitHubParser) githubAPIGetContent(path string) ([]RepositoryEntry, error) {
	b, err := p.ghc.GetContents(path)
	if err != nil {
		return nil, err
	}

	var entries []RepositoryEntry
	err = json.Unmarshal(b, &entries)
	if err != nil {
		log.L.DebugWithData("external error", log.Data{"error": err.Error()})
		return nil, errors.New("failed to decoded github api json response")
	}

	return entries, nil
}

func (p *GitHubParser) githubAPIFileContentDownload(url string) ([]byte, error) {
	return p.ghc.DoGetRequest(url)
}

func (p *GitHubParser) parse(path string, files []file) error {
	entries, err := p.githubAPIGetContent(path)
	if err != nil {
		return err
	}
	for {
		if len(entries) == 0 {
			break
		}
		entry := entries[0]
		entries = entries[1:]
		switch entry.Type {

		case contentTypeFile:
			content, err := p.githubAPIFileContentDownload(*entry.DownloadURL)
			if err != nil {
				return err
			}

			files = append(files, file{
				path: entry.Path,
				mode: defaultMode,
				r:    bytes.NewReader(content),
			})

		case contentTypeDir:
			err = p.parse(entry.Path, files)
			if err != nil {
				return err
			}

		case contentTypeSymlink:
			log.L.DebugWithData("symlink type is currently not support", log.Data{"path": entry.Path})
			return errors.New("github repo parsing error: symlink type is currently not support")

		default:
			log.L.DebugWithData("symlink type is currently not support", log.Data{"path": entry.Path, "type": entry.Type})
			return errors.New("github repo parsing error: unknown type")
		}
	}

	return nil
}
