package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

type Task struct {
	PKey        string    `dynamo:"pk"`
	SKey        string    `dynamo:"sk"`
	UserID      string    `dynamo:"user_id"`
	TaskID      string    `dynamo:"task_id"`
	CreatedAt   time.Time `dynamo:"created_at"`
	Title       string    `dynamo:"title"`
	TomatoNum   string    `dynamo:"tomato_num"`
	Description string    `dynamo:"description"`
	table       dynamo.Table
	deleted     bool
}

type TaskManager struct {
	ssn       *session.Session
	cfg       *aws.Config
	table     dynamo.Table
	tableName string
}

func NewTaskManager(region, tableName string) TaskManager {
	cfg := &aws.Config{Region: aws.String(region)}
	ssn := session.Must(session.NewSession(cfg))
	db := dynamo.New(session.New(), cfg)

	taskMgr := TaskManager{
		cfg:       cfg,
		ssn:       ssn,
		table:     db.Table(tableName),
		tableName: tableName,
	}

	return taskMgr
}

func toTaskKey(userID string, date time.Time, taskID string) (string, string) {
	pk := fmt.Sprintf("%s/task/%s", userID, date.Format("20060102"))
	sk := taskID
	return pk, sk

}

func (x TaskManager) NewTask(userID string, date time.Time) (*Task, error) {
	task := Task{
		UserID:    userID,
		TaskID:    strings.Replace(uuid.New().String(), "-", "", -1),
		CreatedAt: date,
		table:     x.table,
	}

	task.PKey, task.SKey = toTaskKey(task.UserID, task.CreatedAt, task.TaskID)
	if err := task.Save(); err != nil {
		return nil, err
	}

	return &task, nil
}

func (x TaskManager) GetTask(userID string, date time.Time, taskID string) (*Task, error) {
	var task Task
	pk, sk := toTaskKey(userID, date, taskID)

	if err := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).One(&task); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Fail to get task")
	}

	return &task, nil
}

func (x TaskManager) FetchTasks(userID string, date time.Time) ([]Task, error) {
	var tasks []Task
	pk, _ := toTaskKey(userID, date, "")

	if err := x.table.Get("pk", pk).All(&tasks); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Fail to get task")
	}

	return tasks, nil
}

func (x *Task) Save() error {
	if err := x.table.Put(x).Run(); err != nil {
		return errors.Wrapf(err, "Fail to save task: %s", x.PKey)
	}

	return nil
}

func (x *Task) Delete() error {
	if err := x.table.Delete("pk", x.PKey).Range("sk", x.SKey).Run(); err != nil {
		return errors.Wrapf(err, "Fail to delete task: %s", x.PKey)
	}

	x.deleted = true
	return nil
}

func (x *Task) StartPomodoro() (*Pomodoro, error) {

	return nil, nil
}
