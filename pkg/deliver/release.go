package deliver

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/google/go-github/v32/github"
	oauth "golang.org/x/oauth2"
)

// Execute releases to Github
func (c *Config) Execute() error {
	ctx := context.Background()
	ts := oauth.StaticTokenSource(&oauth.Token{AccessToken: c.APIToken})
	tc := oauth.NewClient(ctx, ts)

	ghClient := github.NewClient(tc)
	ghClient.BaseURL = c.apiURL
	ghClient.UploadURL = c.uploadURL

	// ListOptions set to 200 because the default value can be too short sometimes
	repos, _, err := ghClient.Repositories.List(ctx, "", &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 200}})
	if err != nil {
		return err
	}

	var repoOwner string
	for _, repo := range repos {
		if *repo.Name == c.Repo {
			repoOwner = *repo.Owner.Login
		}
	}

	dl := deliverer{
		Client:     ghClient,
		Context:    ctx,
		Owner:      repoOwner,
		Repo:       c.Repo,
		Tag:        c.Tag,
		Draft:      c.Draft,
		Prerelease: c.Prerelease,
		FileExists: c.FileExists,
		Title:      c.Title,
		Note:       c.Note,
		Overwrite:  c.Override,
	}

	release, err := dl.setRelease()

	if err != nil {
		return fmt.Errorf("Failed to create release: %w", err)
	}

	if err := dl.uploadFiles(*release.ID, c.uploads); err != nil {
		return fmt.Errorf("Failed to upload release files: %w", err)
	}

	return nil
}

func (d *deliverer) setRelease() (*github.RepositoryRelease, error) {
	release, err := d.getRelease()

	if err != nil && release == nil {
		fmt.Println(err)
		release, err = d.newRelease()
	} else if release != nil && d.Overwrite {
		release, err = d.editRelease(*release.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve or create a release: %w", err)
	}

	return release, nil
}

func (d *deliverer) getRelease() (*github.RepositoryRelease, error) {
	release, _, err := d.Client.Repositories.GetReleaseByTag(d.Context, d.Owner, d.Repo, d.Tag)

	if err != nil {
		return nil, fmt.Errorf("Release %s not found", d.Tag)
	}

	fmt.Printf("Succesfully retrieved %s release\n", d.Tag)

	return release, nil
}

func (d *deliverer) editRelease(releaseID int64) (*github.RepositoryRelease, error) {
	repoRelease := &github.RepositoryRelease{
		Name: &d.Title,
		Body: &d.Note,
	}

	release, _, err := d.Client.Repositories.EditRelease(d.Context, d.Owner, d.Repo, releaseID, repoRelease)

	if err != nil {
		return nil, fmt.Errorf("Failed to update release: %w", err)
	}

	fmt.Printf("Successfully updated %s release\n", d.Tag)
	return release, nil
}

func (d *deliverer) newRelease() (*github.RepositoryRelease, error) {
	repoRelease := &github.RepositoryRelease{
		Body:       &d.Note,
		Draft:      &d.Draft,
		Name:       &d.Title,
		Prerelease: &d.Prerelease,
		TagName:    github.String(d.Tag),
	}

	release, _, err := d.Client.Repositories.CreateRelease(d.Context, d.Owner, d.Repo, repoRelease)

	if err != nil {
		return nil, fmt.Errorf("Failed to create release: %w", err)
	}

	fmt.Printf("Successfully create %s release\n", d.Tag)
	return release, nil
}

func (d *deliverer) uploadFiles(id int64, files []string) error {
	assets, _, err := d.Client.Repositories.ListReleaseAssets(d.Context, d.Owner, d.Repo, id, &github.ListOptions{})

	if err != nil {
		return fmt.Errorf("Failed to fecth existing assets: %w", err)
	}

	var uploadFiles []string

files:
	for _, file := range files {
		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				switch d.FileExists {
				case "overwrite":
					// do nothing
				case "fail":
					return fmt.Errorf("Asset file %s already exists", path.Base(file))
				case "skip":
					fmt.Printf("Skipping pre-existing %s artifact\n", *asset.Name)
					continue files
				default:
					return fmt.Errorf("Internal error, unknown file_exists value %s", d.FileExists)
				}
			}
		}

		uploadFiles = append(uploadFiles, file)
	}

	for _, file := range uploadFiles {
		handle, err := os.Open(file)

		if err != nil {
			return fmt.Errorf("Failed to read %s artifact: %w", file, err)
		}

		for _, asset := range assets {
			if *asset.Name == path.Base(file) {
				if _, err := d.Client.Repositories.DeleteReleaseAsset(d.Context, d.Owner, d.Repo, *asset.ID); err != nil {
					return fmt.Errorf("Failed to delete %s artifact: %w", file, err)
				}

				fmt.Printf("Successfully deleted old %s artifact\n", *asset.Name)
			}
		}

		upOpts := &github.UploadOptions{Name: path.Base(file)}

		if _, _, err = d.Client.Repositories.UploadReleaseAsset(d.Context, d.Owner, d.Repo, id, upOpts, handle); err != nil {
			return fmt.Errorf("Failed to upload %s artifact: %w", file, err)
		}

		fmt.Printf("Successfully uploaded %s artifact\n", file)
	}

	return nil
}
