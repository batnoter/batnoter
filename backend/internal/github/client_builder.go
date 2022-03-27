package github

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

// ClientBuilder represents an oauth2 github client builder.
// It provides methods to build the github oauth2 client.
//go:generate mockgen -source=client_builder.go -package=github -destination=mock_client_builder.go
type ClientBuilder interface {
	Build(ctx context.Context, token *oauth2.Token) *github.Client
	GetOAuth2Config() *oauth2.Config
}

type clientBuilder struct {
	oauth2Config *oauth2.Config
}

// NewClientBuilder creates and returns a new oauth2 client builder containing oauth2 config.
func NewClientBuilder(oauth2Config *oauth2.Config) ClientBuilder {
	return &clientBuilder{
		oauth2Config: oauth2Config,
	}
}

// Build creates and returns a github oauth2 client using oauth2 token.
func (c *clientBuilder) Build(ctx context.Context, token *oauth2.Token) *github.Client {
	return github.NewClient(c.oauth2Config.Client(ctx, token))
}

// GetOAuth2Config returns oauth2 config.
func (c *clientBuilder) GetOAuth2Config() *oauth2.Config {
	return c.oauth2Config
}
