package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

var Cfg Config

type (
	Duration struct {
		time.Duration
	}

	Config struct {
		IntervalSyncCompanies         string
		IntervalSyncProjects          string
		IntervalSyncProjectsCompanies string
		IntervalSyncJobs              string
		IntervalSyncDumps             string
		IntervalSyncTickets           string
		IntervalSyncFolder            string
		IntervalSyncDownloads         string
		IntervalClearToday            string
		GitlabToken                   string
		MovideskUrl                   string
		MovideskToken                 string
		MovideskDebug                 bool
	}
)

func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

func parseConfig(filename string) {
	logrus.Info("=> Definindo as configurações...")

	// Parse the raw TOML file.
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Fatalf("não foi possível ler o arquivo de configuração %s \n %s", filename, err)
	}
	if _, err = toml.Decode(string(raw), &Cfg); err != nil {
		logrus.Fatalf("não foi possível decodificar o arquivo de configuração %s, %s", filename, err)
	}

	if Cfg.IntervalSyncCompanies == "" {
		Cfg.IntervalSyncCompanies = "0 6 * * 1-5"
	}
	if Cfg.IntervalSyncProjects == "" {
		Cfg.IntervalSyncProjects = "0 7,12,18 * * 1-5"
	}
	if Cfg.IntervalSyncProjectsCompanies == "" {
		Cfg.IntervalSyncProjectsCompanies = "15 7,12,18 * * 1-5"
	}
	if Cfg.IntervalSyncJobs == "" {
		Cfg.IntervalSyncJobs = "*/10 7-19 * * 1-5"
	}
	if Cfg.IntervalSyncDumps == "" {
		Cfg.IntervalSyncDumps = "*/5 7-19 * * 1-5"
	}
	if Cfg.IntervalSyncTickets == "" {
		Cfg.IntervalSyncTickets = "*/19 6-18 * * 1-5"
	}
	if Cfg.IntervalSyncFolder == "" {
		Cfg.IntervalSyncFolder = "0 0 * * 1-5"
	}
	if Cfg.IntervalSyncDownloads == "" {
		Cfg.IntervalSyncDownloads = "*/5 7-19 * * 1-5"
	}
	if Cfg.IntervalClearToday == "" {
		Cfg.IntervalClearToday = "*/5 7-19 * * 1-5"
	}
	if Cfg.GitlabToken == "" {
		logrus.Fatalf("token do Gitlab não informado, %v", err)
	}
	if Cfg.MovideskUrl == "" {
		Cfg.MovideskUrl = "https://api.movidesk.com/public/v1"
	}
	if Cfg.MovideskToken == "" {
		err = fmt.Errorf("token de acesso a API do Movidesk não informado, %s", err)
	}
}
