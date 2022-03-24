package application

import (
	"context"
	"fmt"
	"time"

	"github.com/bdemirpolat/kubecd/pkg/models"
)

var cancel context.CancelFunc

// Loop starts loop for every application with given interval, if application config changes cancels to all go routines and creates again
func Loop(repo RepoInterface) error {
	// do not trigger cancel func at the first call fo Renew function
	if cancel != nil {
		cancel()
	}
	var c context.Context
	c, cancel = context.WithCancel(context.Background())
	return startLoop(c, repo)
}

func startLoop(c context.Context, repo RepoInterface) error {
	applications, err := repo.List(0, 0)
	if err != nil {
		return err
	}
	for _, app := range *applications {
		go func(app models.Application) {
			t := time.NewTicker(time.Millisecond * app.Interval)
			for {
				select {
				case <-t.C:
					fmt.Println(app.Name, "tick at", time.Now())
					go clone(repo, app)
				case <-c.Done():
					fmt.Println(app.Name, "done at", time.Now())
					return
				}
			}
		}(app)
	}

	return nil
}
