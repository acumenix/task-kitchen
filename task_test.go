package main_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/m-mizutani/task-kitchen"
)

func TestNewTask(t *testing.T) {
	uid1 := uuid.New().String()
	now := time.Now()

	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(testCfg.TableRegion)})
	table := db.Table(testCfg.TableName)

	t1 := main.NewTask(table, uid1, now)
	err := t1.Save()
	require.NoError(t, err)

	t2, err := main.GetTask(table, uid1, now, t1.TaskID)
	require.NoError(t, err)
	require.NotNil(t, t2)
	assert.Equal(t, uid1, t2.UserID)

	err = t1.Delete()
	require.NoError(t, err)

	t3, err := main.GetTask(table, uid1, now, t1.TaskID)
	require.NoError(t, err)
	require.Nil(t, t3)
}
