package cli

import (
	"github.com/ulm0/deliver/pkg/deliver"
	"github.com/urfave/cli/v2"
)

func configFlags(config *deliver.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "api-key",
			Aliases:     []string{"k"},
			EnvVars:     []string{"GITHUB_TOKEN", "GITHUB_API_TOKEN"},
			Usage:       "API token to access github",
			Destination: &config.APIToken,
		},
		&cli.StringSliceFlag{
			Name:        "files",
			Aliases:     []string{"f"},
			EnvVars:     []string{"RELEASE_FILES"},
			Usage:       "List of files to release",
			Destination: &config.Files,
		},
		&cli.StringFlag{
			Name:        "file-exists",
			EnvVars:     []string{"RELEASE_FILE_EXISTS"},
			Value:       "overwrite",
			Usage:       "Behavior in case a file previously exists",
			Destination: &config.FileExists,
		},
		&cli.StringSliceFlag{
			Name:        "checksum",
			Aliases:     []string{"s"},
			EnvVars:     []string{"RELEASE_CHECKSUM"},
			Usage:       "Methods for generating files checksums",
			Destination: &config.Checksum,
		},
		&cli.StringFlag{
			Name:        "checksum-file",
			EnvVars:     []string{"CHECKSUM_FILE"},
			Usage:       "Name for checksum file. \"CHECKSUM\" is replaced with chosen method",
			Value:       "CHECKSUMsum.txt",
			Destination: &config.ChecksumFile,
		},
		&cli.BoolFlag{
			Name:        "checksum-flatten",
			Aliases:     []string{"l"},
			EnvVars:     []string{"CHECKSUM_FLATTEN"},
			Usage:       "Include only the basename of the file in the checksum file",
			Value:       true,
			Destination: &config.ChecksumFlatten,
		},
		&cli.BoolFlag{
			Name:        "draft",
			Aliases:     []string{"t"},
			EnvVars:     []string{"DRAFT_RELEASE"},
			Usage:       "This is a draft release",
			Destination: &config.Draft,
		},
		&cli.BoolFlag{
			Name:        "pre-release",
			Aliases:     []string{"p"},
			EnvVars:     []string{"PRE_RELEASE"},
			Usage:       "This is a pre-release",
			Destination: &config.Prerelease,
		},
		&cli.StringFlag{
			Name:        "api-url",
			Usage:       "API endpoint",
			Aliases:     []string{"a"},
			EnvVars:     []string{"GITHUB_API_URL"},
			Value:       "https://api.github.com/",
			Destination: &config.APIURL,
		},
		&cli.StringFlag{
			Name:        "upload-url",
			Aliases:     []string{"u"},
			EnvVars:     []string{"GITHUB_UPLOAD_URL"},
			Usage:       "API endpoint for uploading assets",
			Value:       "https://uploads.github.com/",
			Destination: &config.UploadURL,
		},
		&cli.StringFlag{
			Name:        "name",
			Aliases:     []string{"n"},
			EnvVars:     []string{"RELEASE_NAME"},
			Usage:       "Name for this release",
			Destination: &config.Title,
		},
		&cli.StringFlag{
			Name:        "description",
			Aliases:     []string{"d"},
			EnvVars:     []string{"RELEASE_DESCRIPTION"},
			Usage:       "File or string containing release description",
			Destination: &config.Note,
		},
		&cli.BoolFlag{
			Name:        "override",
			Aliases:     []string{"o"},
			EnvVars:     []string{"OVERRIDE_RELEASE"},
			Usage:       "Override existing release information",
			Destination: &config.Override,
		},
	}
}
