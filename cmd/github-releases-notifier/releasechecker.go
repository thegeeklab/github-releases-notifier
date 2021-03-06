package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/shurcooL/githubv4"
	"github.com/thegeeklab/github-releases-notifier/internal/model"
)

// Checker has a githubv4 client to run queries and also knows about
// the current repositories releases to compare against.
type Checker struct {
	logger   log.Logger
	client   *githubv4.Client
	releases map[string]model.Repository
}

// Run the queries and comparisons for the given repositories in a given interval.
func (c *Checker) Run(interval time.Duration, repositories []string,
	ignorePre bool, releases chan<- model.Repository) {
	if c.releases == nil {
		c.releases = make(map[string]model.Repository)
	}

	for {
		for _, repoName := range repositories {
			s := strings.Split(repoName, "/")
			owner, name := s[0], s[1]
			msg := "no new release for repository"

			nextRepo, err := c.query(owner, name)
			if err != nil {
				level.Warn(c.logger).Log(
					"msg", "failed to query the repository's releases",
					"owner", owner,
					"name", name,
					"err", err,
				)
				continue
			}

			// For debugging uncomment this next line
			//releases <- nextRepo

			currRepo, ok := c.releases[repoName]

			// We've queried the repository for the first time.
			// Saving the current state to compare with the next iteration.
			if !ok {
				c.releases[repoName] = nextRepo
				continue
			}

			if nextRepo.Release.PublishedAt.After(currRepo.Release.PublishedAt) {
				if !(ignorePre && nextRepo.Release.IsPrerelease) {
					releases <- nextRepo
					c.releases[repoName] = nextRepo
					msg = "found new release for repository"
				} else {
					msg = "ignoring new pre-release for repository"
				}
			}

			level.Debug(c.logger).Log(
				"msg", msg,
				"owner", owner,
				"name", name,
			)
		}
		time.Sleep(interval)
	}
}

// This should be improved in the future to make batch requests for all watched repositories at once
// TODO: https://github.com/shurcooL/githubv4/issues/17

func (c *Checker) query(owner, name string) (model.Repository, error) {
	var query struct {
		Repository struct {
			ID          githubv4.ID
			Name        githubv4.String
			Description githubv4.String
			URL         githubv4.URI

			Releases struct {
				Edges []struct {
					Node struct {
						ID           githubv4.ID
						Name         githubv4.String
						Description  githubv4.String
						URL          githubv4.URI
						PublishedAt  githubv4.DateTime
						IsPrerelease githubv4.Boolean
					}
				}
			} `graphql:"releases(last: 1)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.client.Query(ctx, &query, variables); err != nil {
		return model.Repository{}, err
	}

	repositoryID, ok := query.Repository.ID.(string)
	if !ok {
		return model.Repository{}, fmt.Errorf("can't convert repository id to string: %v", query.Repository.ID)
	}

	if len(query.Repository.Releases.Edges) == 0 {
		return model.Repository{}, fmt.Errorf("can't find any releases for %s/%s", owner, name)
	}
	latestRelease := query.Repository.Releases.Edges[0].Node

	releaseID, ok := latestRelease.ID.(string)
	if !ok {
		return model.Repository{}, fmt.Errorf("can't convert release id to string: %v", query.Repository.ID)
	}

	return model.Repository{
		ID:          repositoryID,
		Name:        string(query.Repository.Name),
		Owner:       owner,
		Description: string(query.Repository.Description),
		URL:         *query.Repository.URL.URL,

		Release: model.Release{
			ID:           releaseID,
			Name:         string(latestRelease.Name),
			Description:  string(latestRelease.Description),
			URL:          *latestRelease.URL.URL,
			PublishedAt:  latestRelease.PublishedAt.Time,
			IsPrerelease: bool(latestRelease.IsPrerelease),
		},
	}, nil
}
