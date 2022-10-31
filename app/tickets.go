package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/resty.v1"
	"strconv"
)

type (
	Organization struct {
		Id           string `json:"id"`
		BusinessName string `json:"businessName"`
	}

	Client struct {
		Organization Organization `json:"organization"`
	}

	Ticket struct {
		Clients []Client `json:"clients"`
	}
)

func SyncTickets() {
	if lock {
		return
	}

	lock = true

	logrus.Info("=> Sincronizando tickets novos e em andamento")

	load.Start()

	var companyIds = make(map[string]string)
	var projectsIds = make(map[string]string)
	var tickets []Ticket
	var skip = 0

	r := resty.R()
	r.SetQueryParam("token", Cfg.MovideskToken)
	r.SetQueryParam("$filter", "status eq 'Novo' or status eq 'Em atendimento'")
	r.SetQueryParam("$select", "clients")
	r.SetQueryParam("$expand", "clients($expand=organization($select=id,businessName))")

	for {

		r.SetQueryParam("$skip", strconv.Itoa(skip))

		resp, err := r.Get("/tickets")
		if err != nil {
			logrus.Errorf("erro ao listar os tickets do Movidesk \n %s", err)
			return
		}

		err = json.Unmarshal(resp.Body(), &tickets)
		if err != nil {
			logrus.Errorf("erro ao parsear os tickets \n %v", err)
			return
		}

		ticketsLen := len(tickets)
		if ticketsLen <= 0 {
			break
		} else {
			skip += ticketsLen
		}

		for _, ticket := range tickets {
			for _, client := range ticket.Clients {
				companyIds[client.Organization.Id] = client.Organization.Id
			}
		}
	}

	for _, companyId := range companyIds {
		for projectId, project := range Projects {
			if project.CompanyId == companyId {
				projectsIds[projectId] = projectId
				break
			}
		}
	}

	syncNewJobs(projectsIds)

	lock = false

	load.Stop()

	logrus.Info("=> Tickets sincronizados \n\n")
}

func syncNewJobs(projectsIds map[string]string) {
	for _, projectId := range projectsIds {
		if Projects[projectId] == nil && Projects[projectId].LastJobSuccess != nil {
			continue
		}

		lj, _, err := glc.Jobs.ListProjectJobs(
			projectId,
			&gitlab.ListJobsOptions{
				ListOptions: gitlab.ListOptions{
					PerPage: 10,
					Page:    1,
				},
				Scope: []gitlab.BuildStateValue{BuildStateManual},
			},
		)
		if err != nil {
			logrus.Errorf("erro ao listar os jobs manuais \n %s", err)
			continue
		}

		for _, j := range lj {
			if j.Stage == "dump" {
				_, _, err := glc.Jobs.PlayJob(projectId, j.ID)
				if err != nil {
					logrus.Errorf("erro ao executar o job %s \n %s", j.ID, err)
				}
				break
			}
		}
	}
}
