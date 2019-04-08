package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Task Endpoint
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
