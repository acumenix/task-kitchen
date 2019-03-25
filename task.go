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
	ssn   *session.Session
	cfg   *aws.Config
	table dynamo.Table
}

func NewTaskManager(region, tableName string) TaskManager {
	cfg := &aws.Config{Region: aws.String(region)}
	ssn := session.Must(session.NewSession(cfg))
	db := dynamo.New(session.New(), cfg)

	taskMgr := TaskManager{
		cfg:   cfg,
		ssn:   ssn,
		table: db.Table(tableName),
	}

	return taskMgr
}

func toPKey(userID string, date time.Time, taskID string) string {
	return fmt.Sprintf("%s/task/%s/%s", userID, date.Format("20060102"), taskID)
}
func toPKeyPrefix(userID string, date time.Time) string {
	return fmt.Sprintf("%s/task/%s/", userID, date.Format("20060102"))
}

func (x TaskManager) NewTask(userID string, date time.Time) *Task {
	task := Task{
		UserID:    userID,
		TaskID:    strings.Replace(uuid.New().String(), "-", "", -1),
		CreatedAt: date,
		table:     x.table,
	}

	task.PKey = toPKey(task.UserID, task.CreatedAt, task.TaskID)

	return &task
}

func (x TaskManager) GetTask(userID string, date time.Time, taskID string) (*Task, error) {
	var task Task
	if err := x.table.Get("pk", toPKey(userID, date, taskID)).One(&task); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Fail to get task")
	}

	return &task, nil
}

func (x TaskManager) FetchTasks(userID string, date time.Time) ([]Task, error) {
	var tasks []Task

	query := x.table.Scan().Filter("begins_with($, ?)", "pk", toPKeyPrefix(userID, date))

	if err := query.All(&tasks); err != nil {
		return nil, errors.Wrapf(err, "Fail to fetch tasks: %s %s", userID, date)
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
	if err := x.table.Delete("pk", x.PKey).Run(); err != nil {
		return errors.Wrapf(err, "Fail to delete task: %s", x.PKey)
	}

	x.deleted = true
	return nil
}
