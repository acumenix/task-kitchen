package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

type Pomodoro struct {
	PKey       string `dynamo:"pk"`
	SKey       string `dynamo:"sk"`
	PomodoroID string `dynamo:"pomorodo_id"`

	Status      string    `dynamo:"status"`
	StartedAt   time.Time `dynamo:"started_at"`
	FinishedAt  time.Time `dynamo:"finished_at"`
	IsScheduled bool      `dynamo:"is_scheduled"`

	table   dynamo.Table
	deleted bool
}

func toPomodoroKey(userID string, date time.Time, taskID string, pomodoroID string) (string, string) {
	pk := fmt.Sprintf("%s/pomodoro/%s", userID, date.Format("20060102"))
	sk := fmt.Sprintf("%s/%s", taskID, pomodoroID)
	return pk, sk
}

func newPomodoro(task *Task) (*Pomodoro, error) {
	pID := uuid.New().String()
	pk, sk := toPomodoroKey(task.UserID, task.CreatedAt, task.TaskID, pID)
	p := new(Pomodoro)

	p.PKey = pk
	p.SKey = sk
	p.PomodoroID = pID
	p.Status = "started"
	p.StartedAt = time.Now().UTC()

	p.table = task.table

	if err := p.table.Put(p).Run(); err != nil {
		return p, errors.Wrapf(err, "Fail to put a new promodoro: %s, %s", pk, sk)
	}

	return p, nil
}

func fetchPomodoros(task *Task) ([]Pomodoro, error) {
	var pomodoros []Pomodoro
	pk, sk := toPomodoroKey(task.UserID, task.CreatedAt, task.TaskID, "")

	if err := task.table.Get("pk", pk).Range("sk", dynamo.BeginsWith, sk).All(&pomodoros); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Fail to get task")
	}

	return pomodoros, nil
}

func (x *Pomodoro) Finish() error {
	if x.deleted {
		Logger.WithField("pomodoro", x).Fatal("Already deleted")
	}

	x.FinishedAt = time.Now().UTC()
	x.Status = "finished"

	if err := x.table.Put(x).Run(); err != nil {
		return errors.Wrapf(err, "Fail to update the promodoro to finish: %s, %s", x.PKey, x.SKey)
	}

	return nil
}

func (x *Pomodoro) Delete() error {
	if x.deleted {
		Logger.WithField("pomodoro", x).Fatal("Already deleted")
	}

	if err := x.table.Delete("pk", x.PKey).Range("sk", x.SKey).Run(); err != nil {
		return errors.Wrapf(err, "Fail to delete pomodoro: %s", x.PKey)
	}

	x.deleted = true
	return nil

}
