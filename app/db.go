package app

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type (
	CompanyApi struct {
		Id                      string `json:"id"`
		BusinessName            string `json:"businessName"`
		CodeReferenceAdditional string `json:"codeReferenceAdditional"`
		ChangedDate             string `json:"changedDate"`
	}

	Company struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		UpdatedAt time.Time `json:"UpdatedAt"`
	}

	Project struct {
		ID                int       `json:"id"`
		CompanyId         string    `json:"company_id"`
		Name              string    `json:"name"`
		Path              string    `json:"path"`
		PathWithNamespace string    `json:"path_with_namespace"`
		Ignore            bool      `json:"ignore"`
		LastActivityAt    time.Time `json:"last_activity_at"`
		LastJobSuccess    *Job      `json:"last_job_success"`
	}

	Job struct {
		ID         int       `json:"id"`
		CompanyId  string    `json:"company_id"`
		ProjectId  int       `json:"project_id"`
		Status     string    `json:"status"`
		Stage      string    `json:"stage"`
		FinishedAt time.Time `json:"finished_at"`
		Dump       *Dump     `json:"dump"`
	}

	Dump struct {
		CompanyId     string `json:"company_id"`
		ProjectId     int    `json:"project_id"`
		JobId         int    `json:"job_id"`
		Url           string `json:"url"`
		Name          string `json:"name"`
		NameEasy      string `json:"name_easy"`
		FilePath      string `json:"file_path"`
		FilePathToday string `json:"file_path_today"`
	}

	CompaniesT map[string]*Company
	ProjectsT map[string]*Project
	JobsT map[string]*Job
	DumpsT map[string]*Dump
)

var (
	Companies = make(CompaniesT)
	Projects  = make(ProjectsT)
	Jobs      = make(JobsT)
	Dumps     = make(DumpsT)
)

func LoadCompanies() {
	err := loadDbFile("db/companies.json", &Companies)
	if err != nil {
		logrus.Panicf("erro ao carregar as empresas do banco de dados \n %v", err)
	}
}

func LoadProjects() {
	err := loadDbFile("db/projects.json", &Projects)
	if err != nil {
		logrus.Panicf("erro ao carregar os projetos do banco de dados \n %v", err)
	}
}

func LoadJobs() {
	err := loadDbFile("db/jobs.json", &Jobs)
	if err != nil {
		logrus.Panicf("erro ao carregar os jobs do banco de dados \n %v", err)
	}
}

func LoadDumps() {
	err := loadDbFile("db/dumps.json", &Dumps)
	if err != nil {
		logrus.Panicf("erro ao carregar os dumps do banco de dados \n %v", err)
	}
}

func SaveCompanies(nc *CompaniesT) {
	err := saveDbFile("db/companies.json", &nc)
	if err != nil {
		logrus.Panicf("erro ao salvar as empresas no banco de dados \n %v", err)
	}

	Companies = *nc
}

func SaveProjects(np *ProjectsT) {
	err := saveDbFile("db/projects.json", &np)
	if err != nil {
		logrus.Panicf("erro ao salvar os projetos no banco de dados \n %v", err)
	}

	Projects = *np
}

func SaveJobs(nj *JobsT) {
	err := saveDbFile("db/jobs.json", &nj)
	if err != nil {
		logrus.Panicf("erro ao salvar os jobs dos projetos no banco de dados \n %v", err)
	}

	Jobs = *nj
}

func SaveDumps(nd *DumpsT) {
	err := saveDbFile("db/dumps.json", &nd)
	if err != nil {
		logrus.Panicf("erro ao salvar os dumps dos projetos no banco de dados \n %v", err)
	}

	Dumps = *nd
}

func loadDbFile(jfName string, v interface{}) error {
	jFile, err := os.OpenFile(jfName, os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer jFile.Close()

	// le o conteúdo do json
	byteJson, _ := ioutil.ReadAll(jFile)
	_ = json.Unmarshal([]byte(byteJson), &v)

	return nil
}

func saveDbFile(jfName string, v interface{}) error {
	jm, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	jFile, err := os.OpenFile(jfName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer jFile.Close()

	err = jFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = jFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = jFile.Write(jm)
	if err != nil {
		return err
	}

	return nil
}
