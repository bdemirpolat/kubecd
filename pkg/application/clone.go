package application

import (
	"fmt"
	"github.com/bdemirpolat/kubecd/pkg/application/k8apply"
	"github.com/bdemirpolat/kubecd/pkg/logger"
	"github.com/bdemirpolat/kubecd/pkg/models"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io/ioutil"
	"os"
	"time"
)

var clonePath = "/var/tmp/kubecd"

const failedLastStatusMessage = "Last sync for application (Name: %s) failed. Description: %s"
const successLastStatusMessage = "Last sync for application (Name: %s) synced successfully."

func init() {
	if cp := os.Getenv("KUBECD_CLONE_PATH"); cp != "" {
		clonePath = cp
	}
}

func clone(repo RepoInterface, application models.Application) {
	defer func() {
		application.LastCheck = time.Now()
		_, err := repo.Update(&application)
		if err != nil {
			logger.SugarLogger.Error(err)
		}
	}()
	repoPath := fmt.Sprintf("%s/%s", clonePath, application.Name)
	fmt.Println("repo path:", repoPath)
	_ = os.RemoveAll(repoPath)
	r, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{Username: application.Username, Password: application.Token},
		URL:  application.URL,
	})

	if err != nil {
		logger.SugarLogger.Error(err)
		application.LastStatusMessage = fmt.Sprintf(failedLastStatusMessage, application.Name, err.Error())
		return
	}

	head, err := r.Head()
	fmt.Println("headerr:", err)
	if err != nil {
		logger.SugarLogger.Error(err)
		application.LastStatusMessage = fmt.Sprintf(failedLastStatusMessage, application.Name, err.Error())
		return
	}

	application.Head = head.String()

	files, readDirErr := readDir(fmt.Sprintf("%s/%s", repoPath, application.ManifestDir))
	fmt.Println("readDirErr:", readDirErr)
	if readDirErr != nil {
		logger.SugarLogger.Error(readDirErr)
		application.LastStatusMessage = fmt.Sprintf(failedLastStatusMessage, application.Name, err.Error())
		return
	}

	for _, f := range files {
		resultChan := make(chan error)
		k8apply.AddToApplyQueue(f, resultChan)
		applyErr := <-resultChan
		fmt.Println("apply err:", applyErr)
		if applyErr != nil {
			logger.SugarLogger.Error(applyErr)
			application.LastStatusMessage = fmt.Sprintf(failedLastStatusMessage, application.Name, err.Error())
			return
		}
	}

	application.LastStatusMessage = fmt.Sprintf(successLastStatusMessage, application.Name)

}

func readDir(dir string) ([][]byte, error) {
	var reads [][]byte
	fmt.Println("dir reading:", dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			b, readDirErr := readDir(fmt.Sprintf("%s/%s", dir, file.Name()))
			if readDirErr != nil {
				return nil, err
			}
			reads = append(reads, b...)
		} else {
			f, openErr := os.Open(fmt.Sprintf("%s/%s", dir, file.Name()))
			if openErr != nil {
				return nil, openErr
			}

			d, readAllErr := ioutil.ReadAll(f)
			if readAllErr != nil {
				return nil, readAllErr
			}
			reads = append(reads, d)
		}
	}

	return reads, nil
}
