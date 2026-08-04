package handler

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Server struct {
	gh *githubv4.Client
}

func NewServer() *Server {
	output, err := exec.Command("gh", "auth", "token").Output()
	if err != nil {
		log.Fatal("failed to get github token: ", err)
	}
	token := strings.TrimSpace(string(output))

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)
	return &Server{gh: githubv4.NewClient(httpClient)}
}
