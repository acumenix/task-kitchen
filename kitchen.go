package kitchen

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

type KitchenManager struct {
	table     dynamo.Table
	tableName string
}

func newKitchenManager(region, tableName string) KitchenManager {
	cfg := &aws.Config{Region: aws.String(region)}
	db := dynamo.New(session.New(), cfg)

	kitchenMgr := KitchenManager{
		table:     db.Table(tableName),
		tableName: tableName,
	}

	return kitchenMgr
}

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
}
