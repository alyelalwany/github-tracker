package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shurcooL/githubv4"
)

type userDetailsQuery struct {
	User struct {
		Login           githubv4.String
		Name            githubv4.String
		Email           githubv4.String
		Bio             githubv4.String
		AvatarURL       githubv4.URI `graphql:"avatarUrl"`
		URL             githubv4.URI
		Company         githubv4.String
		Location        githubv4.String
		WebsiteURL      githubv4.URI `graphql:"websiteUrl"`
		TwitterUsername githubv4.String
		CreatedAt       githubv4.DateTime
		UpdatedAt       githubv4.DateTime
		Followers       struct {
			TotalCount githubv4.Int
		}
		Following struct {
			TotalCount githubv4.Int
		}
		Repositories struct {
			TotalCount githubv4.Int
		} `graphql:"repositories(privacy: PUBLIC)"`
	} `graphql:"user(login: $login)"`
}

func (s *Server) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	login := r.PathValue("login")
	if login == "" {
		http.Error(w, "login required", http.StatusBadRequest)
		return
	}

	vars := map[string]interface{}{
		"login": githubv4.String(login),
	}

	var query userDetailsQuery
	if err := s.gh.Query(r.Context(), &query, vars); err != nil {
		slog.Error("query failed", "error", err)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}

	userDetails := query.User

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userDetails); err != nil {
		slog.Error("encode failed", "error", err)
		http.Error(w, "Encoding failed", http.StatusBadGateway)
	}

}
