package api_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/m-mizutani/task-kitchen"
)

func TestNewTask(t *testing.T) {
	mgr := main.NewKitchenManager(testCfg.TableRegion, testCfg.TableName)
	uid1 := uuid.New().String()
	now := time.Now()

	t1, err := mgr.NewTask(uid1, now)
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
	mgr := main.NewKitchenManager(testCfg.TableRegion, testCfg.TableName)
	uid1 := uuid.New().String()
	now := time.Now()

	t1, err := mgr.NewTask(uid1, now)
	require.NoError(t, err)

	t2, err := mgr.NewTask(uid1, now)
	require.NoError(t, err)

	tset, err := mgr.FetchTasks(uid1, now)
	require.NoError(t, err)
	require.Equal(t, 2, len(tset))
	assert.True(t, tset[0].TaskID == t1.TaskID || tset[1].TaskID == t1.TaskID)

	assert.NoError(t, t1.Delete())
	assert.NoError(t, t2.Delete())
}

func TestPomodoro(t *testing.T) {
	mgr := main.NewKitchenManager(testCfg.TableRegion, testCfg.TableName)
	uid1 := uuid.New().String()
	now := time.Now()

	t1, err := mgr.NewTask(uid1, now)
	require.NoError(t, err)
	t2, err := mgr.NewTask(uid1, now)
	require.NoError(t, err)

	// Create a pomodoro
	p1, err := main.NewPomodoro(t1)
	require.NoError(t, err)
	assert.Equal(t, "started", p1.Status)

	err = p1.Finish()
	require.NoError(t, err)

	// Create another pomodoro
	p2, err := main.NewPomodoro(t1)
	require.NoError(t, err)

	// Create yet another pomodoro for t2
	p3, err := main.NewPomodoro(t2)
	require.NoError(t, err)

	// Check fetch action and isolation
	pset, err := main.FetchPomodoros(t1)
	require.NoError(t, err)
	require.Equal(t, 2, len(pset))
	assert.True(t, pset[0].PomodoroID == p1.PomodoroID || pset[1].PomodoroID == p1.PomodoroID)

	for _, p := range pset {
		assert.NotEqual(t, p3.PomodoroID, p.PomodoroID)
	}

	// Teardown
	require.NoError(t, t1.Delete())
	require.NoError(t, t2.Delete())
	require.NoError(t, p1.Delete())
	require.NoError(t, p2.Delete())
	require.NoError(t, p3.Delete())
}
