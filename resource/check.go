package resource

import (
	"context"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v60/github"
	"github.com/jghiloni/custom-resource-type/framework"
)

func (r GithubTagsResourceType[S, V, G, P]) Check(request framework.CheckRequest[SourceInfo, Version]) ([]Version, error) {
	r.client = github.NewClient(nil)
	if request.Source.AccessToken != "" {
		r.client = r.client.WithAuthToken(request.Source.AccessToken)
	}

	opts := &github.ListOptions{PerPage: 100}
	var allTags []*github.RepositoryTag
	for {
		tags, resp, err := r.client.Repositories.ListTags(context.Background(), request.Source.Owner, request.Source.Repository, opts)
		if err != nil {
			return nil, err
		}

		allTags = append(allTags, tags...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	versions := make([]Version, len(allTags))
	for i := range allTags {
		versions[i] = Version{
			TagName: allTags[i].GetName(),
			Ref:     allTags[i].GetCommit().GetSHA(),
		}
	}

	sort.Slice(versions, func(i, j int) bool {
		vi, vj := versions[i].TagName, versions[j].TagName
		svi, svj := semver.MustParse(vi), semver.MustParse(vj)

		// sort in descending order
		return svi.GreaterThan(svj)
	})

	if request.Version != nil {
		parsed := semver.MustParse(request.Version.TagName)
		versions = sliceFilter(versions, func(v Version) bool {
			parsedV := semver.MustParse(v.TagName)
			return parsedV.GreaterThan(parsed) || parsedV.Equal(parsed)
		}, true)
	}

	if request.Source.TagSemverRange != "" {
		constraint, err := semver.NewConstraint(request.Source.TagSemverRange)
		if err != nil {
			return nil, err
		}

		versions = sliceFilter(versions, func(v Version) bool {
			parsed := semver.MustParse(v.TagName)
			return constraint.Check(parsed)
		}, false)
	}

	return versions, nil
}

func sliceFilter[T any](slice []T, filterFunc func(item T) bool, stopOnFirstFilter bool) []T {
	newSlice := make([]T, 0, len(slice))
	for i := range slice {
		if filterFunc(slice[i]) {
			newSlice = append(newSlice, slice[i])
			continue
		}

		if stopOnFirstFilter {
			break
		}
	}

	return newSlice
}
