package app

import (
	"github.com/briandowns/spinner"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
	"time"
	"unicode"
)

var (
	timeLocation *time.Location
	timeNow      time.Time
	load         *spinner.Spinner
	dirName      string
	lock         bool
)

func Start() {
	timeLocation, _ = time.LoadLocation("America/Sao_Paulo")

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	load = spinner.New(
		spinner.CharSets[6],
		100*time.Millisecond,
		spinner.WithColor("yellow"),
	)

	load.Stop()

	// definições de configurações
	parseConfig("config.toml")

	LoadCompanies()
	LoadProjects()
	LoadJobs()
	LoadDumps()

	movideskConnect()
	gitLabConnect()

	SyncFolder()
	SyncCompanies()
	SyncProjects()
	SyncProjectsCompanies()
	SyncJobs()
	SyncDumps()
	SyncDownloads()

	c := cron.New()

	_, err := c.AddFunc(Cfg.IntervalClearToday, SyncClearToday)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncFolder, SyncFolder)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncCompanies, SyncCompanies)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncProjects, SyncProjects)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncProjectsCompanies, SyncProjectsCompanies)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncJobs, SyncJobs)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncDumps, SyncDumps)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncTickets, SyncTickets)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	_, err = c.AddFunc(Cfg.IntervalSyncDownloads, SyncDownloads)
	if err != nil {
		logrus.Panicf("%v", err)
	}

	c.Start()

	select {}
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// Remove todos os caracteres que não são Alfanumericos
//func removeNonAlphanumeric(s string) string {
//	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
//	if err != nil {
//		logrus.Errorf("não foi possível remover os caracteres da string: %s \n %s", s, err)
//	}
//
//	return reg.ReplaceAllString(s, "")
//}

// Helper para debugar variaveis
func pr(i interface{}) {
	config := spew.ConfigState{
		Indent:                  "\t",
		DisableMethods:          true,
		DisablePointerMethods:   true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
	}

	config.Dump(i)
}
