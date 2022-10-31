package app

import "github.com/xanzy/go-gitlab"

var (
	BuildStateManual gitlab.BuildStateValue = "manual"
	glc              *gitlab.Client
)

func gitLabConnect() {
	glc = gitlab.NewClient(nil, Cfg.GitlabToken)
}
