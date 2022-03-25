package github

type GitRepoProps struct {
	Repository    string
	DefaultBranch string
	Owner         string
}

// used to provide request details to github client
type GitFileProps struct {
	SHA         string // this is a blob sha (not commit sha)
	Path        string
	Content     string
	AuthorName  string
	AuthorEmail string
	RepoDetails GitRepoProps
}

// used to provide response to file request
type GitFile struct {
	SHA     string // this is a blob sha (not commit sha)
	Path    string
	Content string
	Size    int
	IsDir   bool
}

// used to provide response to repos request
type GitRepo struct {
	Name          string
	Visibility    string
	DefaultBranch string
}
