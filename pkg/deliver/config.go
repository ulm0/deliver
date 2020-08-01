package deliver

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/urfave/cli/v2"
)

// Config contains configuration regarding repo release
type Config struct {
	APIToken        string
	APIURL          string
	Checksum        cli.StringSlice
	ChecksumFile    string
	ChecksumFlatten bool
	Draft           bool
	FileExists      string
	Files           cli.StringSlice
	Note            string
	Override        bool
	Prerelease      bool
	Repo            string
	Tag             string
	Title           string
	UploadURL       string

	apiURL    *url.URL
	uploadURL *url.URL
	uploads   []string
}

type deliverer struct {
	*github.Client
	context.Context
	Draft      bool
	FileExists string
	Note       string
	Overwrite  bool
	Owner      string
	Prerelease bool
	Repo       string
	Tag        string
	Title      string
}

// Check checks if everything is as supposed
func (c *Config) Check() error {
	var err error
	repoTag := os.Getenv("CI_COMMIT_TAG")
	repoName := os.Getenv("CI_PROJECT_NAME")
	if repoTag == "" {
		return fmt.Errorf("deliver only works on tag refs")
	}
	c.Tag = repoTag
	if repoName == "" {
		return fmt.Errorf("Repo name not set")
	}
	c.Repo = repoName
	if c.APIToken == "" {
		return fmt.Errorf("API key not set")
	}
	if !fileExistsValues[c.FileExists] {
		return fmt.Errorf("Invalid value for file-exists")
	}
	if !strings.HasSuffix(c.APIURL, "/") {
		c.APIURL = c.APIURL + "/"
	}
	c.apiURL, err = url.Parse(c.APIURL)
	if err != nil {
		return fmt.Errorf("Failed to parse api URL: %w", err)

	}
	if !strings.HasSuffix(c.UploadURL, "/") {
		c.UploadURL = c.UploadURL + "/"
	}
	c.uploadURL, err = url.Parse(c.UploadURL)
	if err != nil {
		return fmt.Errorf("Failed to parse upload URL: %w", err)
	}
	if c.Note != "" {
		if c.Note, err = readStringOrFile(c.Note); err != nil {
			return fmt.Errorf("Error while reading %s: %w", c.Note, err)
		}
	}
	files := c.Files.Value()
	for _, glob := range files {
		g, err := filepath.Glob(glob)
		if err != nil {
			return fmt.Errorf("Failed to glob %s: %w", glob, err)
		}
		if g != nil {
			c.uploads = append(c.uploads, g...)
		}
	}

	if len(files) > 0 && len(c.uploads) < 1 {
		return fmt.Errorf("Failed to find any file to release")
	}

	checksum := c.Checksum.Value()
	if len(checksum) > 0 {
		c.uploads, err = writeChecksums(c.uploads, checksum, c.ChecksumFile, c.ChecksumFlatten)
		if err != nil {
			return fmt.Errorf("Failed to write checksums: %w", err)
		}
	}
	return nil
}
