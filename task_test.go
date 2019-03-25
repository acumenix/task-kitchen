package main_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/m-mizutani/task-kitchen"
)

func TestNewTask(t *testing.T) {
	mgr := main.NewTaskManager(testCfg.TableRegion, testCfg.TableName)
	uid1 := uuid.New().String()
	now := time.Now()

	t1 := mgr.NewTask(uid1, now)
	err := t1.Save()
	require.NoError(t, err)

	t2, err := mgr.GetTask(uid1, now, t1.TaskID)
	require.NoError(t, err)
	require.NotNil(t, t2)
	assert.Equal(t, uid1, t2.UserID)

	err = t1.Delete()
	require.NoError(t, err)

	t3, err := mgr.GetTask(uid1, now, t1.TaskID)
	require.NoError(t, err)
	require.Nil(t, t3)
}

func TestFetchTasks(t *testing.T) {
	mgr := main.NewTaskManager(testCfg.TableRegion, testCfg.TableName)
	uid1 := uuid.New().String()
	now := time.Now()

	t1 := mgr.NewTask(uid1, now)
	require.NoError(t, t1.Save())

	t2 := mgr.NewTask(uid1, now)
	require.NoError(t, t2.Save())

	tset, err := mgr.FetchTasks(uid1, now)
	require.NoError(t, err)
	require.Equal(t, 2, len(tset))
	assert.True(t, tset[0].TaskID == t1.TaskID || tset[1].TaskID == t1.TaskID)

	assert.NoError(t, t1.Delete())
	assert.NoError(t, t2.Delete())
}
