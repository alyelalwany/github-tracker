package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os/exec"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type reposQuery struct {
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

func GetRepoForOwner(w http.ResponseWriter, r *http.Request) {

	username := r.PathValue("username")

	variables := map[string]any{
		"login": githubv4.String(username),
		"after": (*githubv4.String)(nil),
	}

	var query reposQuery
	output, err := exec.Command("gh", "auth", "token").Output()

	if err != nil {
		slog.Error("Failed to get token", "error", err)
		return
	}

	token := strings.TrimSpace(string(output))

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		slog.Error("Failed to make request", "error", err)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}

	slog.Info("Login:", query.User)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(query.User.Repositories); err != nil {
		slog.Error("Failed to encode the response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
