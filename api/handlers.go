package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// --------------------------------
// Utilities
// --------------------------------

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
		return "", errors.New("user parameter is empty")
	}

	return user, nil
}

func getTime(params gin.Params, key string) (time.Time, error) {
	date := getParam(params, key)
	ts, err := time.Parse("2006-01-02", date)
	if err != nil {
		Logger.WithError(err).Error("Format error")
		return ts, fmt.Errorf("Invalid date format '%s', should be like 2006-01-02", date)
	}

	return ts, nil
}

func getSpace(params gin.Params) (user string, ts time.Time, err error) {
	if user, err = getUser(params); err != nil {
		return
	}

	ts, err = getTime(params, "date")
	return
}

// --------------------------------
// Report endpoints
// --------------------------------

func fetchReportHandler(c *gin.Context, mgr *KitchenManager) {
	user, err := getUser(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
	}

	begin, err := getTime(c.Params, "begin")
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
	}
	end, err := getTime(c.Params, "end")
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
	}

	reports, err := mgr.FetchReport(user, begin, end)
	if err != nil {
		Logger.WithError(err).Error("Fail to get report")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"", reports})
}

func getReportHandler(c *gin.Context, mgr *KitchenManager) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	report, err := mgr.GetReport(user, ts)
	if err != nil {
		Logger.WithError(err).Error("Fail to get report")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	if report == nil {
		report, err = mgr.NewReport(user, ts)
		if err != nil {
			Logger.WithError(err).Error("Fail to get report")
			c.JSON(500, Response{"Internal server error", nil})
		}
	}

	c.JSON(200, Response{"", report})
}

func getReportRoutine(c *gin.Context, mgr *KitchenManager) *Report {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return nil
	}

	report, err := mgr.GetReport(user, ts)
	if err != nil {
		Logger.WithError(err).Error("Fail to lookup a report")
		c.JSON(500, Response{"Internal server error", nil})
		return nil
	} else if report == nil {
		c.JSON(404, Response{"No available report", nil})
		return nil
	}

	return report
}

func updateReportHandler(c *gin.Context, mgr *KitchenManager) {
	report := getReportRoutine(c, mgr)
	if report == nil {
		return
	}

	var updatedReport Report
	c.BindJSON(&updatedReport)
	report.Status = updatedReport.Status

	if err := report.Save(); err != nil {
		Logger.WithError(err).WithField("report", report).Error("Fail to update a report")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"ok", nil})
	return
}

func deleteReportHandler(c *gin.Context, mgr *KitchenManager) {
	report := getReportRoutine(c, mgr)
	if report == nil {
		return
	}

	if err := report.Delete(); err != nil {
		Logger.WithError(err).WithField("report", report).Error("Fail to delete a report")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"ok", nil})
	return
}

// --------------------------------
// Task endpoints
// --------------------------------

func getTasksHandler(c *gin.Context, mgr *KitchenManager) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	tasks, err := mgr.FetchTasks(user, ts)
	if err != nil {
		Logger.WithError(err).Error("Fail to fetch tasks")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"", tasks})
}

func createTaskHandler(c *gin.Context, mgr *KitchenManager) {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	task, err := mgr.NewTask(user, ts)
	if err != nil {
		Logger.WithError(err).Error("Fail to create a task")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	var reqTask Task
	if err := c.ShouldBindJSON(&reqTask); err == nil && reqTask.Title != "" {
		task.Title = reqTask.Title
		if err := task.Save(); err != nil {
			Logger.WithField("task", task).WithError(err).Error("Fail to update title of task")
			c.JSON(500, Response{"Internal server error", nil})
		}
	}

	c.JSON(200, Response{"ok", task})
}

func updateTaskHandler(c *gin.Context, mgr *KitchenManager) {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	taskID := getParam(c.Params, "task_id")
	task, err := mgr.GetTask(user, ts, taskID)
	if err != nil {
		Logger.WithError(err).Error("Fail to lookup a task")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	if task == nil {
		c.JSON(404, Response{fmt.Sprintf("Task not found: %s", taskID), nil})
		return
	}

	var updatedTask Task
	c.BindJSON(&updatedTask)
	task.Title = updatedTask.Title
	task.TomatoNum = updatedTask.TomatoNum
	task.Description = updatedTask.Description

	if err := task.Save(); err != nil {
		Logger.WithError(err).Error("Fail to update task")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"ok", nil})
	return
}

func deleteTaskHandler(c *gin.Context, mgr *KitchenManager) {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	taskID := getParam(c.Params, "task_id")
	task, err := mgr.GetTask(user, ts, taskID)
	if err != nil {
		Logger.WithError(err).Error("Fail to lookup a task")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	if task == nil {
		c.JSON(404, Response{fmt.Sprintf("Task not found: %s", taskID), nil})
		return
	}

	if err := task.Delete(); err != nil {
		Logger.WithError(err).Error("Fail to delete a task")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"ok", nil})
	return
}

// --------------------------------
// Chore endpoints
// --------------------------------

func fetchChoresHandler(c *gin.Context, mgr *KitchenManager) {
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	chores, err := mgr.FetchChores(user, ts)
	if err != nil {
		Logger.WithError(err).Error("Fail to fetch chores")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"", chores})
}

func createChoreHandler(c *gin.Context, mgr *KitchenManager) {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	chore, err := mgr.NewChore(user, ts)
	if err != nil {
		Logger.WithError(err).Error("Fail to create a Chore")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	var reqChore Chore
	if err := c.ShouldBindJSON(&reqChore); err == nil && reqChore.Title != "" {
		chore.Title = reqChore.Title
		if err := chore.Save(); err != nil {
			Logger.WithField("Chore", chore).WithError(err).Error("Fail to update title of Chore")
			c.JSON(500, Response{"Internal server error", nil})
		}
	}

	c.JSON(200, Response{"ok", chore})
}

func updateChoreHandler(c *gin.Context, mgr *KitchenManager) {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	choreID := getParam(c.Params, "chore_id")
	chore, err := mgr.GetChore(user, ts, choreID)
	if err != nil {
		Logger.WithError(err).Error("Fail to lookup a Chore")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	if chore == nil {
		c.JSON(404, Response{fmt.Sprintf("Chore not found: %s", choreID), nil})
		return
	}

	var updatedChore Chore
	c.BindJSON(&updatedChore)
	chore.Title = updatedChore.Title

	if err := chore.Save(); err != nil {
		Logger.WithError(err).Error("Fail to update Chore")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"ok", nil})
	return
}

func deleteChoreHandler(c *gin.Context, mgr *KitchenManager) {
	Logger.WithField("param", c.Params).Info("Request")
	user, ts, err := getSpace(c.Params)
	if err != nil {
		c.JSON(400, Response{err.Error(), nil})
		return
	}

	ChoreID := getParam(c.Params, "Chore_id")
	Chore, err := mgr.GetChore(user, ts, ChoreID)
	if err != nil {
		Logger.WithError(err).Error("Fail to lookup a Chore")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	if Chore == nil {
		c.JSON(404, Response{fmt.Sprintf("Chore not found: %s", ChoreID), nil})
		return
	}

	if err := Chore.Delete(); err != nil {
		Logger.WithError(err).Error("Fail to delete a Chore")
		c.JSON(500, Response{"Internal server error", nil})
		return
	}

	c.JSON(200, Response{"ok", nil})
	return
}
