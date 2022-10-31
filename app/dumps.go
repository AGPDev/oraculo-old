package app

import (
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func SyncDumps() {
	if lock {
		return
	}

	lock = true

	logrus.Info("=> Sincronizando Dumps.")

	load.Start()

	timeNow = time.Now()

	var nDumps = make(DumpsT)
	var nDumpsC = 0

	for i, job := range Jobs {

		if job.Dump != nil && job.FinishedAt.Day() == timeNow.Day() {
			nDumps[i] = job.Dump
			continue
		}

		// pega o arquivo com o resultado do job
		tf, _, err := glc.Jobs.GetTraceFile(job.ProjectId, job.ID)
		if err != nil {
			logrus.Errorf("erro ao buscar o resultado do job %v \n %v", job.ID, err)
			continue
		}

		// decodifica o arquivo do job
		tfb, err := io.ReadAll(tf)
		if err != nil {
			logrus.Errorf(
				"erro ao decodificar o resultado do Job %v \n %v",
				Projects[strconv.Itoa(job.ProjectId)].Name,
				err,
			)
			continue
		}

		s := string(tfb)

		// busca pelo link do dump
		r, err := regexp.Compile("((https://dbdumps)([[:punct:]])(.*)(.gz\"))")
		if err != nil {
			logrus.Errorf("erro ao compilar o regex \n %v", err)
			continue
		}
		s = r.FindString(s)
		if s == "" {
			continue
		}

		fileUrl := strings.Replace(s, "\"", "", -1)

		// prepara os nomes do arquivo
		urlSlice := strings.Split(fileUrl, "/")
		fileName := urlSlice[len(urlSlice)-1]

		fileNameEasy := strings.Split(fileName, "_")[0]

		nDumps[i] = &Dump{
			CompanyId:     Projects[strconv.Itoa(job.ProjectId)].CompanyId,
			ProjectId:     job.ProjectId,
			JobId:         job.ID,
			Url:           fileUrl,
			Name:          fileName,
			NameEasy:      fileNameEasy,
			FilePath:      dirName + "/" + fileName,
			FilePathToday: "www/dumps/today/" + fileNameEasy + ".sql.gz",
		}

		Jobs[i].Dump = nDumps[i]

		nDumpsC++
	}

	SaveJobs(&Jobs)
	SaveDumps(&nDumps)

	lock = false

	load.Stop()

	logrus.Infof("=> Novos dumps encontrados: %v", nDumpsC)
	logrus.Info("=> Dumps sincronizados \n\n")
}
