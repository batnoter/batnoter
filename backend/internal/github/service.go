package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v43/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

//go:generate mockgen -source=service.go -package=github -destination=mock_service.go
type Service interface {
	GetAuthCodeURL(state string) string
	GetToken(ctx context.Context, code string) (oauth2.Token, error)
	GetUser(ctx context.Context, ghToken oauth2.Token) (github.User, error)

	GetRepos(ctx context.Context, ghToken oauth2.Token) ([]GitRepo, error)
	SearchFiles(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps, query string, pageNo int) ([]GitFile, int, error)
	GetFile(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps) (GitFile, error)
	SaveFile(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps) (GitFile, error)
	DeleteFile(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps) error
}

type service struct {
	clientBuilder ClientBuilder
}

func NewService(clientBuilder ClientBuilder) Service {
	return &service{
		clientBuilder: clientBuilder,
	}
}

const (
	dirType       = "dir"
	commitMessage = "Created with GitNoter"
	affiliation   = "owner"
	fileExtension = "md"
	pageSize      = 20
)

func (s *service) GetAuthCodeURL(state string) string {
	// AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
	// Generate a random `state` string and validate that it matches the `state` query parameter
	// on redirect callback
	return s.clientBuilder.GetOAuth2Config().AuthCodeURL(state)
}

func (s *service) GetToken(ctx context.Context, code string) (oauth2.Token, error) {
	ghToken, err := s.clientBuilder.GetOAuth2Config().Exchange(ctx, code)
	if err != nil {
		return oauth2.Token{}, errors.Wrap(err, "retrieving user token from github failed")
	}
	return *ghToken, nil
}

func (s *service) GetUser(ctx context.Context, ghToken oauth2.Token) (github.User, error) {
	client := s.clientBuilder.Build(ctx, &ghToken)

	githubUser, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return github.User{}, errors.Wrap(err, "retrieving user from github failed")
	}
	if githubUser.Email == nil {
		return github.User{}, errors.New("processing github user object failed")
	}
	return *githubUser, nil
}

func (s *service) GetRepos(ctx context.Context, ghToken oauth2.Token) ([]GitRepo, error) {
	client := s.clientBuilder.Build(ctx, &ghToken)
	opts := &github.RepositoryListOptions{
		Affiliation: affiliation,
	}
	gitRepos, _, err := client.Repositories.List(ctx, "", opts)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving user's repos from github failed")
	}
	repos := make([]GitRepo, 0, len(gitRepos))
	for _, gitRepo := range gitRepos {
		repos = append(repos, GitRepo{
			Name:          gitRepo.GetName(),
			Visibility:    gitRepo.GetVisibility(),
			DefaultBranch: gitRepo.GetDefaultBranch(),
		})
	}
	return repos, nil
}

func (s *service) SearchFiles(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps, query string, pageNo int) ([]GitFile, int, error) {
	client := s.clientBuilder.Build(ctx, &ghToken)

	opts := &github.SearchOptions{
		TextMatch: true,
		ListOptions: github.ListOptions{
			Page:    pageNo,
			PerPage: pageSize,
		},
	}
	pathQualifier := ""
	if fileProps.Path != "" {
		pathQualifier = "path:" + fileProps.Path
	}
	ghQuery := fmt.Sprintf("%s %s extension:md repo:%s/%s", query, pathQualifier, fileProps.RepoDetails.Owner, fileProps.RepoDetails.Repository)
	cs, _, err := client.Search.Code(ctx, ghQuery, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, "searching on github failed")
	}
	gitFiles := make([]GitFile, 0, len(cs.CodeResults))
	for _, item := range cs.CodeResults {
		gitFile, err := s.getFileInternal(ctx, client, fileProps.RepoDetails.Owner, fileProps.RepoDetails.Repository, fileProps.RepoDetails.DefaultBranch, item.GetPath())
		if err != nil {
			return nil, 0, err
		}
		gitFiles = append(gitFiles, gitFile)
	}
	return gitFiles, cs.GetTotal(), nil
}

func (s *service) GetFile(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps) (GitFile, error) {
	client := s.clientBuilder.Build(ctx, &ghToken)

	gitFile, err := s.getFileInternal(ctx, client, fileProps.RepoDetails.Owner, fileProps.RepoDetails.Repository, fileProps.RepoDetails.DefaultBranch, fileProps.Path)
	if err != nil {
		return GitFile{}, err
	}

	return gitFile, nil
}

func (*service) getFileInternal(ctx context.Context, client *github.Client, owner string, repo string, branch string, path string) (GitFile, error) {
	opts := &github.RepositoryContentGetOptions{
		Ref: branch,
	}

	fc, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opts)
	if err != nil {
		return GitFile{}, errors.Wrap(err, "retrieving file from github failed")
	}

	if fc == nil {
		return GitFile{}, errors.New("file with matching path not found. retrieving file from github failed")
	}

	contents, err := fc.GetContent()
	if err != nil {
		return GitFile{}, errors.Wrap(err, "parsing github file content failed")
	}

	gitFile := GitFile{
		SHA:     fc.GetSHA(),
		IsDir:   fc.GetType() == dirType,
		Content: contents,
		Size:    fc.GetSize(),
		Path:    fc.GetPath(),
	}
	return gitFile, nil
}

func (s *service) SaveFile(ctx context.Context, ghToken oauth2.Token, fp GitFileProps) (GitFile, error) {
	client := s.clientBuilder.Build(ctx, &ghToken)
	fileContent := []byte(fp.Content)

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(commitMessage),
		Content:   fileContent,
		Branch:    github.String(fp.RepoDetails.DefaultBranch),
		Committer: &github.CommitAuthor{Name: github.String(fp.AuthorName), Email: github.String(fp.AuthorEmail)},
	}
	if fp.SHA != "" {
		// providing the blob sha will update the file on github otherwise new file is created
		opts.SHA = &fp.SHA
	}
	rc, _, err := client.Repositories.UpdateFile(ctx, fp.RepoDetails.Owner, fp.RepoDetails.Repository, fp.Path, opts)
	if err != nil {
		return GitFile{}, errors.Wrap(err, "saving file to github failed")
	}
	return GitFile{
		SHA:   rc.Content.GetSHA(),
		IsDir: false,
		Size:  rc.Content.GetSize(),
		Path:  rc.Content.GetPath(),
	}, nil
}

func (s *service) DeleteFile(ctx context.Context, ghToken oauth2.Token, fileProps GitFileProps) error {
	client := s.clientBuilder.Build(ctx, &ghToken)
	fileContent := []byte(fileProps.Content)

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(commitMessage),
		Content:   fileContent,
		Branch:    github.String(fileProps.RepoDetails.DefaultBranch),
		Committer: &github.CommitAuthor{Name: github.String(fileProps.AuthorName), Email: github.String(fileProps.AuthorEmail)},
		SHA:       &fileProps.SHA,
	}
	_, _, err := client.Repositories.DeleteFile(ctx, fileProps.RepoDetails.Owner, fileProps.RepoDetails.Repository, fileProps.Path, opts)
	if err != nil {
		return errors.Wrap(err, "deleting file from github failed")
	}
	return nil
}
