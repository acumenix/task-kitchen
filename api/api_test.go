package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m-mizutani/task-kitchen/api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var apiEndPoint = "127.0.0.1:23456"

func runTestServer() {
	api.Logger = logrus.New()
	api.Logger.SetLevel(logrus.DebugLevel)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	api.SetupRouter(v1, testCfg.TableRegion, testCfg.TableName)

	go func() {
		r.Run(apiEndPoint)
	}()
}

func httpRequest(method, path string, input interface{}, response interface{}) (int, error) {
	url := fmt.Sprintf("http://%s/api/v1/%s", apiEndPoint, path)
	var reader io.Reader
	if input != nil {
		raw, err := json.Marshal(input)
		if err != nil {
			return 0, err
		}

		reader = bytes.NewBuffer(raw)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	if response != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, err
		}
		if err := json.Unmarshal(respBody, response); err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}

func TestReportAPI(t *testing.T) {
	type Report struct {
		Results api.Report `json:"results,omitempty"`
	}
	type Reports struct {
		Results []api.Report `json:"results,omitempty"`
	}
	type Error struct {
		Error string `json:"error,omitempty"`
	}
	var (
		code int
		err  error
	)
	uid := strings.Replace(uuid.New().String(), "-", "", -1)

	var resp1 Reports
	code, err = httpRequest("GET", uid+"?begin=2018-03-21&end=2018-04-01", nil, &resp1)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 0, len(resp1.Results))

	var resp2 Report
	code, err = httpRequest("GET", uid+"/2018-03-22", nil, &resp2)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, uid, resp2.Results.UserID)
	assert.Equal(t, "edit", string(resp2.Results.Status))

	var resp3 Reports
	code, err = httpRequest("GET", uid+"?begin=2018-03-21&end=2018-04-01", nil, &resp3)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(resp3.Results))

	code, err = httpRequest("GET", uid+"?begin=2018-03-22&end=2018-04-01", nil, &resp3)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(resp3.Results))

	code, err = httpRequest("GET", uid+"?begin=2018-03-23&end=2018-04-01", nil, &resp3)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 0, len(resp3.Results))

	code, err = httpRequest("GET", uid+"?begin=2018-03-22&end=2018-03-22", nil, &resp3)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 1, len(resp3.Results))

	report := api.Report{
		Status: api.ReportWorking,
	}
	code, err = httpRequest("PUT", uid+"/2018-03-22", report, nil)
	require.NoError(t, err)
	assert.Equal(t, 200, code)

	var resp4 Report
	code, err = httpRequest("GET", uid+"/2018-03-22", nil, &resp4)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, "work", string(resp4.Results.Status))

	report.Status = api.ReportStatus("invalid")
	code, err = httpRequest("PUT", uid+"/2018-03-22", report, nil)
	require.NoError(t, err)
	assert.Equal(t, 400, code)

	var resp5 Error
	code, err = httpRequest("DELETE", uid+"/2018-03-22", nil, &resp5)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, "", resp5.Error)
}

func TestTaskAPI(t *testing.T) {
	type Task struct {
		Results api.Task `json:"results,omitempty"`
	}
	type Tasks struct {
		Results []api.Task `json:"results,omitempty"`
	}
	var (
		code int
		err  error
	)

	var tasks Tasks
	uid := strings.Replace(uuid.New().String(), "-", "", -1)
	code, err = httpRequest("GET", uid+"/2018-03-22/task", nil, &tasks)
	require.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 0, len(tasks.Results))
}

func TestPomodoroAPI(t *testing.T) {
	type Task struct {
		Results api.Task `json:"results,omitempty"`
	}
	type Tasks struct {
		Results []api.Task `json:"results,omitempty"`
	}
	type Pomodoro struct {
		Results api.Pomodoro `json:"results,omitempty"`
	}
	type Pomodoros struct {
		Results []api.Pomodoro `json:"results,omitempty"`
	}
	var (
		code      int
		err       error
		task      Task
		pomodoro  Pomodoro
		pomodoros Pomodoros
	)
	uid := strings.Replace(uuid.New().String(), "-", "", -1)

	code, err = httpRequest("GET", uid+"/1983-04-20/pomodoro", nil, &pomodoros)
	require.NoError(t, err)
	require.Equal(t, 200, code)
	assert.Equal(t, 0, len(pomodoros.Results))

	code, err = httpRequest("POST", uid+"/1983-04-20/task", nil, &task)
	require.NoError(t, err)
	require.Equal(t, 200, code)

	code, err = httpRequest("POST", uid+"/1983-04-20/pomodoro/"+task.Results.TaskID, nil, &pomodoro)
	require.NoError(t, err)
	require.Equal(t, 200, code)

	code, err = httpRequest("GET", uid+"/1983-04-20/pomodoro/"+task.Results.TaskID+"/"+pomodoro.Results.PomodoroID, nil, &pomodoro)
	require.NoError(t, err)
	require.Equal(t, 200, code)

	code, err = httpRequest("GET", uid+"/1983-04-20/pomodoro/"+task.Results.TaskID+"/zatsu", nil, &pomodoro)
	require.NoError(t, err)
	require.Equal(t, 404, code)

	// Remote task
	code, err = httpRequest("DELETE", uid+"/1983-04-20/task/"+task.Results.TaskID, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 200, code)

	// Pomodoro belong to deleted task should be hidden
	code, err = httpRequest("GET", uid+"/1983-04-20/pomodoro/"+task.Results.TaskID+"/"+pomodoro.Results.PomodoroID, nil, &pomodoro)
	require.NoError(t, err)
	require.Equal(t, 404, code)
}
