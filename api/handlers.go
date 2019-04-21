package api

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// --------------------------------
// Utilities
// --------------------------------

const (
	charsetAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigit    = "0123456789"
)

func containsOnly(s, chars string) bool {
	for _, c := range s {
		if !strings.ContainsAny(string(c), chars) {
			return false
		}
	}

	return true
}

func getParam(params gin.Params, key string) string {
	for _, p := range params {
		if p.Key == key {
			return p.Value
		}
	}

	return ""
}

func getUser(params gin.Params) (string, error) {
	user := getParam(params, "user")

	if user == "" {
		return "", newUserError(400, "user parameter is empty")
	}

	if !containsOnly(user, charsetAlphabet+charsetDigit+"@_-") {
		return "", newUserError(400, "user parameter has invalid charactor")
	}

	return user, nil
}

func getTime(c *gin.Context, key string) (time.Time, error) {
	date, ok := c.GetQuery(key)
	if !ok {
		ts := time.Now()
		return ts, newUserError(400, "Missing required query string: %s", key)
	}

	ts, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ts, newUserError(400, "Invalid date format '%s', should be like 2006-01-02", date)
	}

	return ts, nil
}

func getSpace(params gin.Params) (user string, ts time.Time, err error) {
	if user, err = getUser(params); err != nil {
		return
	}

	date := getParam(params, "date")
	ts, err = time.Parse("2006-01-02", date)
	if err != nil {
		err = newUserError(400, "Invalid date format '%s', should be like 2006-01-02", date).setCause(err)
		return
	}

	return
}

type Response struct {
	Error     string      `json:"error,omitempty"`
	Results   interface{} `json:"results,omitempty"`
	RequestID string      `json:"request_id"`
}

type handler func(c *gin.Context, mgr *KitchenManager) (interface{}, error)

func handle(hdlr handler, c *gin.Context, mgr *KitchenManager) {
	reqID := uuid.New().String()
	result, err := hdlr(c, mgr)
	if err != nil {
		if userErr, ok := err.(*userError); ok {
			// User oriented error
			c.JSON(userErr.code, Response{userErr.Error(), nil, reqID})
		} else {
			// System oriented error
			c.JSON(500, Response{"Internal server error", nil, reqID})
		}
	} else {
		c.JSON(200, Response{"", result, reqID})
	}
}

// --------------------------------
// Report endpoints
// --------------------------------

func fetchReportHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, err := getUser(c.Params)
	if err != nil {
		return nil, err
	}

	begin, err := getTime(c, "begin")
	if err != nil {
		return nil, err
	}
	end, err := getTime(c, "end")
	if err != nil {
		return nil, err
	}

	reports, err := mgr.FetchReport(user, begin, end)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func getReportHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	report, err := mgr.GetReport(user, ts)
	if err != nil {
		return nil, err
	}

	if report == nil {
		report, err = mgr.NewReport(user, ts)
		if err != nil {
			return nil, err
		}
	}

	return report, nil
}

func getReportRoutine(c *gin.Context, mgr *KitchenManager) (*Report, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	report, err := mgr.GetReport(user, ts)
	if err != nil {
		return nil, err
	} else if report == nil {
		return nil, newUserError(404, "The report is not found")
	}

	return report, nil
}

func updateReportHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	report, err := getReportRoutine(c, mgr)
	if err != nil {
		return nil, err
	}

	var updatedReport Report
	c.BindJSON(&updatedReport)
	report.Status = updatedReport.Status

	if err := report.Save(); err != nil {
		return nil, err
	}

	return nil, nil
}

func deleteReportHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	report, err := getReportRoutine(c, mgr)
	if err != nil {
		return nil, err
	}

	if err := report.Delete(); err != nil {
		return nil, err
	}

	return nil, nil
}

// --------------------------------
// Task endpoints
// --------------------------------

func getTasksHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	tasks, err := mgr.FetchTasks(user, ts)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func createTaskHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	task, err := mgr.NewTask(user, ts)
	if err != nil {
		return nil, err
	}

	var reqTask Task
	if err := c.ShouldBindJSON(&reqTask); err == nil && reqTask.Title != "" {
		task.Title = reqTask.Title
		if err := task.Save(); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func updateTaskHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	taskID := getParam(c.Params, "task_id")
	task, err := mgr.GetTask(user, ts, taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, newUserError(404, "Task not found: %s", taskID)
	}

	var updatedTask Task
	c.BindJSON(&updatedTask)
	task.Title = updatedTask.Title
	task.TomatoNum = updatedTask.TomatoNum
	task.Description = updatedTask.Description

	if err := task.Save(); err != nil {
		return nil, err
	}

	return nil, nil
}

func deleteTaskHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	taskID := getParam(c.Params, "task_id")
	task, err := mgr.GetTask(user, ts, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, newUserError(404, "Task not found: %s", taskID)
	}

	if err := task.Delete(); err != nil {
		return nil, err
	}

	return nil, nil
}

// --------------------------------
// Chore endpoints
// --------------------------------

func fetchChoresHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	chores, err := mgr.FetchChores(user, ts)
	if err != nil {
		return nil, err
	}

	return chores, nil
}

func createChoreHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	chore, err := mgr.NewChore(user, ts)
	if err != nil {
		return nil, err
	}

	var reqChore Chore
	if err := c.ShouldBindJSON(&reqChore); err == nil && reqChore.Title != "" {
		chore.Title = reqChore.Title
		if err := chore.Save(); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return chore, nil
}

func updateChoreHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	choreID := getParam(c.Params, "chore_id")
	chore, err := mgr.GetChore(user, ts, choreID)
	if err != nil {
		return nil, err
	}
	if chore == nil {
		return nil, newUserError(404, "Chore not found: %s", choreID)
	}

	var updatedChore Chore
	c.BindJSON(&updatedChore)
	chore.Title = updatedChore.Title

	if err := chore.Save(); err != nil {
		return nil, err
	}

	return nil, nil
}

func deleteChoreHandler(c *gin.Context, mgr *KitchenManager) (interface{}, error) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		return nil, err
	}

	choreID := getParam(c.Params, "chore_id")
	chore, err := mgr.GetChore(user, ts, choreID)
	if err != nil {
		return nil, err
	}
	if chore == nil {
		return nil, newUserError(404, "Chore not found: %s", choreID)
	}

	if err := chore.Delete(); err != nil {
		return nil, err
	}

	return nil, nil
}
