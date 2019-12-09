package internal

type ProwJobSpec struct {
	Type string   `json:"type"`
	Job  string   `json:"job"`
	Refs ProwRefs `json:"refs"`
}

type ProwRefs struct {
	Org      string      `json:"org"`
	Repo     string      `json:"repo"`
	RepoLink string      `json:"repo_link"`
	BaseRef  string      `json:"base_ref"`
	Pulls    []ProwPulls `json:"pulls"`
}

type ProwPulls struct {
	Number     int    `json:"number"`
	Author     string `json:"author"`
	Link       string `json:"link"`
	AuthorLink string `json:"author_link"`
}
