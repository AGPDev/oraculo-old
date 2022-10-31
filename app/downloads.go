package app

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/sirupsen/logrus"
)

func SyncFolder() {

	logrus.Info("=> Verificando pasta de downloads do dia")

	load.Start()

	timeNow = time.Now()
	dirName = fmt.Sprintf("www/dumps/%v_%v_%v", timeNow.Year(), int(timeNow.Month()), timeNow.Day())

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		logrus.Errorf("não foi possível criar o diretório %v", dirName)
	}

	load.Stop()

	logrus.Info("=> Pasta de downloads do dia sincronizada \n\n")
}

func SyncDownloads() {
	if lock {
		return
	}

	lock = true

	logrus.Info("=> Sincronizando Downloads")

	load.Start()

	timeNow = time.Now()

	for _, dump := range Dumps {

		if Jobs[strconv.Itoa(dump.JobId)].FinishedAt.Day() != timeNow.Day() {
			continue
		}

		fileInfo, _ := os.Stat(dump.FilePath)
		fileTodayInfo, _ := os.Stat(dump.FilePathToday)

		if fileInfo == nil {
			load.Stop()

			err := downloadFile(dump)
			if err != nil {

				logrus.Error(err)
				load.Start()

				fileInfo, _ := os.Stat(dump.FilePath)
				if fileInfo != nil {
					_ = os.Remove(dump.FilePath)
				}
				continue
			}

			load.Start()
		}

		if fileTodayInfo != nil {
			_ = os.Remove(dump.FilePathToday)
		}

		copyDumpToday(dump)
	}

	lock = false

	load.Stop()

	logrus.Info("=> Downloads sincronizados \n\n")
}

func SyncClearToday() {
	logrus.Info("=> Limpando pasta de dumps do dia")

	load.Start()

	files, err := filepath.Glob(filepath.Join("www/dumps/today/", "*.sql.gz"))

	if err != nil {
		logrus.Errorf("não foi possível listar os arquivos de dumps do dia \n %v", err)
	}

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			logrus.Errorf("não foi possível limpar pasta de dumps do dia \n %v", err)
		}
	}

	load.Stop()

	logrus.Info("=> Pasta de dumps do dia limpa \n\n")
}

func downloadFile(dump *Dump) error {
	// prepara a request para download
	req, err := grab.NewRequest(dump.FilePath, dump.Url)
	if err != nil {
		return err
	}

	// prepara o cliente para download
	client := grab.NewClient()
	client.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.96 Safari/537.36"
	//client.BufferSize = 512

	// executa o download
	logrus.Infof("=> Baixando %v", req.URL())
	resp := client.Do(req)

	// máximo de tempo de download de acordo com o tamanho do arquivo
	maxEta := (resp.Size / 1048576) / 10
	if maxEta < 5 {
		maxEta = 5
	}

	// inicializa o loop UI
	t := time.NewTicker(time.Second)
	defer t.Stop()

	// contador de tempo do download lento
	lowDtry := 0

Loop:
	for {
		select {
		case <-t.C:

			eta := resp.ETA()
			minEta := eta.Sub(time.Now()).Minutes()
			minEta = math.Round(minEta)

			// exibe o status do download
			fmt.Printf(" transferido %v / %v M (%.2f%%) - %.2f M/s ETA %vm \n",
				resp.BytesComplete()/1048576,
				resp.Size/1048576,
				100*resp.Progress(),
				resp.BytesPerSecond()/524288,
				minEta,
			)

			// verifica a velocidade do download
			if int64(minEta) > maxEta {
				// aumenta o contador caso esteja lento
				lowDtry++
			} else {
				// zera o contador se estiver boa a velocidade
				lowDtry = 0
			}

			if lowDtry >= 10 {
				err = errors.New("download muito lento")
				break Loop
			}
		case <-resp.Done:
			// Download completo
			break Loop
		}
	}

	if err == nil {
		err = resp.Err()
	}

	if err != nil {
		return err
	}

	return nil
}

func copyDumpToday(dump *Dump) {
	cmd := exec.Command("cp", dump.FilePath, dump.FilePathToday)
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("não foi possível compiar o arquivo para a pasta do dia \n %v", err)
	}
}
