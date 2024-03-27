package resource

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/go-github/v60/github"
	"github.com/jghiloni/custom-resource-type/framework"
)

func (r *GithubTagsResourceType[S, V, G, P]) Get(baseDir string, request framework.GetRequest[SourceInfo, Version, any]) (framework.Response[Version], error) {
	var err error
	if err = os.WriteFile(filepath.Join(baseDir, "tag"), []byte(request.Version.TagName), 0o666); err != nil {
		return framework.Response[Version]{}, err
	}

	if err = os.WriteFile(filepath.Join(baseDir, "ref"), []byte(request.Version.Ref), 0o666); err != nil {
		return framework.Response[Version]{}, err
	}

	r.client = github.NewClient(nil)
	if request.Source.AccessToken != "" {
		r.client = r.client.WithAuthToken(request.Source.AccessToken)
	}

	commit, _, err := r.client.Repositories.GetCommit(context.Background(), request.Source.Owner, request.Source.Repository, request.Version.Ref, nil)
	if err != nil {
		return framework.Response[Version]{}, err
	}

	commitInfo, err := os.Create(filepath.Join(baseDir, "commit-info"))
	if err != nil {
		return framework.Response[Version]{}, err
	}
	defer commitInfo.Close()

	encoder := json.NewEncoder(commitInfo)
	encoder.SetIndent("", "    ")
	if err = encoder.Encode(commit); err != nil {
		return framework.Response[Version]{}, err
	}

	return framework.Response[Version]{
		Version: request.Version,
		Metadata: []framework.MetadataField{
			{
				Name:  "commit-author",
				Value: commit.GetAuthor().GetName(),
			},
		},
	}, nil
}
