package api

type CommitStatus struct {
	State string `json:"state"`
	Key   string `json:"key"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type DataMiddleWare struct {
	CommitStatus map[string]CommitStatus
}

type Repository struct {
	Slug          string `json:"slug"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ScmID         string `json:"scmId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Forkable      bool   `json:"forkable"`
}
