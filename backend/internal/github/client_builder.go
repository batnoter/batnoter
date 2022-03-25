package github

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

//go:generate mockgen -source=client_builder.go -package=github -destination=mock_client_builder.go
type ClientBuilder interface {
	Build(ctx context.Context, token *oauth2.Token) *github.Client
	GetOAuth2Config() *oauth2.Config
}

type clientBuilder struct {
	oauth2Config *oauth2.Config
}

func NewClientBuilder(oauth2Config *oauth2.Config) ClientBuilder {
	return &clientBuilder{
		oauth2Config: oauth2Config,
	}
}

func (c *clientBuilder) Build(ctx context.Context, token *oauth2.Token) *github.Client {
	return github.NewClient(c.oauth2Config.Client(ctx, token))
}

func (c *clientBuilder) GetOAuth2Config() *oauth2.Config {
	return c.oauth2Config
}
