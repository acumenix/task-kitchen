package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

type Chore struct {
	PKey        string    `dynamo:"pk" json:"-"`
	SKey        string    `dynamo:"sk" json:"-"`
	UserID      string    `dynamo:"user_id" json:"user_id"`
	ChoreID     string    `dynamo:"chore_id" json:"chore_id"`
	CreatedAt   time.Time `dynamo:"created_at" json:"created_at"`
	Title       string    `dynamo:"title" json:"title"`
	Done        bool      `dynamo:"done" json:"done"`
	Description string    `dynamo:"description" json:"description"`

	table   dynamo.Table
	deleted bool
}

func toChoreKey(userID string, date time.Time, choreID string) (string, string) {
	pk := fmt.Sprintf("%s/chore/%s", userID, date.Format("20060102"))
	sk := choreID
	return pk, sk
}

func (x KitchenManager) NewChore(userID string, date time.Time) (*Chore, error) {
	chore := Chore{
		UserID:    userID,
		ChoreID:   strings.Replace(uuid.New().String(), "-", "", -1),
		CreatedAt: date,
		table:     x.table,
	}

	chore.PKey, chore.SKey = toChoreKey(chore.UserID, chore.CreatedAt, chore.ChoreID)
	if err := chore.Save(); err != nil {
		return nil, err
	}

	return &chore, nil
}

func (x KitchenManager) GetChore(userID string, date time.Time, ChoreID string) (*Chore, error) {
	var Chore Chore
	pk, sk := toChoreKey(userID, date, ChoreID)

	if err := x.table.Get("pk", pk).Range("sk", dynamo.Equal, sk).One(&Chore); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Fail to get Chore")
	}

	Chore.table = x.table

	return &Chore, nil
}

func (x KitchenManager) FetchChores(userID string, date time.Time) ([]Chore, error) {
	var chores []Chore
	pk, _ := toChoreKey(userID, date, "")

	if err := x.table.Get("pk", pk).All(&chores); err != nil {
		if err.Error() == "dynamo: no item found" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "Fail to get chore")
	}

	return chores, nil
}

func (x *Chore) Save() error {
	if err := x.table.Put(x).Run(); err != nil {
		return errors.Wrapf(err, "Fail to save chore: %s", x.PKey)
	}

	return nil
}

func (x *Chore) Delete() error {
	if err := x.table.Delete("pk", x.PKey).Range("sk", x.SKey).Run(); err != nil {
		return errors.Wrapf(err, "Fail to delete chore: %s", x.PKey)
	}

	x.deleted = true
	return nil
}
