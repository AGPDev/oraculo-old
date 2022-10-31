package app

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"strconv"
	"time"
)

func SyncJobs() {
	if lock {
		return
	}

	lock = true

	logrus.Info("=> Sincronizando jobs dos projetos no Gitlab.")

	load.Start()

	var nJobs = make(JobsT)
	var nJobsC = 0

	timeNow = time.Now()

	for i, project := range Projects {

		if project.LastJobSuccess != nil &&
			project.LastJobSuccess.FinishedAt.Day() == timeNow.Day() {
			newJobId := strconv.Itoa(project.LastJobSuccess.ID)
			nJobs[newJobId] = project.LastJobSuccess
			continue
		}

		lj, _, err := glc.Jobs.ListProjectJobs(
			project.ID,
			&gitlab.ListJobsOptions{
				ListOptions: gitlab.ListOptions{
					PerPage: 10,
					Page:    1,
				},
				Scope: []gitlab.BuildStateValue{gitlab.Success},
			},
		)
		if err != nil {
			load.Stop()
			logrus.Errorf("erro ao listar os jobs de: %s \n %s", project.Name, err)
			load.Start()
			continue
		}

		for _, j := range lj {
			if j.Stage == "dump" &&
				j.FinishedAt.Day() == timeNow.Day() &&
				j.FinishedAt.Month() == timeNow.Month() &&
				j.FinishedAt.Year() == timeNow.Year() {

				newJobId := strconv.Itoa(j.ID)
				nJobs[newJobId] = &Job{
					ID:         j.ID,
					CompanyId:  project.CompanyId,
					ProjectId:  project.ID,
					Status:     j.Status,
					Stage:      j.Stage,
					FinishedAt: j.FinishedAt.UTC().In(timeLocation),
				}

				Projects[i].LastJobSuccess = nJobs[newJobId]
				nJobsC++

				break
			}
		}
	}

	SaveProjects(&Projects)
	SaveJobs(&nJobs)

	lock = false

	load.Stop()

	logrus.Infof("=> Novos jobs encontradas: %v", nJobsC)
	logrus.Info("=> Jobs sincronizados \n\n")
}
