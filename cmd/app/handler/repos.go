package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shurcooL/githubv4"
)

type orgReposQuery struct {
	Organization struct {
		Repositories struct {
			TotalCount int
			PageInfo   struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
			Nodes []struct {
				Name           githubv4.String
				StargazerCount githubv4.Int
				CreatedAt      githubv4.DateTime
				Visibility     githubv4.RepositoryVisibility
				ForkCount      githubv4.Int
			}
		} `graphql:"repositories(first: 100, after: $after, orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"organization(login: $login)"`
}

type userReposQuery struct {
	User struct {
		CreatedAt    string
		Repositories struct {
			TotalCount int
			PageInfo   struct {
				HasNextPage bool
				EndCursor   githubv4.String
			}
			Nodes []struct {
				Name           githubv4.String
				StargazerCount githubv4.Int
				CreatedAt      githubv4.DateTime
				Visibility     githubv4.RepositoryVisibility
				ForkCount      githubv4.Int
			}
		} `graphql:"repositories(first: 100, after: $after, ownerAffiliations: [OWNER], orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"user(login: $login)"`
}

func (s *Server) GetRepos(w http.ResponseWriter, r *http.Request) {
	kind := r.PathValue("kind") // "user" or "org"
	login := r.PathValue("login")
	if login == "" {
		http.Error(w, "login required", http.StatusBadRequest)
		return
	}

	vars := map[string]any{
		"login": githubv4.String(login),
		"after": (*githubv4.String)(nil),
	}

	var repos any
	switch kind {
	case "user":
		var q userReposQuery
		if err := s.gh.Query(r.Context(), &q, vars); err != nil {
			slog.Error("query failed", "error", err)
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		repos = q.User.Repositories
	case "org":
		var q orgReposQuery
		if err := s.gh.Query(r.Context(), &q, vars); err != nil {
			slog.Error("query failed", "error", err)
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}
		repos = q.Organization.Repositories
	default:
		http.Error(w, "kind must be user or org", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(repos); err != nil {
		slog.Error("encode failed", "error", err)
	}
}
