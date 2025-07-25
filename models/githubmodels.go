package models

type GithubUser struct {
	Login string `json:"login"`
}

type GithubRepo struct {
	Name            string     `json:"name"`
	FullName        string     `json:"full_name"`
	HTMLURL         string     `json:"html_url"`
	UpdatedAt       string     `json:"updated_at"`
	Owner           GithubUser `json:"owner"`
	CommitCount     int        `json:"commit_count,omitempty"`
	IssueFixCount   int        `json:"issue_fix_count,omitempty"`
	StargazersCount int        `json:"stargazers_count,omitempty"`
	ForksCount      int        `json:"forks_count,omitempty"`
	WatchersCount   int        `json:"watchers_count,omitempty"`
	ReviewedPRs     int        `json:"reviewed_prs,omitempty"`
	ReviewCount     int        `json:"review_count,omitempty"`
}

type PullRequest struct {
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	MergedAt     string     `json:"merged_at"`
	CreatedAt    string     `json:"created_at"`
	State        string     `json:"state"`
	Number       int        `json:"number"`
	Additions    int        `json:"additions"`
	Deletions    int        `json:"deletions"`
	ChangedFiles int        `json:"changed_files"`
	User         GithubUser `json:"user"`
}
