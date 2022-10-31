package app

import (
	"gopkg.in/resty.v1"
)

type MovideskQuery struct {
	Filter  string
	Select  string
	Expand  string
	OrderBy string
	Skip    int
}

func movideskConnect() {
	resty.SetDebug(Cfg.MovideskDebug)
	resty.SetHostURL(Cfg.MovideskUrl)
}