package resource

import "github.com/google/go-github/v60/github"

type Version struct {
	TagName string `json:"tag_name"`
	Ref     string `json:"ref"`
}

type SourceInfo struct {
	Owner          string `json:"owner"`
	Repository     string `json:"repository"`
	AccessToken    string `json:"access_token,omitempty"`
	TagSemverRange string `json:"tag_semver_range,omitempty"`
}

type GithubTagsResourceType[S SourceInfo, V Version, G any, P any] struct {
	client *github.Client
}
