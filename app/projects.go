package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gopkg.in/toqueteos/substring.v1"
)

func SyncProjects() {
	logrus.Info("=> Sincronizando projetos no Gitlab.")

	load.Start()

	var nProjects = make(ProjectsT)
	var lpCount = 0
	var page = 1
	var tentatives = 1

	for {
		lpo := gitlab.ListProjectsOptions{
			Membership: gitlab.Bool(true),
			Archived:   gitlab.Bool(false),
			Visibility: gitlab.Visibility("private"),
			Simple:     gitlab.Bool(true),
			OrderBy:    gitlab.String("name"),
			ListOptions: gitlab.ListOptions{
				PerPage: 50,
				Page:    page,
			},
		}

		lp, resp, err := glc.Projects.ListProjects(&lpo)
		if err != nil {
			logrus.Errorf("erro ao listar projetos \n %s", err)
			if tentatives <= 3 {
				tentatives++
				time.Sleep(5 * time.Second)
				continue
			} else {
				load.Stop()
				return
			}
		}

		for _, p := range lp {
			if p.PathWithNamespace == p.Path {

				id := strconv.Itoa(p.ID)

				if Projects[id] != nil {
					nProjects[id] = Projects[id]
					nProjects[id].Name = p.Name
					nProjects[id].LastActivityAt = p.LastActivityAt.UTC().In(timeLocation)
				} else {
					nProjects[id] = &Project{
						ID:                p.ID,
						Name:              p.Name,
						Path:              p.Path,
						PathWithNamespace: p.PathWithNamespace,
						LastActivityAt:    p.LastActivityAt.UTC().In(timeLocation),
					}

					lpCount++
				}
			}
		}

		page, _ = strconv.Atoi(resp.Header.Get("X-Next-Page"))
		if page <= 0 {
			break
		}
	}

	SaveProjects(&nProjects)

	load.Stop()

	logrus.Infof("=> Novos projetos encontrados: %v", lpCount)
	logrus.Info("=> Projetos sincronizadas do Gitlab \n\n")
}

func SyncProjectsCompanies() {
	logrus.Info("=> Sincronizando projetos com as empresas.")

	load.Start()

	for i, project := range Projects {
		if project.CompanyId != "" {
			continue
		}

		projectNameLow := strings.ToLower(project.Name)

		nameSplit := strings.Split(projectNameLow, " ")
		nameS := projectNameLow

		if len(nameSplit) >= 2 {
			nameS = nameSplit[0] + " " + nameSplit[1]
		}

		m := substring.Or(
			substring.Has(projectNameLow),
			substring.Has(nameS),
			substring.Has(strings.Replace(projectNameLow, "", "", -1)),
			substring.Has(project.Path),
		)

		for _, company := range Companies {

			companyNameLow := strings.ToLower(company.Name)

			b := make([]byte, len(companyNameLow))
			t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
			nDst, _, err := t.Transform(b, []byte(companyNameLow), true)

			if err != nil {
				logrus.Errorf("erro ao remover caracteres especiais do nome da empresa. \n %s", err)
				continue
			}

			if m.Match(companyNameLow) || m.Match(string(b[:nDst])) || m.Match(strings.Replace(companyNameLow, " ", "", -1)) {
				Projects[i].CompanyId = company.Id
			}
		}

		//if Projects[i].CompanyId == "" {
		//	fmt.Println(Projects[i].Name)
		//}
	}

	SaveProjects(&Projects)

	load.Stop()

	logrus.Info("=> Projetos com as empresas sincronizadas \n\n")
}
