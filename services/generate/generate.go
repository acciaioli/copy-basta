package generate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"copy-basta/internal/clients/github"
	"copy-basta/internal/common"
	"copy-basta/internal/common/log"
	"copy-basta/internal/crawl"
	"copy-basta/internal/load"
	"copy-basta/internal/specification"
	"copy-basta/internal/write"
)

type Params struct {
	Src       string
	Dest      string
	SpecYAML  string
	InputYAML string
}

func Generate(params *Params) error {
	log.L.DebugWithData("params", log.Data{
		"src":       params.Src,
		"dest":      params.Dest,
		"specYAML":  params.SpecYAML,
		"inputYAML": params.InputYAML,
	})

	log.L.Info("validating params...")
	err := validate(params)
	if err != nil {
		return err
	}
	log.L.Info("params are valid!")

	log.L.Info("crawling files...")
	crawler, err := getCrawler(params.Src)
	if err != nil {
		return err
	}
	crawledFiles, err := crawler.Crawl()
	if err != nil {
		return err
	}
	log.L.Info("files crawled!")

	log.L.Info("loading specification...")
	specLoadedPath := common.TrimRootDir(filepath.Join(params.Src, params.SpecYAML))
	spec, err := specification.New(specLoadedPath, crawledFiles)
	if err != nil {
		return err
	}
	log.L.Info("spec loaded!")

	log.L.Info("loading files...")
	loader, err := load.New(spec.Ignorer, spec.Passer)
	if err != nil {
		return err
	}
	files, err := loader.Load(crawledFiles)
	if err != nil {
		return err
	}
	{
		logData := log.Data{}
		for _, f := range files {
			logData[f.Path] = fmt.Sprintf("mode=%v, is-template=%T, byte-counts=%d", f.Mode, f.Template, len(f.Content))
		}
		log.L.DebugWithData("loaded files", logData)
	}
	log.L.Info("files loaded!")

	var input common.InputVariables
	if params.InputYAML != "" {
		log.L.InfoWithData("loading template variables from file", log.Data{"location": params.InputYAML})
		fileInput, err := spec.Variables.InputFromFile(params.InputYAML)
		if err != nil {
			return err
		}
		input = fileInput
	} else {
		log.L.Info("getting template variables dynamically")
		stdinInput, err := spec.Variables.InputFromStdIn()
		if err != nil {
			return err
		}
		input = stdinInput
	}

	log.L.InfoWithData("writing to new project", log.Data{"location": params.Dest})
	writer := write.NewDiskWriter(params.Dest)
	err = writer.Write(files, input)
	if err != nil {
		return err
	}

	log.L.Info("done!")
	return nil
}

func validate(params *Params) error {
	if params.Src == "" {
		return errors.New("params validation error - src can't be empty")
	}

	if strings.HasPrefix(params.Src, common.GithubPrefix) {
		log.L.Warn("src is a remote location, skipping validations...")
		return nil
	}

	// src
	{
		_, err := os.Stat(params.Src)
		if os.IsNotExist(err) {
			return fmt.Errorf("params validation error - src directory (%s) not found", params.Src)
		}
		if err != nil {
			return err
		}

	}

	// specYAML
	{
		if params.SpecYAML == "" {
			return fmt.Errorf("params validation error - specYAML can't be empty")
		}
		specYAMLFullPath := filepath.Join(params.Src, params.SpecYAML)
		fInfo, err := os.Stat(specYAMLFullPath)
		if err != nil {
			return fmt.Errorf("params validation error - specYAML file (%s) not found", specYAMLFullPath)
		}
		if fInfo.IsDir() {
			return fmt.Errorf("params validation error - specYAML file (%s) is not a file", specYAMLFullPath)
		}
	}

	// inputYAML
	{
		if params.InputYAML != "" {
			fInfo, err := os.Stat(params.InputYAML)
			if err != nil {
				return fmt.Errorf("params validation error - inputYAML file (%s) not found", params.InputYAML)
			}
			if fInfo.IsDir() {
				return fmt.Errorf("params validation error - inputYAML file (%s) is not a file", params.InputYAML)
			}
		}
	}

	return nil
}

func getCrawler(src string) (crawl.Crawler, error) {
	switch {
	case strings.HasPrefix(src, common.GithubPrefix):
		log.L.Debug("using github crawler")
		ghc, err := github.NewClient(strings.TrimPrefix(src, common.GithubPrefix))
		if err != nil {
			return nil, err
		}
		return crawl.NewGithubCrawler(ghc), nil
	default:
		log.L.Debug("using disk crawler")
		return crawl.NewLocalCrawler(src), nil
	}
}
