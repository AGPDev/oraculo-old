package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
	"time"
)

func SyncCompanies() {
	logrus.Info("=> Sincronizando empresas com o Movidesk.")

	load.Start()

	var companiesApi []CompanyApi
	var nCompanies = make(CompaniesT)
	var newC = 0
	var skip = 0

	r := resty.R()
	r.SetQueryParam("token", Cfg.MovideskToken)
	r.SetQueryParam("$filter", "personType eq 2 and isActive eq true and profileType eq 2")
	r.SetQueryParam("$select", "id,businessName,codeReferenceAdditional,changedDate")
	r.SetQueryParam("$orderBy", "businessName asc")

	for {

		r.SetQueryParam("$skip", strconv.Itoa(skip))

		resp, err := r.Get("/persons")
		if err != nil {
			logrus.Errorf("erro ao listar as empresas do Movidesk \n %v", err)
			return
		}

		err = json.Unmarshal(resp.Body(), &companiesApi)
		if err != nil {
			logrus.Errorf("erro ao parsear as empresas \n %v", err)
			return
		}

		companyLen := len(companiesApi)
		if companyLen <= 0 {
			break
		} else {
			skip += companyLen
		}

		for _, company := range companiesApi {

			changedDate, _ := time.Parse("2006-01-02T15:04:05", company.ChangedDate)

			companyName := strings.TrimSpace(company.BusinessName)

			if Companies[company.Id] != nil {
				nCompanies[company.Id] = Companies[company.Id]
				nCompanies[company.Id].Name = companyName
				nCompanies[company.Id].UpdatedAt = changedDate
			} else {
				nCompanies[company.Id] = &Company{
					Id:        company.Id,
					Name:      companyName,
					UpdatedAt: changedDate,
				}
				newC++
			}
		}
	}

	SaveCompanies(&nCompanies)

	load.Stop()

	logrus.Infof("=> Novas empresas encontradas: %v", newC)
	logrus.Info("=> Empresas sincronizadas do Movidesk \n\n")
}