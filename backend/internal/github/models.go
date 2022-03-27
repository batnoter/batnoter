package github

// GitRepoProps used to provide repo details to github client
type GitRepoProps struct {
	Repository    string
	DefaultBranch string
	Owner         string
}

// GitFileProps used to provide request details to github client
type GitFileProps struct {
	SHA         string // this is a blob sha (not commit sha)
	Path        string
	Content     string
	AuthorName  string
	AuthorEmail string
	RepoDetails GitRepoProps
}

// GitFile used to provide response to file request
type GitFile struct {
	SHA     string // this is a blob sha (not commit sha)
	Path    string
	Content string
	Size    int
	IsDir   bool
}

// GitRepo used to provide response to repos request
type GitRepo struct {
	Name          string
	Visibility    string
	DefaultBranch string
}
