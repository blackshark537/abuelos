package jobs

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

type Job struct {
	Seconds int
}

func (j *Job) Hello() {
	cron := gocron.NewScheduler(time.UTC)
	go func() {
		cron.Every(j.Seconds).
			Seconds().Do(fmt.Printf("%s\n", "Hello Job!"))
	}()
}
